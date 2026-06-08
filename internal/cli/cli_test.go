package cli

import (
	"strings"
	"testing"

	"instahost/internal/cipher"
)

func TestParseArgsDefaults(t *testing.T) {
	result := ParseArgs([]string{"page.html"})
	if result.Error != "" || result.Help {
		t.Fatalf("unexpected result: %+v", result)
	}
	if result.Options.File != "page.html" || result.Options.BaseURL != DefaultBaseURL || result.Options.Key != cipher.DefaultKey {
		t.Fatalf("options = %+v", result.Options)
	}
}

func TestParseArgsBaseURL(t *testing.T) {
	result := ParseArgs([]string{"page.html", "--base-url", "https://x.test/"})
	if result.Options.BaseURL != "https://x.test/" {
		t.Fatalf("base URL = %q", result.Options.BaseURL)
	}
}

func TestParseArgsKey(t *testing.T) {
	result := ParseArgs([]string{"page.html", "--key", "secret"})
	if result.Options.Key != "secret" {
		t.Fatalf("key = %q", result.Options.Key)
	}
}

func TestParseArgsHelp(t *testing.T) {
	result := ParseArgs([]string{"--help"})
	if !result.Help {
		t.Fatal("expected help flag")
	}
}

func TestParseArgsMissingFile(t *testing.T) {
	result := ParseArgs([]string{})
	if result.Error != "Error: exactly one file argument is required" || !result.ShowUsage {
		t.Fatalf("result = %+v", result)
	}
}

func TestParseArgsMissingBaseURLValue(t *testing.T) {
	result := ParseArgs([]string{"page.html", "--base-url"})
	if result.Error != "Error: --base-url requires a value" {
		t.Fatalf("error = %q", result.Error)
	}
}

func TestUsageContainsShare(t *testing.T) {
	if !strings.Contains(Usage(), "Usage: share") {
		t.Fatal("usage text missing expected prefix")
	}
}
