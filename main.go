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
	Structs    []string
	Primitives []string
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
	cfg := &Config{
		Structs: []string{
			"github.com/happenslol/mog/fixtures.Author",
		},
		Primitives: []string{
			".uint",
			".uintptr",
			".uint8",
			".uint16",
			".uint32",
			".uint64",

			".int",
			".int8",
			".int16",
			".int32",
			".int64",

			".float32",
			".float64",

			".complex64",
			".complex128",

			".byte",
			".rune",
			".string",
			".bool",

			"time.Time",
		},
	}

	strcts := map[string]*Strct{}
	knownImports := map[string]string{}
	usedImports := map[string]struct{}{}

	primitives := map[string]struct{}{}
	for _, p := range cfg.Primitives {
		primitives[p] = struct{}{}
	}

	for _, uid := range cfg.Structs {
		err := generateCollections(uid, strcts, primitives, knownImports, usedImports)
		if err != nil {
			panic(err)
		}
	}

	input := &StructsTemplateInput{
		PackageName: "strcts",
		Structs:     []*Strct{},
		Imports:     []string{},
	}

	for _, s := range strcts {
		input.Structs = append(input.Structs, s)
	}

	for imp := range usedImports {
		input.Imports = append(input.Imports, imp)
	}

	if err := writeFormattedOutput(input, os.Stdout); err != nil {
		panic(err)
	}
}

func generateCollections(
	uid string,
	strcts map[string]*Strct,
	primitives map[string]struct{},
	knownImports map[string]string,
	usedImports map[string]struct{},
) error {
	basePkg, _, err := splitPackageUID(uid)
	if err != nil {
		return err
	}

	structQueue := []string{uid}
	dotImports := []string{}

	for len(structQueue) != 0 {
		uid := structQueue[0]
		structQueue = structQueue[1:]

		if _, ok := strcts[uid]; ok {
			continue
		}

		if _, ok := primitives[uid]; ok {
			continue
		}

		if strings.HasPrefix(uid, ".") {
			uid = basePkg + uid
		}

		spec, foundImports, err := findTypeSpec(uid)
		if err != nil {
			return err
		}

		for _, imp := range foundImports {
			path := strings.ReplaceAll(imp.Path.Value, "\"", "")
			parts := strings.Split(path, "/")
			name := parts[len(parts)-1]

			if imp.Name != nil {
				if imp.Name.Name == "." {
					dotImports = append(dotImports, imp.Name.Name)
					continue
				}

				name = imp.Name.Name
			}

			knownImports[name] = path
		}

		strctType, ok := spec.Type.(*ast.StructType)
		if !ok {
			return fmt.Errorf("Type `%s` is not a struct", uid)
		}

		fields, typs := parseFields(strctType.Fields.List,
			knownImports, usedImports, primitives)

		strcts[uid] = &Strct{
			Name:   spec.Name.Name,
			Fields: fields}

		structQueue = append(structQueue, typs...)
	}

	return nil
}

func writeFormattedOutput(input *StructsTemplateInput, out io.Writer) error {
	buf := new(bytes.Buffer)
	if err := structsTemplate.Execute(buf, input); err != nil {
		return err
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	_, err = out.Write(formatted)
	return err
}

func parseFields(
	flds []*ast.Field,
	knownImports map[string]string,
	usedImports map[string]struct{},
	primitives map[string]struct{},
) (
	fields []Field,
	referencedTypes []string,
) {
	fields = []Field{}

	for _, fld := range flds {
		spec := fld.Type
		foundTypes := parseReferencedTypes(spec, knownImports, usedImports, primitives)
		referencedTypes = append(referencedTypes, foundTypes...)

		field := parseFieldName(fld)
		if field == nil {
			continue
		}

		fields = append(fields, *field)
	}

	return
}

func parseFieldName(fld *ast.Field) *Field {
	if fld.Names == nil {
		// TODO: Implement embedded types
		return nil
	}

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
			return nil
		}

		return &Field{
			Name:    name.Name,
			BsonKey: mongoTags.Name,
		}
	}

	return nil
}

func parseReferencedTypes(
	expr ast.Expr,
	knownImports map[string]string,
	usedImports map[string]struct{},
	primitives map[string]struct{},
) []string {
	exprQueue := []ast.Expr{expr}
	referencedTypes := []string{}

	for len(exprQueue) != 0 {
		e := exprQueue[0]
		exprQueue = exprQueue[1:]

		switch t := e.(type) {
		case *ast.Ident:
			// This is a local type or dot import
			referencedTypes = append(referencedTypes, fmt.Sprintf(".%s", t.Name))
		case *ast.SelectorExpr:
			// This is a package import
			pkgName, ok := t.X.(*ast.Ident)
			if !ok {
				panic(fmt.Sprintf("Struct package was not an identifier: %v", t.X))
			}

			typ := fmt.Sprintf("%s.%s", pkgName.Name, t.Sel.Name)
			if _, ok := primitives[typ]; ok {
				continue
			}

			imp, ok := knownImports[pkgName.Name]
			if !ok {
				panic(fmt.Sprintf(
					"Import not found: %s (for type %s)",
					t.Sel.Name, typ))
			}

			typeUID := fmt.Sprintf("%s.%s", imp, t.Sel.Name)
			if _, ok := primitives[typ]; ok {
				continue
			}

			referencedTypes = append(referencedTypes, typeUID)
			usedImports[imp] = struct{}{}
		case *ast.StarExpr:
			exprQueue = append(exprQueue, t.X)
			continue
		case *ast.ArrayType:
			exprQueue = append(exprQueue, t.Elt)
			continue
		case *ast.MapType:
			exprQueue = append(exprQueue, t.Key)
			exprQueue = append(exprQueue, t.Value)
			break
		case *ast.InterfaceType:
			// Can't generate anything for interfaces, so we just skip them.
			// TODO: Do we want to issue a warning for this?
			continue
		default:
			panic(fmt.Sprintf("found unsupported field expression: %T\n", t))
		}
	}

	return referencedTypes
}

func findTypeSpec(uid string) (*ast.TypeSpec, []*ast.ImportSpec, error) {
	pkg, strct, err := splitPackageUID(uid)
	if err != nil {
		return nil, nil, err
	}

	cfg := &packages.Config{Mode: packages.NeedFiles | packages.NeedSyntax}

	pkgs, err := packages.Load(cfg, pkg)
	if err != nil {
		return nil, nil, err
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
						return tspec, syn.Imports, nil
					}
				}
			}
		}
	}

	return nil, nil, fmt.Errorf(
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
