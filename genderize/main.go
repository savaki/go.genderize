package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/savaki/go.genderize"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	filenameFlag = "filename"
	csvFlag      = "csv"
)

func main() {
	app := cli.NewApp()
	app.Name = "genderize"
	app.Usage = "a command line interface for genderize.io"
	app.Version = "0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{filenameFlag, "", "the optional filename to read names from (one request)", ""},
		cli.BoolFlag{csvFlag, "output data as csv rather than json", ""},
	}
	app.Action = Run
	app.Run(os.Args)
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func Run(c *cli.Context) {
	names := c.Args()

	if len(names) == 0 {
		cli.ShowAppHelp(c)
		os.Exit(1)
	}

	filename := c.String(filenameFlag)
	if filename != "" {
		data, err := ioutil.ReadFile(filename)
		names = strings.Split(strings.TrimSpace(string(data)), "\n")
		checkErr(err)
	}

	client := genderize.New()
	results, err := client.Query(names...)
	checkErr(err)

	if c.Bool(csvFlag) {
		displayAsCsv(results)
	} else {
		displayAsJson(results)
	}
}

func displayAsCsv(results genderize.Results) {
	fmt.Println("name,gender,probability,count")
	for _, result := range results {
		fmt.Printf("%s,%s,%s,%d\n", result.Name, result.Gender, result.Probability, result.Count)
	}
}

func displayAsJson(results genderize.Results) {
	json.NewEncoder(os.Stdout).Encode(results)
}
