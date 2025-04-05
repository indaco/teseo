package teseo

import (
	"context"
	"fmt"
	"html"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/a-h/templ"
)

var (
	charset      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand   = rand.New(rand.NewSource(time.Now().UnixNano()))
	seededRandMu sync.Mutex
)

// GenerateUniqueKey generates a unique key using math/rand.
func GenerateUniqueKey() string {
	const length = 16
	b := make([]byte, length)
	seededRandMu.Lock()
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	seededRandMu.Unlock()
	return string(b)
}

// GetFullURL constructs the full URL from the http.Request object.
func GetFullURL(r *http.Request) string {
	// Determine the scheme. If r.TLS is non-nil, the scheme is https, otherwise, it's http.
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	// Construct the full URL using the scheme, host, and path including query string.
	return fmt.Sprintf("%s://%s%s", scheme, r.Host, r.URL.RequestURI())
}

// WriteMetaTag writes a single HTML meta tag to the provided writer.
func WriteMetaTag(w io.Writer, property, content string) error {
	if content == "" {
		return nil
	}
	_, err := io.WriteString(w, fmt.Sprintf(`<meta property="%s" content="%s" >`, html.EscapeString(property), html.EscapeString(content)))
	if err != nil {
		return fmt.Errorf("failed to write %s meta tag: %w", property, err)
	}
	return nil
}

func RenderToHTML(c templ.Component) (template.HTML, error) {
	html, err := templ.ToGoHTML(context.Background(), c)
	if err != nil {
		log.Printf("failed to convert to html: %v", err)
		return "", err
	}
	return html, nil
}
