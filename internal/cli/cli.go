package cli

import (
	"fmt"
	"strings"

	"instahost/internal/cipher"
)

const DefaultBaseURL = "https://timdelange.github.io/instahost/"

func Usage() string {
	return `Usage: share <file> [--base-url <url>] [--key <passphrase>]

  Minify, compress, obfuscate, and encode an HTML file into a shareable URL.

Options:
  --base-url <url>   Base URL for the static page (default: https://timdelange.github.io/instahost/)
  --key <passphrase> XOR obfuscation key (default: built-in key)
  -h, --help         Show this help`
}

type Options struct {
	File    string
	BaseURL string
	Key     string
}

type ParseResult struct {
	Help       bool
	Error      string
	ShowUsage  bool
	Options    Options
}

func ParseArgs(argv []string) ParseResult {
	args := argv
	baseURL := DefaultBaseURL
	key := cipher.DefaultKey
	var positional []string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-h", "--help":
			return ParseResult{Help: true}
		case "--base-url":
			if i+1 >= len(args) || strings.HasPrefix(args[i+1], "-") {
				return ParseResult{Error: "Error: --base-url requires a value"}
			}
			i++
			baseURL = args[i]
		case "--key":
			if i+1 >= len(args) || strings.HasPrefix(args[i+1], "-") {
				return ParseResult{Error: "Error: --key requires a value"}
			}
			i++
			key = args[i]
		default:
			positional = append(positional, arg)
		}
	}

	if len(positional) != 1 {
		return ParseResult{
			Error:     "Error: exactly one file argument is required",
			ShowUsage: true,
		}
	}

	return ParseResult{
		Options: Options{
			File:    positional[0],
			BaseURL: baseURL,
			Key:     key,
		},
	}
}

func FormatReadError(fileName, message string) string {
	return fmt.Sprintf("Error: cannot read %s: %s", fileName, message)
}
