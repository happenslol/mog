package templates

import "text/template"

var Models = template.Must(template.New(
	"mog-codegen-models").Delims("[[", "]]").Parse(modelsTmplRaw))

const modelsTmplRaw = `// Code generated by github.com/happenslol/mog, DO NOT EDIT.

package [[ .PackageName ]]

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
