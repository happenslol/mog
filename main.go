package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"io"
	"os"
	"reflect"
	"strings"
	"text/template"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"golang.org/x/tools/go/packages"
)

const structsTemplateRaw = `package [[ .PackageName ]]

import "go.mongodb.org/mongo-driver/bson"

[[ range $strct := .Structs ]]
type _[[ $strct.Name ]] struct{}

var [[ $strct.Name ]] = new(_[[ $strct.Name ]])

[[ range $fld := $strct.Fields ]]
func (_[[ $strct.Name ]]) [[ $fld.Name ]](filter interface{}) bson.D {
	return bson.D{{
		Key: "[[ $fld.BsonKey ]]",
		Value: filter}}
}
[[ end ]]
[[ end ]]
`

type StructsTemplateInput struct {
	PackageName string
	Structs     []*Strct
}

type Strct struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name    string
	BsonKey string
	Type    string
}

var structsTemplate = template.Must(
	template.New("structs").Delims("[[", "]]").Parse(structsTemplateRaw))

func main() {
	testPath := "github.com/happenslol/mog/fixtures.Book"
	err := generateCollectionForStruct(testPath, os.Stdout)
	if err != nil {
		panic(err)
	}
}

func generateCollectionForStruct(uid string, out io.Writer) error {
	spec, err := findTypeSpec(uid)
	if err != nil {
		return err
	}

	strct, err := parseStruct(uid, spec)
	if err != nil {
		return err
	}

	input := StructsTemplateInput{
		PackageName: "strcts",
		Structs:     []*Strct{strct},
	}

	buf := new(bytes.Buffer)
	err = structsTemplate.Execute(buf, input)
	if err != nil {
		return err
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	_, err = out.Write(formatted)
	return err
}

func parseStruct(uid string, spec *ast.TypeSpec) (*Strct, error) {
	strctType, ok := spec.Type.(*ast.StructType)
	if !ok {
		return nil, fmt.Errorf("Type `%s` is not a struct", uid)
	}

	fields := parseFields(strctType.Fields.List)

	return &Strct{
		Name:   spec.Name.Name,
		Fields: fields,
	}, nil
}

func parseFields(flds []*ast.Field) []Field {
	result := []Field{}

	for _, fld := range flds {
		ide := fld.Type.(*ast.Ident)

		for _, name := range fld.Names {
			var tag reflect.StructTag
			if fld.Tag != nil {
				tag = reflect.StructTag(fld.Tag.Value)
			}

			dummyStructField := reflect.StructField{
				Name: name.Name,
				Tag:  tag,
			}

			// NOTE: This function never returns an err, so
			// we can safely ignore it.
			mongoTags, _ := bsoncodec.DefaultStructTagParser(dummyStructField)
			bsonKey := mongoTags.Name

			result = append(result, Field{
				Name:    name.Name,
				BsonKey: bsonKey,
				Type:    ide.Name})
		}
	}

	return result
}

func findTypeSpec(uid string) (*ast.TypeSpec, error) {
	pkg, strct, err := splitPackageUID(uid)
	if err != nil {
		return nil, err
	}

	cfg := &packages.Config{Mode: packages.NeedFiles | packages.NeedSyntax}

	pkgs, err := packages.Load(cfg, pkg)
	if err != nil {
		return nil, err
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
						return tspec, nil
					}
				}
			}
		}
	}

	return nil, fmt.Errorf(
		"Struct `%s` not found in module `%s`",
		strct, pkg,
	)
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
