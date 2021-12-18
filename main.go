package main

import (
	"fmt"
	"log"
	"os"

	"github.com/happenslol/mog/codegen"
	"github.com/happenslol/mog/templates"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func main() {
	log.SetFlags(0)

	app := &cli.App{
		Name:   "mog",
		Action: generate,
		Usage:  "Generate mongodb collections and query helpers",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				DefaultText: "mog.yml",
				Value:       "mog.yml",
				Usage:       "Config file location",
				Aliases:     []string{"c"}}},

		Commands: []*cli.Command{
			{
				Name:   "init",
				Usage:  "Generate a file containing the default configuration",
				Action: initConfig}}}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func generate(c *cli.Context) error {
	configPath := c.String("config")
	rawConfig, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("Failed to read config: %s", err.Error())
	}

	config := codegen.DefaultConfig()
	if err := yaml.Unmarshal(rawConfig, &config); err != nil {
		return fmt.Errorf("Failed to parse config: %s", err.Error())
	}

	config.Generate()
	return nil
}

func initConfig(c *cli.Context) error {
	configPath := c.String("config")
	templates.WriteDefaultConfig(configPath)
	return nil
}
