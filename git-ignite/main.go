package main

import (
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "gitignite",
		Usage: "generate .gitignore file from a template",
		Commands: []*cli.Command{
			{
				Name: "init",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "lang",
						Aliases:  []string{"l"},
						Usage:    "programming language template",
						Required: true,
						Action: func(ctx *cli.Context, s string) error {
							// Lowercase the string so that it matches the template file name
							s = strings.ToTitle(s)
							// Trim the string
							s = strings.Trim(s, " ")
							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
