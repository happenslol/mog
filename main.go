package main

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/packages"
)

func main() {
	loadStruct := "github.com/happenslol/mog/fixtures.Author"
	cfg := &packages.Config{
		Mode: packages.NeedFiles | packages.NeedSyntax,
	}

	parts := strings.Split(loadStruct, ".")
	packagePath := strings.Join(parts[:len(parts)-1], ".")
	strct := parts[len(parts)-1]

	pkgs, err := packages.Load(cfg, packagePath)
	if err != nil {
		panic(err)
	}

	for _, pkg := range pkgs {
		for _, syn := range pkg.Syntax {
			for _, decl := range syn.Decls {
				d, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}

				for _, spec := range d.Specs {
					tspec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}

					if strct == tspec.Name.String() {
						fmt.Printf("found struct: %s\n", tspec.Name.String())
					}
				}
			}
		}
	}
}
