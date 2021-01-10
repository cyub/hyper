package cli

import (
	"os"

	"github.com/gobuffalo/packr/v2"
	"github.com/urfave/cli/v2"
)

// default cli app
var app *App

// App define cli application struct
type App struct {
	*cli.App
	box *packr.Box
	pwd string
}

// Run use for run the cli app
func Run() error {
	app = &App{cli.NewApp(), nil, ""}
	app.Name = "Hyper"
	app.Usage = "hyper is an boilerplate for Golang Gin framework"
	app.Version = "1.0.0"
	app.HelpName = "hyper"
	app.Commands = []*cli.Command{
		newCmd,
	}
	app.Before = before
	return app.Run(os.Args)
}

func before(c *cli.Context) error {
	app.pwd, _ = os.Getwd()
	app.box = packr.New("stubs", "../resource/stubs")
	return nil
}
