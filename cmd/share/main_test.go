package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"instahost/internal/encode"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func exampleFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(repoRoot(t), "example.html")
}

func runShare(t *testing.T, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	root := repoRoot(t)
	cmd := exec.Command("go", append([]string{"run", "./cmd/share"}, args...)...)
	cmd.Dir = root
	var outBuf, errBuf strings.Builder
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	exitCode = 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatal(err)
		}
	}
	return outBuf.String(), errBuf.String(), exitCode
}

func TestCLISuccess(t *testing.T) {
	stdout, stderr, code := runShare(t, exampleFile(t))
	if code != 0 {
		t.Fatalf("exit code = %d, stderr = %s", code, stderr)
	}
	url := strings.TrimSpace(stdout)
	if !regexp.MustCompile(`^index\.html\?d=[A-Za-z0-9_-]+$`).MatchString(url) {
		t.Fatalf("stdout = %q", url)
	}
	for _, line := range []string{"Original:", "Minified:", "Compressed:", "URL length:"} {
		if !strings.Contains(stderr, line) {
			t.Fatalf("stderr missing %q: %s", line, stderr)
		}
	}
}

func TestCLIBaseURL(t *testing.T) {
	stdout, _, code := runShare(t, exampleFile(t), "--base-url", "https://example.com/host/")
	if code != 0 {
		t.Fatal("expected success")
	}
	if !regexp.MustCompile(`^https://example\.com/host/\?d=[A-Za-z0-9_-]+$`).MatchString(strings.TrimSpace(stdout)) {
		t.Fatalf("stdout = %q", stdout)
	}
}

func TestCLIKeyParam(t *testing.T) {
	stdout, _, code := runShare(t, exampleFile(t), "--key", "secret")
	if code != 0 {
		t.Fatal("expected success")
	}
	if !regexp.MustCompile(`^index\.html\?d=[A-Za-z0-9_-]+&k=c2VjcmV0$`).MatchString(strings.TrimSpace(stdout)) {
		t.Fatalf("stdout = %q", stdout)
	}
}

func TestCLIRoundTrip(t *testing.T) {
	stdout, _, code := runShare(t, exampleFile(t))
	if code != 0 {
		t.Fatal("expected success")
	}
	url := strings.TrimSpace(stdout)
	payload := strings.TrimPrefix(url, "index.html?d=")
	decoded, err := encode.DecodePayload(payload, "instahost")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(decoded, "Hello from InstaHost") {
		t.Fatal("decoded HTML missing expected content")
	}
}

func TestCLIMissingFile(t *testing.T) {
	_, stderr, code := runShare(t, "does-not-exist.html")
	if code == 0 {
		t.Fatal("expected failure")
	}
	if !strings.Contains(stderr, "cannot read") {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestCLINoFileArg(t *testing.T) {
	_, stderr, code := runShare(t)
	if code == 0 {
		t.Fatal("expected failure")
	}
	if !strings.Contains(stderr, "exactly one file argument is required") {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestCLIHelp(t *testing.T) {
	_, stderr, code := runShare(t, "--help")
	if code != 0 {
		t.Fatalf("exit code = %d", code)
	}
	if !strings.Contains(stderr, "Usage: share") {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
