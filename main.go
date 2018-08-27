package main

import (
	"os"
)

const Name = "gommy"
const Usage = "The command create a dummy data file based on a definition file."
const Version = "0.0.1"

func main() {
	cli := &CLI{outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(cli.Run(os.Args))
}
