package main

import (
	"fmt"
	"github.com/ngerakines/oacp/client"
	"github.com/ngerakines/oacp/server"
	"log"
	"os"
	"sort"
	"time"

	"github.com/urfave/cli"
)

var ReleaseCode string
var GitCommit string
var BuildTime string

func main() {
	compiledAt, err := time.Parse(time.RFC822Z, BuildTime)
	if err != nil {
		compiledAt = time.Now()
	}
	if ReleaseCode == "" {
		ReleaseCode = "na"
	}
	if GitCommit == "" {
		GitCommit = "na"
	}

	app := cli.NewApp()
	app.Name = "oacp"
	app.Usage = "The oauth callback proxy application."
	app.Version = fmt.Sprintf("%s-%s", ReleaseCode, GitCommit)
	app.Compiled = compiledAt
	app.Copyright = "(c) 2020 Nick Gerakines"

	app.Commands = []cli.Command{
		server.Command,
		client.RecordCommand,
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
