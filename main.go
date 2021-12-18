package main

import (
	"flag"
	"log"
	"os"

	"github.com/happenslol/mog/codegen"
	"gopkg.in/yaml.v3"
)

func main() {
	log.SetFlags(0)

	configPath := ""
	flag.StringVar(&configPath, "c", "mog.yml", "set config file location")
	flag.Parse()

	rawConfig, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config: %s", err.Error())
	}

	config := codegen.DefaultConfig()
	if err := yaml.Unmarshal(rawConfig, &config); err != nil {
		log.Fatalf("Failed to parse config: %s", err.Error())
	}

	config.Generate()
}
