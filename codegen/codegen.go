package codegen

import (
	"fmt"
	"log"
	"strings"

	"github.com/happenslol/mog/templates"
)

var defaultPrimitives = []string{
	"time.Time",
}

var defaultModelImports = []string{
	"time",
}

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

type Models struct {
	PackageName string
	Structs     []Strct
	Imports     []string
}

type Collections struct {
	PackageName string
	Collections map[string]Collection
	Imports     []string
}

type Collection struct {
	ModelType      string
	CollectionType string
	IDType         string
}

func DefaultConfig() *Config {
	cfg := &Config{
		Output: struct {
			Models      OutputConfig
			Collections OutputConfig
		}{
			Models: OutputConfig{
				Package:  "modelgen",
				Filename: "modelgen/gen.go"},
			Collections: OutputConfig{
				Package:  "colgen",
				Filename: "colgen/gen.go"}},

		Primitives:  []string{},
		Collections: map[string]string{}}

	for _, p := range defaultPrimitives {
		cfg.Primitives = append(cfg.Primitives, p)
	}

	return cfg
}

func (c *Config) Generate() {
	collections := Collections{
		PackageName: c.Output.Collections.Package,
		Collections: make(map[string]Collection),
		Imports:     make([]string, 0),
	}

	models := Models{
		PackageName: c.Output.Models.Package,
		Structs:     make([]Strct, 0),
		Imports:     make([]string, 0),
	}

	collectionImports := map[string]struct{}{}
	for _, imp := range defaultModelImports {
		collectionImports[imp] = struct{}{}
	}

	// If we just mark all primitives as already
	// known, they will be skipped when we reach them
	strcts := map[string]Strct{}
	for _, s := range c.Primitives {
		strcts[s] = Strct{Skip: true}
	}

	for colName, col := range c.Collections {
		pkg, name, err := splitPackageUID(col)
		if err != nil {
			log.Fatalf("Invalid struct identifier: %s", err.Error())
		}

		LoadModel(pkg, name, strcts)

		parts := strings.Split(pkg, "/")
		modelType := fmt.Sprintf("%s.%s", parts[len(parts)-1], name)

		// Find our struct from the struct list. It's guaranteed to
		// be there if the previous function did not error out.
		strct := strcts[col]
		idField := findIDType(strct)
		collectionImports[strct.StructImport] = struct{}{}

		collectionType := fmt.Sprintf("%sCollection", strings.Title(colName))

		result := Collection{
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

	templates.WriteFormattedOutput(collections, templates.Collections,
		c.Output.Collections.Filename)

	templates.WriteFormattedOutput(models, templates.Models,
		c.Output.Models.Filename)
}
