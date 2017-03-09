package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/kardianos/service"
)

func runCli(c *cli.Context) {
	svcConfig := &service.Config{
		Name:        "NoteKeeper.io",
		DisplayName: "NoteKeeper.io",
		Description: "This is the NoteKeeper.io backend service.",
	}

	backend := NewBackend()
	backend.Logger.Debug("Starting Service...")
	svc, err := service.New(backend, svcConfig)
	if err != nil {
		backend.Logger.Fatal(err)
	}
	logger, err := svc.Logger(nil)
	if err != nil {
		backend.Logger.Fatal(err)
	}
	err = svc.Run()
	if err != nil {
		logger.Error(err)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "notekeeper-backend"
	app.Usage = "NoteKeeper.io Backend"
	app.Action = runCli
	app.Run(os.Args)
}
