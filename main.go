package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

const RepositoryURL = "https://api.github.com/repos/github/gitignore"

// decodeBase64Response decodes a base64-encoded response and returns the decoded contents.
//
// It takes in a string containing the base64-encoded response and returns a byte slice with the decoded contents.
// The function returns an error if there is an issue decoding the base64 content or if the language is not found.
func decodeBase64Response(content string) ([]byte, error) {
	if content == "" {
		return nil, fmt.Errorf("language not found")
	}

	decodedContent, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return nil, fmt.Errorf("error decoding base64: %s", err)
	}

	return decodedContent, nil
}

func decodeJSONResponse(response []byte) (string, error) {
	var data struct {
		Content string `json:"content"`
	}

	if len(response) == 0 {
		return "", errors.New("empty response")
	}

	if err := json.Unmarshal(response, &data); err != nil {
		return "", fmt.Errorf("error decoding JSON: %w", err)
	}

	return data.Content, nil
}

func decodeJSONResponseArray(response []byte) ([]string, error) {
	var results []struct {
		Name string `json:"name"`
	}

	if err := json.Unmarshal(response, &results); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %s", err)
	}

	decodedResponses := make([]string, len(results))
	for i, result := range results {
		decodedResponses[i] = result.Name
	}

	return decodedResponses, nil
}

// addCredits adds credits to the given contents.
//
// It takes a byte slice as input parameter.
// It returns a byte slice.
func addCredits(contents []byte, noCredits bool) []byte {
	if noCredits {
		return contents
	}
	return append([]byte("# Generated by gitignite\n# Template: https://github.com/github/gitignore\n\n"), contents...)
}

// fetchGitignore fetches the .gitignore file for a given language.
//
// It takes a string parameter `lang` which specifies the language for which the .gitignore file should be fetched.
// The function returns a byte slice containing the contents of the .gitignore file for the specified language, along with an error if any.
func fetchGitignore(lang string) ([]byte, error) {
	languages, err := fetchLanguages()
	if err != nil {
		return nil, err
	}

	for _, language := range languages {
		if strings.HasSuffix(language, ".gitignore") {
			language = strings.TrimSuffix(language, ".gitignore")
			if strings.EqualFold(strings.ToLower(language), strings.ToLower(lang)) {
				url := fmt.Sprintf("%s/contents/%s.gitignore", RepositoryURL, language)
				response, err := http.Get(url)
				if err != nil {
					return nil, err
				}
				defer response.Body.Close()

				body, err := io.ReadAll(response.Body)
				if err != nil {
					return nil, err
				}

				jsonContents, err := decodeJSONResponse(body)
				contents, err := decodeBase64Response(jsonContents)
				if err != nil {
					return nil, err
				}

				return contents, nil
			}
		}
	}

	return nil, fmt.Errorf("language not found")
}

// Fetch available languages.
func fetchLanguages() ([]string, error) {
	url := fmt.Sprintf("%s/contents", RepositoryURL)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	contents, err := decodeJSONResponseArray(body)
	if err != nil {
		return nil, err
	}

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
				Action:  initCommand,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "no-credits",
						Aliases: []string{"nc"},
						Usage:   "Do not add credits to the generated .gitignore file",
						Value:   false,
					},
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "Specify the output directory",
					},
				},
			},
			{
				Name:    "langs",
				Aliases: []string{"l"},
				Usage:   "list available languages",
				Action:  listLanguages,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func initCommand(ctx *cli.Context) error {
	// Get the language argument from the command line
	language := strings.TrimSpace(ctx.Args().First())
	if language == "" {
		return errors.New("language is required")
	}

	// Check if the no-credits flag is set
	noCredits := ctx.Bool("no-credits")

	// output directory
	outputDir := ctx.String("output")
	if outputDir == "" {
		outputDir = "."
	}

	// Fetch the gitignore content for the specified language
	content, err := fetchGitignore(language)
	if err != nil {
		return fmt.Errorf("error fetching gitignore content: %w", err)
	}

	// Add credits to the gitignore content if needed
	content = addCredits(content, noCredits)

	// Write the gitignore content to a file
	outputFile := filepath.Join(outputDir, ".gitignore")
	err = os.WriteFile(outputFile, content, 0644)
	if err != nil {
		return fmt.Errorf("error creating .gitignore file: %w", err)
	}

	fmt.Println(".gitignore file created successfully!")
	return nil
}

func listLanguages(ctx *cli.Context) error {
	languages, err := fetchLanguages()
	if err != nil {
		return fmt.Errorf("error fetching languages: %w", err)
	}

	fmt.Println("Available languages:")
	for _, language := range languages {
		if strings.HasSuffix(language, ".gitignore") {
			language = strings.TrimSuffix(language, ".gitignore")
			fmt.Println(language)
		}
	}

	return nil
}
