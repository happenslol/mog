package templates

import "text/template"

var Config = template.Must(template.New(
	"mog-init-config").Parse(configTmplRaw))

const configTmplRaw = `output:
  models:
    package: modelgen
    filename: modelgen/gen.go
  collections:
    package: colgen
    filename: colgen/gen.go

primitives: []
collections: {}
`
