package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func runCli(c *cli.Context) {
	backend := NewBackend()
	backend.Run()
	backend.Shutdown()
}

func main() {
	app := cli.NewApp()
	app.Name = "notekeeper-backend"
	app.Usage = "NoteKeeper.io Backend"
	app.Action = runCli
	app.Run(os.Args)
}
