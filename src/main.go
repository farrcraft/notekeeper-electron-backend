package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func RunCli(c *cli.Context) error {
	backend := NewBackend()
	backend.Run()
	backend.Shutdown()

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "notekeeper-backend"
	app.Usage = "NoteKeeper.io Backend"
	app.Action = RunCli
	app.Run(os.Args)
}
