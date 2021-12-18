package templates

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"text/template"
)

func WriteCollections(input CollectionsInput, filename string) error {
	return writeFormattedOutput(input, collectionsTmpl, filename)
}

func WriteModels(input ModelsInput, filename string) error {
	return writeFormattedOutput(input, modelsTmpl, filename)
}

func writeFormattedOutput(
	input interface{},
	tmpl *template.Template,
	filename string,
) error {
	out, err := os.Create(filename)
	if err != nil {
		panic(fmt.Sprintf("Failed to create output file: %s", err.Error()))
	}

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, input); err != nil {
		return err
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	_, err = out.Write(formatted)
	return err
}
