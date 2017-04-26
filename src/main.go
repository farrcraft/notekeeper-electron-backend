package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func runCli(c *cli.Context) {
	backend := NewBackend()
	backend.Logger.Debug("Starting Service...")
	backend.Run()
}

func main() {
	app := cli.NewApp()
	app.Name = "notekeeper"
	app.Usage = "NoteKeeper.io"
	app.Action = runCli
	app.Run(os.Args)
}
