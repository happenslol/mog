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

type Config struct {
	Structs []string
}

type StructsTemplateInput struct {
	PackageName string
	Structs     []*Strct
	Imports     []string
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
	testPath := "github.com/happenslol/mog/fixtures.Author"
	err := generateCollections([]string{testPath}, os.Stdout)
	if err != nil {
		panic(err)
	}
}

func generateCollections(uids []string, out io.Writer) error {
	input := StructsTemplateInput{
		PackageName: "strcts",
		Structs:     []*Strct{},
		Imports:     []string{},
	}

	imports := map[string]struct{}{}

	for _, uid := range uids {
		spec, err := findTypeSpec(uid)
		if err != nil {
			return err
		}

		strctType, ok := spec.Type.(*ast.StructType)
		if !ok {
			return fmt.Errorf("Type `%s` is not a struct", uid)
		}

		fields := parseFields(strctType.Fields.List)
		input.Structs = append(input.Structs, &Strct{
			Name:   spec.Name.Name,
			Fields: fields})
	}

	buf := new(bytes.Buffer)
	if err := structsTemplate.Execute(buf, input); err != nil {
		return err
	}

	for imp := range imports {
		input.Imports = append(input.Imports, imp)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	_, err = out.Write(formatted)
	return err
}

func parseFields(flds []*ast.Field) []Field {
	result := []Field{}

	for _, fld := range flds {
		for _, name := range fld.Names {
			var tag reflect.StructTag
			if fld.Tag != nil {
				unquoted := strings.ReplaceAll(fld.Tag.Value, "`", "")
				tag = reflect.StructTag(unquoted)
			}

			dummyStructField := reflect.StructField{
				Name: name.Name,
				Tag:  tag,
			}

			// NOTE: This function never returns an err, so
			// we can safely ignore it.
			mongoTags, _ := bsoncodec.DefaultStructTagParser(dummyStructField)
			if mongoTags.Skip {
				continue
			}

			bsonKey := mongoTags.Name

			result = append(result, Field{
				Name:    name.Name,
				BsonKey: bsonKey})
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
