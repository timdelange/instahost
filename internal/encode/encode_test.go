package encode

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"instahost/internal/cipher"
)

func exampleHTML(t *testing.T) string {
	t.Helper()
	path := filepath.Join("..", "..", "example.html")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestEncodeHTML(t *testing.T) {
	html := exampleHTML(t)
	result, err := EncodeHTML(html, cipher.DefaultKey)
	if err != nil {
		t.Fatal(err)
	}

	if len(result.Minified) >= len(html) {
		t.Fatalf("minified length = %d, want < %d", len(result.Minified), len(html))
	}
	if strings.Contains(result.Minified, "comment") {
		t.Fatal("minified output should not contain comments")
	}
	if !regexpMatch(result.Encoded, `^[A-Za-z0-9_-]+$`) {
		t.Fatalf("encoded = %q, want URL-safe base64", result.Encoded)
	}
	if strings.ContainsAny(result.Encoded, "+/=") {
		t.Fatalf("encoded contains non-url-safe chars: %q", result.Encoded)
	}
	if len(result.Encoded) == 0 {
		t.Fatal("encoded payload is empty")
	}
}

func TestEncodeRoundTrip(t *testing.T) {
	html := exampleHTML(t)
	result, err := EncodeHTML(html, cipher.DefaultKey)
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := DecodePayload(result.Encoded, cipher.DefaultKey)
	if err != nil {
		t.Fatal(err)
	}
	if decoded != result.Minified {
		t.Fatal("decoded HTML does not match minified output")
	}
	if !strings.Contains(decoded, "Hello from InstaHost") {
		t.Fatal("decoded HTML missing expected content")
	}
}

func TestEncodeCustomKey(t *testing.T) {
	html := exampleHTML(t)
	result, err := EncodeHTML(html, "custom-key")
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := DecodePayload(result.Encoded, "custom-key")
	if err != nil {
		t.Fatal(err)
	}
	if decoded != result.Minified {
		t.Fatal("decoded HTML does not match minified output")
	}
}

func TestEncodeStable(t *testing.T) {
	html := exampleHTML(t)
	first, err := EncodeHTML(html, cipher.DefaultKey)
	if err != nil {
		t.Fatal(err)
	}
	second, err := EncodeHTML(html, cipher.DefaultKey)
	if err != nil {
		t.Fatal(err)
	}
	if first.Encoded != second.Encoded {
		t.Fatal("encoded output is not stable")
	}
	if first.Minified != second.Minified {
		t.Fatal("minified output is not stable")
	}
}

func TestURLLengthWarning(t *testing.T) {
	if got := URLLengthWarning(strings.Repeat("a", MaxShareURLLength)); got != "" {
		t.Fatalf("unexpected warning: %q", got)
	}
	got := URLLengthWarning(strings.Repeat("a", MaxShareURLLength+1))
	if got == "" {
		t.Fatal("expected warning for long URL")
	}
	if !strings.Contains(got, "8192") {
		t.Fatalf("warning = %q", got)
	}
}

func TestBuildShareURL(t *testing.T) {
	if got := BuildShareURL("index.html", "abc123", cipher.DefaultKey); got != "index.html#d=abc123" {
		t.Fatalf("url = %q", got)
	}
	if got := BuildShareURL("https://example.com/page?foo=bar", "abc123", cipher.DefaultKey); got != "https://example.com/page?foo=bar#d=abc123" {
		t.Fatalf("url = %q", got)
	}
	if got := BuildShareURL("index.html", "abc123", "secret"); got != "index.html#d=abc123&k=c2VjcmV0" {
		t.Fatalf("url = %q", got)
	}
}

func TestExtractPayloadFromURL(t *testing.T) {
	payload, key, err := ExtractPayloadFromURL("https://example.com/#d=abc123")
	if err != nil {
		t.Fatal(err)
	}
	if payload != "abc123" || key != cipher.DefaultKey {
		t.Fatalf("payload = %q key = %q", payload, key)
	}

	payload, key, err = ExtractPayloadFromURL("https://example.com/#d=abc123&k=c2VjcmV0")
	if err != nil {
		t.Fatal(err)
	}
	if payload != "abc123" || key != "secret" {
		t.Fatalf("payload = %q key = %q", payload, key)
	}
}

func TestMinifyMatchesNodeExample(t *testing.T) {
	want := `<!doctype html><html lang="en"><head><meta charset="UTF-8"><title>Hello InstaHost</title><style>body{font-family:Georgia,serif;max-width:40rem;margin:4rem auto;padding:0 1rem;line-height:1.6}h1{color:#2563eb}</style></head><body><h1>Hello from InstaHost</h1><p>This page was shared via a compressed, URL-encoded link.</p></body></html>`
	got, err := MinifyHTML(exampleHTML(t))
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("minified output does not match Node\nwant: %s\ngot:  %s", want, got)
	}
}

func regexpMatch(s, pattern string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_' || c == '-') {
			return false
		}
	}
	return len(s) > 0 || pattern != `^[A-Za-z0-9_-]+$`
}
