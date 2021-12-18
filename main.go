package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io"
	"os"
	"strings"

	"github.com/happenslol/mog/codegen"
	"github.com/happenslol/mog/templates"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Output      OutputConfig
	Primitives  []string
	Collections map[string]CollectionConfig
}

type OutputConfig struct {
	Package  string
	Filename string
}

type CollectionConfig struct {
	Model string
}

type TemplateInput struct {
	PackageName string
	Collections map[string]CollectionInput
	Structs     []codegen.Strct
	Imports     []string
}

type CollectionInput struct {
	ModelType      string
	CollectionType string
	IDType         string
	Imports        []string
}

var defaultPrimitives = []string{
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

	modelImports := map[string]struct{}{}

	input := &TemplateInput{
		PackageName: config.Output.Package,
		Collections: map[string]CollectionInput{},
		Structs:     []codegen.Strct{},
		Imports:     []string{}}

	// If we just mark all primitives as already
	// known, they will be skipped when we reach them
	strcts := map[string]codegen.Strct{}
	for _, s := range config.Primitives {
		strcts[s] = codegen.Strct{}
	}

	for _, s := range defaultPrimitives {
		strcts[s] = codegen.Strct{}
	}

	for colName, col := range config.Collections {
		pkg, name, err := splitPackageUID(col.Model)
		if err != nil {
			panic(err)
		}

		codegen.LoadModel(pkg, name, strcts)

		parts := strings.Split(pkg, "/")
		modelType := fmt.Sprintf("%s.%s", parts[len(parts)-1], name)

		// Find our struct from the struct list. It's guaranteed to
		// be there if the previous function did not error out.
		strct := strcts[col.Model]
		idField := findIDType(strct)

		collectionType := fmt.Sprintf("%sCollection", strings.Title(colName))

		result := CollectionInput{
			ModelType:      modelType,
			IDType:         idField.Type,
			CollectionType: collectionType,
		}

		for imp := range idField.Imports {
			result.Imports = append(result.Imports, imp)
		}

		input.Collections[colName] = result
	}

	for _, s := range strcts {
		input.Structs = append(input.Structs, s)
	}

	for imp := range modelImports {
		input.Imports = append(input.Imports, imp)
	}

	outputFile, err := os.Create(config.Output.Filename)
	if err != nil {
		panic(fmt.Sprintf("Failed to create output file: %s", err.Error()))
	}

	if err := writeFormattedOutput(input, outputFile); err != nil {
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

func writeFormattedOutput(input *TemplateInput, out io.Writer) error {
	buf := new(bytes.Buffer)
	if err := templates.Tmpl.Execute(buf, input); err != nil {
		return err
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	_, err = out.Write(formatted)
	return err
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
