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

// fetchGitignore fetches the gitignore file for a given programming language.
//
// lang: the programming language for which to fetch the gitignore file.
// Returns the contents of the gitignore file as a byte slice and an error if any.
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
						Action: func(ctx *cli.Context) error {
							// Titlecase and trim the string so that it matches the template file name
							language = strings.ToTitle(strings.Trim(ctx.Args().First(), " "))

							if language == "" {
								return fmt.Errorf("lang is required")
							}

							content, err := fetchGitignore(language)
							if err != nil {
								fmt.Println("Error fetching gitignore content:", err)
								return nil
							}

							if err := os.WriteFile(".gitignore", content, 0644); err != nil {
								fmt.Println("Error creating .gitignore file:", err)
								return nil
							} else {
								fmt.Println(".gitignore file created successfully!")
							}
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
