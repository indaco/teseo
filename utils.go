package teseo

import (
	"fmt"
	"html"
	"io"
	"math/rand"
	"net/http"
	"time"
)

// GenerateUniqueKey generates a unique key using math/rand.
func GenerateUniqueKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 16)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// GetFullURL constructs the full URL from the http.Request object.
func GetFullURL(r *http.Request) string {
	// Determine the scheme. If r.TLS is non-nil, the scheme is https, otherwise, it's http.
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	// Construct the full URL using the scheme, host, and path.
	return fmt.Sprintf("%s://%s%s", scheme, r.Host, r.URL.Path)
}

// WriteMetaTag writes a single HTML meta tag to the provided writer.
func WriteMetaTag(w io.Writer, property, content string) error {
	if content == "" {
		return nil
	}
	_, err := io.WriteString(w, fmt.Sprintf(`<meta property="%s" content="%s" />`, html.EscapeString(property), html.EscapeString(content)))
	if err != nil {
		return fmt.Errorf("failed to write %s meta tag: %w", property, err)
	}
	return nil
}
