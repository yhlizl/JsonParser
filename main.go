package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	defaultjson "JsonParse/default"
	jsonuse "JsonParse/json"

	"github.com/urfave/cli"
)

func defaultRun(c *cli.Context) error {
	defaultjson.DefaultJson(location)

	return nil
}

func jsonRun(c *cli.Context) error {
	jsonuse.Jsonuse(location)

	return nil
}

var location string

func main() {
	fmt.Println("-----start:----")

	//Start

	app := cli.NewApp()
	app.Name = "Json Parser"
	app.Usage = "Json Parser v1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "from, f",
			Usage:       "input the location of the files",
			Destination: &location,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "jsonRun",
			Aliases: []string{"j"},
			Usage:   "run script to output summary excel from json files",
			Action:  jsonRun,
		},
		{
			Name:    "defaultRun",
			Aliases: []string{"d"},
			Usage:   "run script to output summary excel from default files",
			Action:  defaultRun,
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}