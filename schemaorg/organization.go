package schemaorg

import (
	"fmt"
	"html/template"

	"github.com/a-h/templ"
	"github.com/indaco/teseo"
)

// Organization represents a Schema.org Organization object.
// For more details about the meaning of the properties, see:https://schema.org/Organization
//
// Example usage:
//
// Pure struct usage:
//
// 	organization := &schemaorg.Organization{
// 		Name:        "Example Organization",
// 		URL:         "https://www.example.com",
// 		Logo:        "https://www.example.com/logo.jpg",
// 		Description: "This is an example organization.",
// 	}
//
// Factory method usage:
//
// 	organization := schemaorg.NewOrganization(
// 		"Example Organization",
// 		"https://www.example.com",
// 		"https://www.example.com/logo.jpg",
// 		"This is an example organization",
// 	)
//
// // Rendering JSON-LD using templ:
//
//	templ Page() {
//		@organization.ToJsonLd()
//	}
//
// // Rendering JSON-LD as `template.HTML` value:
//
//	jsonLdHtml := organization.ToGoHTMLJsonLd()
//
// Expected output:
//
// 	{
// 		"@context": "https://schema.org",
// 		"@type": "Organization",
// 		"name": "Example Organization",
// 		"url": "https://www.example.com",
// 		"logo": "https://www.example.com/logo.jpg",
// 		"description": "This is an example organization"
// 	}

// Organization represents a Schema.org Organization object
// For more details about the meaning of the properties see: https://schema.org/Organization
type Organization struct {
	Context       string         `json:"@context"`
	Type          string         `json:"@type"`
	Name          string         `json:"name,omitempty"`
	URL           string         `json:"url,omitempty"`
	Logo          *ImageObject   `json:"logo,omitempty"`
	ContactPoints []ContactPoint `json:"contactPoint,omitempty"`
	SameAs        []string       `json:"sameAs,omitempty"`
}

// Validate checks for recommended fields in Organization.
func (org *Organization) Validate() []string {
	var warnings []string

	if org.Name == "" {
		warnings = append(warnings, "missing recommended field: name")
	}
	if org.URL == "" {
		warnings = append(warnings, "missing recommended field: url")
	}
	if org.Logo == nil || org.Logo.URL == "" {
		warnings = append(warnings, "missing recommended field: logo.url")
	}

	return warnings
}

// ToJsonLd converts the Organization struct to a JSON-LD `templ.Component`.
func (org *Organization) ToJsonLd() templ.Component {
	org.ensureDefaults()
	id := fmt.Sprintf("%s-%s", "org", teseo.GenerateUniqueKey())
	return templ.JSONScript(id, org).WithType("application/ld+json")
}

// ToGoHTMLJsonLd renders the Organization struct as `template.HTML` value for Go's `html/template`.
func (org *Organization) ToGoHTMLJsonLd() (template.HTML, error) {
	return teseo.RenderToHTML(org.ToJsonLd())
}

// NewOrganization initializes an Organization with default context and type.
func NewOrganization(name string, url string, logoURL string, contactPoints []ContactPoint, sameAs []string) *Organization {
	org := &Organization{
		Name: name,
		URL:  url,
		Logo: &ImageObject{
			Type: "ImageObject",
			URL:  logoURL,
		},
		ContactPoints: contactPoints,
		SameAs:        sameAs,
	}
	org.ensureDefaults()
	return org
}

func (org *Organization) ensureDefaults() {
	if org.Context == "" {
		org.Context = "https://schema.org"
	}

	if org.Type == "" {
		org.Type = "Organization"
	}

	if org.Logo != nil {
		org.Logo.ensureDefaults()
	}
}
