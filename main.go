package main

import (
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

	pkgs, err := packages.Load(cfg, packagePath)
	if err != nil {
		panic(err)
	}

	for _, pkg := range pkgs {
		for _, synt := range pkg.Syntax {
		}
	}
}
