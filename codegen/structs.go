package codegen

import (
	"fmt"
	"go/types"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"golang.org/x/tools/go/loader"
)

type Strct struct {
	Name            string
	Fields          []Field
	StructImport    string
	EmbeddedImports map[string]struct{}
	Skip            bool
}

type Field struct {
	Name    string
	BsonKey string
	Type    string
	Imports map[string]struct{}
}

func LoadModel(
	pkg, strct string,
	strcts map[string]Strct,
) {
	obj, err := findTypeDef(pkg, strct)
	if err != nil {
		panic(err)
	}

	strctNamed, ok := obj.Type().(*types.Named)
	if !ok {
		panic(fmt.Sprintf("Found %T instead of a named type for %s",
			obj.Type(), strct))
	}

	strctQueue := []*types.Named{strctNamed}
	for len(strctQueue) != 0 {
		strct := strctQueue[0]
		strctQueue = strctQueue[1:]

		qual := qualifiedName(strct)
		if _, ok := strcts[qual]; ok {
			continue
		}

		fields := []Field{}
		embeddedImports := map[string]struct{}{}

		strctType, ok := strct.Underlying().(*types.Struct)
		if !ok {
			panic(fmt.Sprintf("Underlying type for %s was %T instead of struct",
				strct.String(), strct.Underlying()))
		}

		fldQueue := []fieldWithTag{}
		for i := 0; i < strctType.NumFields(); i++ {
			fldQueue = append(fldQueue, fieldWithTag{
				strctType.Field(i), strctType.Tag(i)})
		}

		for len(fldQueue) != 0 {
			fwt := fldQueue[0]
			fldQueue = fldQueue[1:]

			parsed, embedded, imports, nested := parseField(fwt)
			if parsed == nil {
				continue
			}

			fields = append(fields, *parsed)
			fldQueue = append(fldQueue, embedded...)
			strctQueue = append(strctQueue, nested...)

			for imp := range imports {
				embeddedImports[imp] = struct{}{}
			}
		}

		strcts[qual] = Strct{
			Name:            strct.Obj().Name(),
			Fields:          fields,
			StructImport:    strct.Obj().Pkg().Path(),
			EmbeddedImports: embeddedImports,
		}
	}
}

type fieldWithTag struct {
	fld *types.Var
	tag string
}

func qualifiedName(t *types.Named) string {
	return fmt.Sprintf("%s.%s", t.Obj().Pkg().Path(), t.Obj().Name())
}

func parseField(fwt fieldWithTag) (
	result *Field,
	embedded []fieldWithTag,
	embeddedImports map[string]struct{},
	nested []*types.Named,
) {
	if !fwt.fld.Exported() {
		return nil, nil, nil, nil
	}

	name := fwt.fld.Name()
	bsonKey, skip := parseBSONKey(name, fwt.tag)
	if skip {
		return nil, nil, nil, nil
	}

	imports := map[string]struct{}{}
	embeddedImports = map[string]struct{}{}

	// To track embedded fields, we need to somehow
	// track if we've gone down into a map or slice
	// type. This doesn't seem like an optimal solution
	// but it works.
	isDirect := true

	// try to find sub structs that we also need to parse
	typeQueue := []types.Type{fwt.fld.Type()}
	for len(typeQueue) != 0 {
		typ := typeQueue[0]
		typeQueue = typeQueue[1:]

		switch t := typ.(type) {
		// We always expect this to be a struct, and we don't expect
		// to ever encounter unnamed structs
		case *types.Named:
			// if there are embedded fields, we need to
			// pass them back and parse them as if they
			// are fields of the parent, no matter if we
			// already know this struct type
			if isDirect && fwt.fld.Embedded() {
				embeddedImports[t.Obj().Pkg().Path()] = struct{}{}

				strct, ok := t.Underlying().(*types.Struct)
				if !ok {
					panic(fmt.Sprintf("Underlying type for %s was %T instead of struct",
						t.String(), t.Underlying()))
				}

				for i := 0; i < strct.NumFields(); i++ {
					embedded = append(embedded, fieldWithTag{
						strct.Field(i), strct.Tag(i)})
				}

				continue
			}

			imports[t.Obj().Pkg().Path()] = struct{}{}
			nested = append(nested, t)
			continue

		// Reached a basic type, stop recursing
		case *types.Basic:
			continue

		// These just mean we need to recurse deeper
		case *types.Pointer:
			typeQueue = append(typeQueue, t.Elem())
			continue

		// These mean we're not dealing with an embedded
		// struct anymore
		case *types.Slice:
			isDirect = false
			typeQueue = append(typeQueue, t.Elem())
			continue
		case *types.Array:
			isDirect = false
			typeQueue = append(typeQueue, t.Elem())
			continue
		case *types.Map:
			isDirect = false
			typeQueue = append(typeQueue, t.Key(), t.Elem())
			continue

		default:
			// TODO: better error reporting for this
			panic(fmt.Sprintf("Unsupported model type: %T", t))
		}
	}

	result = &Field{
		name, bsonKey,
		fwt.fld.Type().String(),
		imports,
	}
	return
}

func parseBSONKey(name, tag string) (key string, skip bool) {
	strctTag := reflect.StructTag(tag)

	dummyStructField := reflect.StructField{
		Name: name,
		Tag:  strctTag,
	}

	// NOTE: This function never returns an err, so
	// we can safely ignore it.
	mongoTags, _ := bsoncodec.DefaultStructTagParser(dummyStructField)
	if mongoTags.Skip {
		return "", true
	}

	return mongoTags.Name, false
}

func findTypeDef(pkg string, strct string) (types.Object, error) {
	ld := loader.Config{}
	ld.Import(pkg)
	lprog, err := ld.Load()
	if err != nil {
		return nil, fmt.Errorf("Failed to load package: %w", err)
	}

	cpkg := lprog.Package(pkg).Pkg
	if t := cpkg.Scope().Lookup(strct); t != nil {
		return t, nil
	}

	return nil, fmt.Errorf("struct %s not found in package %s", strct, pkg)
}
