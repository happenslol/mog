package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/happenslol/mog/codegen"
	"github.com/happenslol/mog/templates"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Output struct {
		Models      OutputConfig
		Collections OutputConfig
	}

	Primitives  []string
	Collections map[string]string
}

type OutputConfig struct {
	Package  string
	Filename string
}

var defaultPrimitives = []string{
	"time.Time",
}

var defaultModelImports = []string{
	"time",
}

func main() {
	configPath := ""
	flag.StringVar(&configPath, "c", "mog.yml", "set config file location")
	flag.Parse()

	rawConfig, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("Failed to read config: %s", err.Error()))
	}

	var config Config
	if err := yaml.Unmarshal(rawConfig, &config); err != nil {
		panic(fmt.Sprintf("Failed to parse config: %s", err.Error()))
	}

	collections := templates.CollectionsInput{
		PackageName: config.Output.Collections.Package,
		Collections: make(map[string]templates.CollectionInput),
		Imports:     make([]string, 0),
	}

	models := templates.ModelsInput{
		PackageName: config.Output.Models.Package,
		Structs:     make([]codegen.Strct, 0),
		Imports:     make([]string, 0),
	}

	collectionImports := map[string]struct{}{}
	for _, imp := range defaultModelImports {
		collectionImports[imp] = struct{}{}
	}

	// If we just mark all primitives as already
	// known, they will be skipped when we reach them
	strcts := map[string]codegen.Strct{}
	for _, s := range config.Primitives {
		strcts[s] = codegen.Strct{Skip: true}
	}

	for _, s := range defaultPrimitives {
		strcts[s] = codegen.Strct{Skip: true}
	}

	for colName, col := range config.Collections {
		pkg, name, err := splitPackageUID(col)
		if err != nil {
			panic(err)
		}

		codegen.LoadModel(pkg, name, strcts)

		parts := strings.Split(pkg, "/")
		modelType := fmt.Sprintf("%s.%s", parts[len(parts)-1], name)

		// Find our struct from the struct list. It's guaranteed to
		// be there if the previous function did not error out.
		strct := strcts[col]
		idField := findIDType(strct)
		collectionImports[strct.StructImport] = struct{}{}

		collectionType := fmt.Sprintf("%sCollection", strings.Title(colName))

		result := templates.CollectionInput{
			ModelType:      modelType,
			IDType:         idField.Type,
			CollectionType: collectionType,
		}

		for imp := range idField.Imports {
			collectionImports[imp] = struct{}{}
		}

		collections.Collections[colName] = result
	}

	for _, s := range strcts {
		if s.Skip {
			continue
		}

		models.Structs = append(models.Structs, s)
	}

	for imp := range collectionImports {
		collections.Imports = append(collections.Imports, imp)
	}

	if err := templates.WriteCollections(collections,
		config.Output.Collections.Filename); err != nil {
		panic(err)
	}

	if err := templates.WriteModels(models,
		config.Output.Models.Filename); err != nil {
		panic(err)
	}
}

func findIDType(strct codegen.Strct) codegen.Field {
	for _, fld := range strct.Fields {
		// If one of the existing fields is the index field, we don't
		// need any extra imports since it will already be imported
		// for the field helper methods
		if fld.BsonKey == "_id" {
			return fld
		}
	}

	// If we drop through, there's no explicit ID key and
	// we use primitive.ObjectID, which also needs an import.
	return codegen.Field{
		Type: "primitive.ObjectID",
		Imports: map[string]struct{}{
			"go.mongodb.org/mongo-driver/bson/primitive": {}}}
}

func splitPackageUID(uid string) (pkg, strct string, err error) {
	parts := strings.Split(uid, ".")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("Invalid struct identifier: %s", uid)
	}

	pkg = strings.Join(parts[:len(parts)-1], ".")
	strct = parts[len(parts)-1]
	return
}
