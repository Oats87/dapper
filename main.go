package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/ibuildthecloud/dapper/file"
)

func main() {
	exit := func(err error) {
		if err != nil {
			logrus.Fatal(err)
		}
	}

	app := cli.NewApp()
	app.Author = "@ibuildthecloud, @imikushin"
	app.EnableBashCompletion = true
	app.Usage = "Docker build wrapper"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Value: "Dockerfile.dapper",
			Usage: "Dockerfile to build from",
		},
		cli.BoolFlag{
			Name:  "socket, k",
			Usage: "Bind in the Docker socket",
		},
		cli.StringFlag{
			Name:   "mode, m",
			Value:  "auto",
			Usage:  "Execution mode for Dapper bind/cp/auto",
			EnvVar: "DAPPER_MODE",
		},
		cli.BoolFlag{
			Name:  "no-out, O",
			Usage: "Do not copy the output back (in --mode cp)",
		},
		cli.StringFlag{
			Name:  "directory, C",
			Value: ".",
			Usage: "The directory in which to run, --file is relative to this",
		},
		cli.BoolFlag{
			Name:  "shell, s",
			Usage: "Launch a shell",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "Print debugging",
		},
	}
	app.Action = func(c *cli.Context) {
		exit(run(c))
	}

	exit(app.Run(os.Args))
}

func run(c *cli.Context) error {
	if c.Bool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	dir := c.String("directory")
	mode := c.String("mode")
	shell := c.Bool("shell")
	socket := c.Bool("socket")
	noOut := c.Bool("no-out")

	if err := os.Chdir(dir); err != nil {
		return fmt.Errorf("Failed to change to directory %s: %v", dir, err)
	}

	dapperFile, err := file.Lookup(c.String("file"))
	if err != nil {
		return err
	}

	dapperFile.SetSocket(socket)
	dapperFile.SetNoOut(noOut)

	if shell {
		return dapperFile.Shell(mode, c.Args())
	}

	return dapperFile.Run(mode, c.Args())
}
