package main

import (
	"fmt"
	"os"
	"path/filepath"

	"instahost/internal/cli"
	"instahost/internal/encode"
)

func main() {
	result := cli.ParseArgs(os.Args[1:])
	if result.Help {
		fmt.Fprintln(os.Stderr, cli.Usage())
		os.Exit(0)
	}
	if result.Error != "" {
		fmt.Fprintln(os.Stderr, result.Error)
		if result.ShowUsage {
			fmt.Fprintln(os.Stderr, cli.Usage())
		}
		os.Exit(1)
	}

	filePath, err := filepath.Abs(result.Options.File)
	if err != nil {
		fmt.Fprintln(os.Stderr, cli.FormatReadError(filepath.Base(result.Options.File), err.Error()))
		os.Exit(1)
	}

	html, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, cli.FormatReadError(filepath.Base(filePath), err.Error()))
		os.Exit(1)
	}

	encoded, err := encode.EncodeHTML(string(html), result.Options.Key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	url := encode.BuildShareURL(result.Options.BaseURL, encoded.Encoded, result.Options.Key)

	fmt.Println(url)
	fmt.Fprintf(os.Stderr, "\nOriginal:  %d bytes\n", len(html))
	fmt.Fprintf(os.Stderr, "Minified:  %d bytes\n", len(encoded.Minified))
	fmt.Fprintf(os.Stderr, "Compressed: %d bytes\n", len(encoded.Compressed))
	fmt.Fprintf(os.Stderr, "URL length: %d chars\n", len(url))
}
