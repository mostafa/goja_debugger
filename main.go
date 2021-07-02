package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "goja_debugger",
		Usage: "Runs or inspects a JS script with Goja",
		Commands: []*cli.Command{
			{
				Name:    "run",
				Aliases: []string{"r"},
				Usage:   "Runs a JS script with Goja",
				Action: func(c *cli.Context) error {
					return debug(false, "", c.Args().First())
				},
			},
			{
				Name:    "inspect",
				Aliases: []string{"i"},
				Usage:   "Debugs a JS script with Goja",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "liveinfo",
						Aliases: []string{"l"},
						Value:   "pc",
						Usage:   "Show program counter (pc) or line number (line) in debug prompt",
					},
				},
				Action: func(c *cli.Context) error {
					return debug(true, c.String("liveinfo"), c.Args().First())
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
