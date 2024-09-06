package schemaorg

import (
	"context"
	"io"
	"log"

	"github.com/a-h/templ"
	"github.com/indaco/teseo"
)

// LocalBusiness represents a Schema.org LocalBusiness object.
//
// Example usage:
//
// Pure struct usage:
//
//	localBusiness := &schemaorg.LocalBusiness{
//		Name:        "Example Business",
//		Address:     &schemaorg.PostalAddress{StreetAddress: "123 Main St", AddressLocality: "Anytown", AddressRegion: "CA", PostalCode: "12345"},
//		Telephone:   "+1-800-555-1234",
//		Description: "This is an example local business.",
//	}
//
// Factory method usage:
//
//	localBusiness := schemaorg.NewLocalBusiness(
//		"Example Business",
//		&schemaorg.PostalAddress{StreetAddress: "123 Main St", AddressLocality: "Anytown", AddressRegion: "CA", PostalCode: "12345"},
//		"+1-800-555-1234",
//		"This is an example local business",
//	)
//
// // Rendering JSON-LD using templ:
//
//	templ Page() {
//		@localBusiness.ToJsonLd()
//	}
//
// // Rendering JSON-LD as `template.HTML` value:
//
//	jsonLdHtml := localBusiness.ToGoHTMLJsonLd()
//
// Expected output:
//
//	{
//		"@context": "https://schema.org",
//		"@type": "LocalBusiness",
//		"name": "Example Business",
//		"address": {
//			"@type": "PostalAddress",
//			"streetAddress": "123 Main St",
//			"addressLocality": "Anytown",
//			"addressRegion": "CA",
//			"postalCode": "12345"
//		},
//		"telephone": "+1-800-555-1234",
//		"description": "This is an example local business"
//	}
type LocalBusiness struct {
	Context         string           `json:"@context"`
	Type            string           `json:"@type"`
	Name            string           `json:"name,omitempty"`
	Description     string           `json:"description,omitempty"`
	URL             string           `json:"url,omitempty"`
	Logo            *ImageObject     `json:"logo,omitempty"`
	Telephone       string           `json:"telephone,omitempty"`
	Address         *PostalAddress   `json:"address,omitempty"`
	OpeningHours    []string         `json:"openingHours,omitempty"`
	Geo             *GeoCoordinates  `json:"geo,omitempty"`
	AggregateRating *AggregateRating `json:"aggregateRating,omitempty"`
	Review          []*Review        `json:"review,omitempty"`
}

// GeoCoordinates represents a Schema.org GeoCoordinates object
type GeoCoordinates struct {
	Type      string  `json:"@type"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

// NewLocalBusiness initializes a LocalBusiness with default context and type.
func NewLocalBusiness(name string, description string, url string, telephone string, logo *ImageObject, address *PostalAddress, openingHours []string, geo *GeoCoordinates, aggregateRating *AggregateRating, reviews []*Review) *LocalBusiness {
	localBusiness := &LocalBusiness{
		Name:            name,
		Description:     description,
		URL:             url,
		Logo:            logo,
		Telephone:       telephone,
		Address:         address,
		OpeningHours:    openingHours,
		Geo:             geo,
		AggregateRating: aggregateRating,
		Review:          reviews,
	}
	localBusiness.ensureDefaults()
	return localBusiness
}

// ToJsonLd converts the LocalBusiness struct to a JSON-LD `templ.Component`.
func (lb *LocalBusiness) ToJsonLd() templ.Component {
	lb.ensureDefaults()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		return templ.JSONScript(teseo.GenerateUniqueKey(), lb).WithType("application/ld+json").Render(ctx, w)
	})
}

// ToGoHTMLJsonLd renders the LocalBusiness struct as a string for Go's `html/template`.
func (lb *LocalBusiness) ToGoHTMLJsonLd() (string, error) {
	lb.ensureDefaults()
	// Create the templ component.
	templComponent := lb.ToJsonLd()
	// Render the templ component to a `template.HTML` value.
	html, err := templ.ToGoHTML(context.Background(), templComponent)
	if err != nil {
		log.Fatalf("failed to convert to html: %v", err)
	}
	return string(html), nil
}

// ensureDefaults sets default values for LocalBusiness and its nested objects if they are not already set.
func (lb *LocalBusiness) ensureDefaults() {
	if lb.Context == "" {
		lb.Context = "https://schema.org"
	}
	if lb.Type == "" {
		lb.Type = "LocalBusiness"
	}

	if lb.Logo != nil {
		lb.Logo.ensureDefaults()
	}

	if lb.Address != nil {
		lb.Address.ensureDefaults()
	}

	if lb.Geo != nil {
		lb.Geo.ensureDefaults()
	}

	if lb.AggregateRating != nil {
		lb.AggregateRating.ensureDefaults()
	}

	for _, review := range lb.Review {
		review.ensureDefaults()
	}
}

// ensureDefaults sets default values for GeoCoordinates if they are not already set.
func (geo *GeoCoordinates) ensureDefaults() {
	if geo.Type == "" {
		geo.Type = "GeoCoordinates"
	}
}
