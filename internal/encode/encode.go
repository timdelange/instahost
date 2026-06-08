package encode

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"instahost/internal/cipher"
)

const MaxShareURLLength = 8192

type Result struct {
	Minified   string
	Compressed []byte
	Encoded    string
	Key        string
}

func MinifyHTML(htmlInput string) (string, error) {
	m := minify.New()
	htmlMin := &html.Minifier{
		KeepDocumentTags: true,
		KeepEndTags:      true,
		KeepQuotes:       true,
	}
	m.Add("text/html", htmlMin)
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("application/javascript", js.Minify)

	var buf bytes.Buffer
	if err := m.Minify("text/html", &buf, strings.NewReader(htmlInput)); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func EncodeHTML(htmlInput string, key string) (Result, error) {
	if key == "" {
		key = cipher.DefaultKey
	}

	minified, err := MinifyHTML(htmlInput)
	if err != nil {
		return Result{}, err
	}

	var compressed bytes.Buffer
	gz := gzip.NewWriter(&compressed)
	if _, err := gz.Write([]byte(minified)); err != nil {
		return Result{}, err
	}
	if err := gz.Close(); err != nil {
		return Result{}, err
	}
	compressedBytes := compressed.Bytes()

	sealed := cipher.SealWithChecksum(compressedBytes)
	obfuscated := cipher.XorBytes(sealed, key)
	encoded := base64.RawURLEncoding.EncodeToString(obfuscated)

	return Result{
		Minified:   minified,
		Compressed: compressedBytes,
		Encoded:    encoded,
		Key:        key,
	}, nil
}

func BuildShareURL(baseURL, encoded, key string) string {
	if key == "" {
		key = cipher.DefaultKey
	}

	shareURL := baseURL + "#d=" + encoded
	if key != cipher.DefaultKey {
		shareURL += "&k=" + base64.RawURLEncoding.EncodeToString([]byte(key))
	}

	return shareURL
}

func URLLengthWarning(url string) string {
	if len(url) <= MaxShareURLLength {
		return ""
	}
	return fmt.Sprintf(
		"\n*** WARNING: URL length is %d chars and exceeds %d characters. ***\n*** Some apps (including Slack) may truncate or refuse to send this link. ***\n\n",
		len(url),
		MaxShareURLLength,
	)
}

func ExtractPayloadFromURL(rawURL string) (payload string, key string, err error) {
	hashIndex := strings.Index(rawURL, "#")
	if hashIndex < 0 {
		return "", "", fmt.Errorf("missing URL fragment")
	}

	values, err := url.ParseQuery(rawURL[hashIndex+1:])
	if err != nil {
		return "", "", err
	}

	payload = values.Get("d")
	if payload == "" {
		return "", "", fmt.Errorf("missing d in URL fragment")
	}

	key = cipher.DefaultKey
	if keyParam := values.Get("k"); keyParam != "" {
		decoded, err := base64.RawURLEncoding.DecodeString(keyParam)
		if err != nil {
			return "", "", err
		}
		key = string(decoded)
	}

	return payload, key, nil
}

func DecodePayload(encoded, key string) (string, error) {
	if key == "" {
		key = cipher.DefaultKey
	}

	bytes, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	deobfuscated := cipher.XorBytes(bytes, key)
	if payload, ok := cipher.OpenWithChecksum(deobfuscated); ok {
		return gunzipString(payload)
	}

	if decoded, err := gunzipString(deobfuscated); err == nil {
		return decoded, nil
	}

	return gunzipString(bytes)
}

func gunzipString(data []byte) (string, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer reader.Close()

	out, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
