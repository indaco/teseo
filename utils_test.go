package teseo

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"io"
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

func TestGetFullURL_HTTPS(t *testing.T) {
	r, _ := http.NewRequest("GET", "/secure", nil)
	r.Host = "secure.example.com"
	r.TLS = &tls.ConnectionState{} // non-nil means HTTPS

	url := GetFullURL(r)
	expected := "https://secure.example.com/secure"

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

type errorWriter struct{}

func (errorWriter) Write(p []byte) (int, error) {
	return 0, errors.New("write failure")
}

func TestWriteMetaTag_Error(t *testing.T) {
	err := WriteMetaTag(errorWriter{}, "og:title", "Test")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	expectedMsg := "failed to write og:title meta tag: write failure"
	if err.Error() != expectedMsg {
		t.Errorf("unexpected error message:\nexpected: %s\ngot:      %s", expectedMsg, err.Error())
	}
}

type dummyComponent struct{}

func (dummyComponent) Render(_ context.Context, w io.Writer) error {
	_, err := w.Write([]byte(`<meta property="og:title" content="Hello">`))
	return err
}

func TestRenderToHTML_Success(t *testing.T) {
	html, err := RenderToHTML(dummyComponent{})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	expected := `<meta property="og:title" content="Hello">`
	if string(html) != expected {
		t.Errorf("unexpected HTML output:\nexpected: %s\ngot:      %s", expected, html)
	}
}

type brokenComponent struct{}

func (brokenComponent) Render(_ context.Context, _ io.Writer) error {
	return errors.New("simulated render error")
}

func TestRenderToHTML_Error(t *testing.T) {
	html, err := RenderToHTML(brokenComponent{})
	if err == nil {
		t.Errorf("expected error from renderToHTML")
	}
	if html != "" {
		t.Errorf("expected empty HTML, got: %s", html)
	}
}
