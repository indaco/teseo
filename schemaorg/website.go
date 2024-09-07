package schemaorg

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"

	"github.com/a-h/templ"
	"github.com/indaco/teseo"
)

// WebSite represents a Schema.org WebSite object.
//
// Example usage:
//
// Pure struct usage:
//
// 	website := &schemaorg.WebSite{
// 		URL:         "https://www.example.com",
// 		Name:        "Example Website",
// 		Description: "This is an example website.",
// 		AlternateName: "Example Site",
// 	}
//
// Factory method usage:
//
// 	website := schemaorg.NewWebSite(
// 		"https://www.example.com",
// 		"Example Website",
// 		"Example Site",
// 		"This is an example website",
// 	)
//
// // Rendering JSON-LD using templ:
//
//	templ Page() {
//		@website.ToJsonLd()
//	}
//
// // Rendering JSON-LD as `template.HTML` value:
//
//	jsonLdHtml := website.ToGoHTMLJsonLd()
//
// Expected output:
//
// 	{
// 		"@context": "https://schema.org",
// 		"@type": "WebSite",
// 		"url": "https://www.example.com",
// 		"name": "Example Website",
// 		"alternateName": "Example Site",
// 		"description": "This is an example website"
// 	}

// Target represents the target of an action in Schema.org
type Target struct {
	Type        string `json:"@type"`
	URLTemplate string `json:"urlTemplate"`
}

// Action represents a Schema.org Action object
type Action struct {
	Type       string  `json:"@type"`
	Target     *Target `json:"target"`
	QueryInput string  `json:"query-input"`
}

// WebSite represents a Schema.org WebSite object
type WebSite struct {
	Context         string  `json:"@context"`
	Type            string  `json:"@type"`
	URL             string  `json:"url,omitempty"`
	Name            string  `json:"name,omitempty"`
	AlternateName   string  `json:"alternateName,omitempty"`
	Description     string  `json:"description,omitempty"`
	PotentialAction *Action `json:"potentialAction,omitempty"`
}

func NewWebSite(url string, name string, alternateName string, description string, potentialAction *Action) *WebSite {
	website := &WebSite{
		URL:             url,
		Name:            name,
		AlternateName:   alternateName,
		Description:     description,
		PotentialAction: potentialAction,
	}
	website.ensureDefaults()
	return website
}

// ToJsonLd converts the WebSite struct to a JSON-LD `templ.Component`.
func (ws *WebSite) ToJsonLd() templ.Component {
	ws.ensureDefaults()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		id := fmt.Sprintf("%s-%s", "website", teseo.GenerateUniqueKey())
		return templ.JSONScript(id, ws).WithType("application/ld+json").Render(ctx, w)
	})
}

// ToGoHTMLJsonLd renders the WebSite struct as `template.HTML` value for Go's `html/template`.
func (ws *WebSite) ToGoHTMLJsonLd() (template.HTML, error) {
	// Create the templ component.
	templComponent := ws.ToJsonLd()

	// Render the templ component to a `template.HTML` value.
	html, err := templ.ToGoHTML(context.Background(), templComponent)
	if err != nil {
		log.Fatalf("failed to convert to html: %v", err)
	}

	return html, nil
}

func (ws *WebSite) ensureDefaults() {
	if ws.Context == "" {
		ws.Context = "https://schema.org"
	}

	if ws.Context == "" {
		ws.Type = "WebSite"
	}

	if ws.PotentialAction != nil {
		ws.PotentialAction.ensureDefaults()
	}
}

func (act *Action) ensureDefaults() {
	if act.Type == "" {
		act.Type = "Action"
	}

	if act.Target != nil {
		act.Target.ensureDefaults()
	}
}

func (tgt *Target) ensureDefaults() {
	if tgt.Type == "" {
		tgt.Type = "EntryPoint"
	}
}
