package main

import (
	"fmt"
	"io"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK = iota
	ExitCodeError
)

// CLI is the command line object
type CLI struct {
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments.
func (c *CLI) Run(args []string) int {

	app := cli.NewApp()
	app.Name = Name
	app.Usage = Usage
	app.Version = Version

	// flags
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "./config.toml",
			Usage: "a config file path",
		},
		cli.StringFlag{
			Name:  "output, o",
			Value: "./dummyData.csv",
			Usage: "a creating dummy data file path",
		},
	}

	// action
	app.Action = func(c *cli.Context) error {
		// parameter check
		param := make(map[string]string)
		param["config"] = c.String("config")
		param["output"] = c.String("output")

		log.Printf("\"config\": %#v\n", param["config"])
		log.Printf("\"output\": %#v\n", param["output"])

		if !Exists(param["config"]) {
			log.Fatalf(color.RedString("%s not find"), param["config"])
		}
		if Exists(param["output"]) {
			if Question(fmt.Sprintf("%s already exists. if you wish to overwrite the file, press enter.[Y/n]", param["output"])) {
				log.Printf(color.GreenString("Overwrite the %s"), param["output"])
			} else {
				log.Fatalf(color.RedString("Rename or delete the %s"), param["output"])
			}
		}

		// config parser
		var config Config
		_, err := toml.DecodeFile(param["config"], &config)
		FailOnError(err)

		// for k, v := range config.Header {
		// 	fmt.Printf("header row%#v %#v\n", k, v.Row)
		// }

		// main function for creating dummy data file
		Generator(param["output"], config)

		return nil
	}

	if app.Run(args) == nil {
		log.Printf("done!")
		return ExitCodeOK
	}
	return ExitCodeError
}
