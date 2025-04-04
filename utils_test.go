package teseo

import (
	"bytes"
	"net/http"
	"testing"
)

func TestGenerateUniqueKey(t *testing.T) {
	key1 := GenerateUniqueKey()
	key2 := GenerateUniqueKey()

	if len(key1) != 16 || len(key2) != 16 {
		t.Errorf("expected length 16, got %d and %d", len(key1), len(key2))
	}

	if key1 == key2 {
		t.Errorf("expected unique keys, got identical: %s", key1)
	}
}

func TestGetFullURL(t *testing.T) {
	r, _ := http.NewRequest("GET", "/path?foo=bar", nil)
	r.Host = "example.com"
	url := GetFullURL(r)

	expected := "http://example.com/path?foo=bar"
	if url != expected {
		t.Errorf("expected %s, got %s", expected, url)
	}
}

func TestWriteMetaTag(t *testing.T) {
	var buf bytes.Buffer
	err := WriteMetaTag(&buf, "og:title", "Hello & welcome")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	expected := `<meta property="og:title" content="Hello &amp; welcome" >`
	if output != expected {
		t.Errorf("expected %s, got %s", expected, output)
	}
}

func TestWriteMetaTagEmptyContent(t *testing.T) {
	var buf bytes.Buffer
	err := WriteMetaTag(&buf, "og:title", "")
	if err != nil {
		t.Errorf("expected no error for empty content, got %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for empty content, got: %s", buf.String())
	}
}
