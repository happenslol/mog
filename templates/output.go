package templates

import (
	"bytes"
	"go/format"
	"log"
	"os"
	"strings"
	"text/template"
)

func WriteFormattedOutput(
	input interface{},
	tmpl *template.Template,
	filename string,
) {
	ensureContainingDir(filename)
	out, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create output file: %s", err.Error())
	}

	defer out.Close()

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, input); err != nil {
		log.Fatalf("Failed to execute template: %s", err.Error())
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatalf("Failed to format output: %s", err.Error())
	}

	if _, err = out.Write(formatted); err != nil {
		log.Fatalf("Failed to write output: %s", err.Error())
	}
}

func WriteDefaultConfig(filename string) {
	ensureContainingDir(filename)
	out, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create config file: %s", err.Error())
	}

	defer out.Close()

	if err := Config.Execute(out, struct{}{}); err != nil {
		log.Fatalf("Failed to execute config template: %s", err.Error())
	}
}

func ensureContainingDir(filename string) {
	parts := strings.Split(filename, "/")
	if len(parts) > 1 {
		pathParts := parts[:len(parts)-1]
		path := strings.Join(pathParts, "/")

		if err := os.MkdirAll(path, 0755); err != nil {
			log.Fatalf("Failed to create directory for path %s: %s",
				path, err.Error())
		}
	}

}
