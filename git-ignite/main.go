package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

// decodeBase64Response decodes a base64-encoded response and returns the decoded contents.
//
// It takes in a byte slice containing the response and returns a byte slice with the decoded contents.
// The function returns an error if there is an issue parsing the JSON or decoding the base64 content.
func decodeBase64Response(response []byte) ([]byte, error) {
	var result struct {
		Content string `json:"content"`
	}

	err := json.Unmarshal(response, &result)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %s", err)
	}

	contents, err := base64.StdEncoding.DecodeString(result.Content)
	if err != nil {
		return nil, fmt.Errorf("error decoding base64: %s", err)
	}

	return contents, nil
}

// fetchGitignore fetches a gitignore file from the GitHub API based on the specified language.
//
// lang: The language for which the gitignore file is requested.
// Returns the contents of the gitignore file as a byte array and an error if any.
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

	contents, err := decodeBase64Response(body)

	return contents, nil
}

func main() {
	app := &cli.App{
		Name:  "gitignite",
		Usage: "generate .gitignore file from a template",
		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "generate .gitignore file from a language template",
				Action: func(ctx *cli.Context) error {
					// Titlecase and trim the string so that it matches the template file name
					language := strings.ToTitle(strings.Trim(ctx.Args().First(), " "))

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
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
