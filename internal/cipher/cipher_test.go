package cipher

import (
	"bytes"
	"compress/gzip"
	"testing"
)

func TestSealAndOpen(t *testing.T) {
	input := []byte("gzip-like bytes")
	sealed := SealWithChecksum(input)
	opened, ok := OpenWithChecksum(sealed)

	if len(sealed) != len(input)+ChecksumBytes {
		t.Fatalf("sealed length = %d, want %d", len(sealed), len(input)+ChecksumBytes)
	}
	if !ok || !bytes.Equal(opened, input) {
		t.Fatalf("opened = %q, want %q", opened, input)
	}
}

func TestOpenChecksumMismatch(t *testing.T) {
	sealed := SealWithChecksum([]byte("payload"))
	sealed[0] ^= 0xff

	if _, ok := OpenWithChecksum(sealed); ok {
		t.Fatal("expected checksum mismatch")
	}
}

func TestCRC32Stable(t *testing.T) {
	if got := CRC32([]byte("test")); got != 0xd87f7e0c {
		t.Fatalf("crc32 = %#x, want %#x", got, 0xd87f7e0c)
	}
}

func TestXorBytes(t *testing.T) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	gz.Write([]byte("hello world"))
	gz.Close()
	input := buf.Bytes()

	output := XorBytes(input, DefaultKey)
	if len(output) != len(input) {
		t.Fatalf("xor length = %d, want %d", len(output), len(input))
	}

	encrypted := XorBytes([]byte("test payload"), "secret")
	decrypted := XorBytes(encrypted, "secret")
	if !bytes.Equal(decrypted, []byte("test payload")) {
		t.Fatalf("decrypted = %q", decrypted)
	}

	obfuscated := XorBytes(input, DefaultKey)
	if bytes.Equal(obfuscated, input) {
		t.Fatal("expected obfuscated output to differ from input")
	}
}
