# Git Ignite

Git Ignite is a command-line tool for generating `.gitignore` files from language templates provided by [github/gitignore](https://github.com/github/gitignore). It simplifies the process of setting up the `.gitignore` file for your project by automatically fetching the appropriate template for the selected programming language.

## Installation

You can install Git Ignite using the following command:

```sh
go get github.com/urfave/cli/v2
go install github.com/xyugen/gitignite
```

## Usage

Generate a `.gitignore` file for a specific programming language using the following command:
```sh
gitignite init <language>
```

Replace <language> with the programming language for which you want to generate the .gitignore file.

For example, to generate a `.gitignore` file for Python, you would run:
```sh
gitignite init python
```

You can also use the `--no-credits` flag to exclude credits from the generated `.gitignore` file.

## Credits

This tool is inspired by and credits [github/gitignore](https://github.com/github/gitignore) for providing the language templates.

## License

This project is licensed under [MIT License](LICENSE)