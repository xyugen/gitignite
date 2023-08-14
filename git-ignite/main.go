package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func fetchGitignore(lang string) ([]byte, error) {
	url := fmt.Sprintf("https://api.github.com/repos/github/gitignore/contents/%s.gitignore", lang)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

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
							// Titlecase and trim the string so that it matches the template file name
							s = strings.ToTitle(strings.Trim(s, " "))
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
