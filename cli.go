package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// Config for toml
type Config struct {
	Datetime DatetimeConfig
	Data     DataConfig
	Header   []HeaderConfig
}

// DatetimeConfig for toml
type DatetimeConfig struct {
	DatetimeFormat string           `toml:"datetimeFormat"`
	Start          time.Time        `toml:"start"`
	End            time.Time        `toml:"end"`
	Column         int              `toml:"column"`
	Sampling       SamplingDatetime `toml:"sampling"`
}

// SamplingDatetime for toml
type SamplingDatetime struct {
	Num  int    `toml:"num"`
	Unit string `toml:"unit"`
}

// DataConfig for toml
type DataConfig struct {
	Min       float64        `toml:"min"`
	Max       float64        `toml:"max"`
	Pointtype string         `toml:"pointtype"`
	Abnormal  []AbnormalData `toml:"abnormal"`
	Tag       []TagData      `toml:"tag"`
	Datetime  []DatetimeData `toml:"datetime"`
}

// AbnormalData for toml
type AbnormalData struct {
	Min        float64                `toml:"min"`
	Max        float64                `toml:"max"`
	Pointtype  string                 `toml:"pointtype"`
	Column     int                    `toml:"column"`
	Start      time.Time              `toml:"start"`
	End        time.Time              `toml:"end"`
	Transition TransitionAbnormalData `toml:"transition"`
}

// TransitionAbnormalData for toml
type TransitionAbnormalData struct {
	Num  int    `toml:"num"`
	Unit string `toml:"unit"`
}

// TagData for toml
type TagData struct {
	Column int      `toml:"column"`
	Rate   int      `toml:"rate"` // TODO only rate=100 now
	Value  []string `toml:"value"`
}

// DatetimeData for toml
type DatetimeData struct {
	Column         int    `toml:"column"`
	DatetimeFormat string `toml:"datetimeFormat"`
	Add            int    `toml:"add"` // [ms]
}

// HeaderConfig for toml
type HeaderConfig struct {
	Row []string `toml:"row"`
}

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
		failOnError(err)

		// for k, v := range config.Header {
		// 	fmt.Printf("header row%#v %#v\n", k, v.Row)
		// }

		// main function for creating dummy data file
		Generator(param["output"], config)

		return nil
	}

	if app.Run(args) == nil {
		return ExitCodeOK
	}
	return ExitCodeError
}
