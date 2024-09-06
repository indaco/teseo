package opengraph

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"

	"github.com/a-h/templ"
	"github.com/indaco/teseo"
)

// Place represents the Open Graph place metadata.
//
// Example usage:
//
// Pure struct usage:
//
//	// Create a place using pure struct
//	place := &opengraph.Place{
//		OpenGraphObject: opengraph.OpenGraphObject{
//			Title:       "Example Place",
//			URL:         "https://www.example.com/place/example-place",
//			Description: "This is an example place description.",
//			Image:       "https://www.example.com/images/place.jpg",
//		},
//		Latitude:  40.7128, // Example latitude
//		Longitude: -74.0060, // Example longitude
//		StreetAddress: "123 Main St",
//		Locality: "New York",
//		Region: "NY",
//		PostalCode: "10001",
//		Country: "USA",
//	}
//
// Factory method usage:
//
//	// Create a place using the factory method
//	place := opengraph.NewPlace(
//		"Example Place",
//		"https://www.example.com/place/example-place",
//		"This is an example place description.",
//		"https://www.example.com/images/place.jpg",
//		40.7128,   // Latitude
//		-74.0060,  // Longitude
//		"123 Main St",
//		"New York",
//		"NY",
//		"10001",
//		"USA",
//	)
//
// // Rendering the HTML meta tags using templ:
//
//	templ Page() {
//		@place.ToMetaTgs()
//	}
//
// // Rendering the HTML meta tags as `template.HTML` value:
//
//	metaTagsHtml := place.ToGoHTMLMetaTgs()
//
// Expected output:
//
//	<meta property="og:type" content="place"/>
//	<meta property="og:title" content="Example Place"/>
//	<meta property="og:url" content="https://www.example.com/place/example-place"/>
//	<meta property="og:description" content="This is an example place description."/>
//	<meta property="og:image" content="https://www.example.com/images/place.jpg"/>
//	<meta property="place:location:latitude" content="40.7128"/>
//	<meta property="place:location:longitude" content="-74.0060"/>
//	<meta property="place:contact_data:street_address" content="123 Main St"/>
//	<meta property="place:contact_data:locality" content="New York"/>
//	<meta property="place:contact_data:region" content="NY"/>
//	<meta property="place:contact_data:postal_code" content="10001"/>
//	<meta property="place:contact_data:country_name" content="USA"/>
type Place struct {
	OpenGraphObject
	Latitude      float64 // place:location:latitude, latitude of the place
	Longitude     float64 // place:location:longitude, longitude of the place
	StreetAddress string  // place:contact_data:street_address, street address of the place
	Locality      string  // place:contact_data:locality, locality or city of the place
	Region        string  // place:contact_data:region, region or state of the place
	PostalCode    string  // place:contact_data:postal_code, postal code of the place
	Country       string  // place:contact_data:country_name, country of the place
}

// NewPlace initializes a Place with the default type "place".
func NewPlace(title, url, description, image string, latitude, longitude float64, streetAddress, locality, region, postalCode, country string) *Place {
	place := &Place{
		OpenGraphObject: OpenGraphObject{
			Title:       title,
			URL:         url,
			Description: description,
			Image:       image,
		},
		Latitude:      latitude,
		Longitude:     longitude,
		StreetAddress: streetAddress,
		Locality:      locality,
		Region:        region,
		PostalCode:    postalCode,
		Country:       country,
	}
	place.ensureDefaults()
	return place
}

// ToMetaTags generates the HTML meta tags for the Open Graph Place as templ.Component.
func (place *Place) ToMetaTags() templ.Component {
	place.ensureDefaults()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, tag := range place.metaTags() {
			if tag.content != "" {
				if err := teseo.WriteMetaTag(w, tag.property, tag.content); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// ToGoHTMLMetaTags generates the HTML meta tags for the Open Graph Place as `template.HTML` value for Go's `html/template`.
func (place *Place) ToGoHTMLMetaTags() (template.HTML, error) {
	// Create the templ component.
	templComponent := place.ToMetaTags()

	// Render the templ component to a `template.HTML` value.
	html, err := templ.ToGoHTML(context.Background(), templComponent)
	if err != nil {
		log.Fatalf("failed to convert to html: %v", err)
	}

	return html, nil
}

// ensureDefaults sets default values for Place.
func (place *Place) ensureDefaults() {
	place.OpenGraphObject.ensureDefaults("place")
}

// metaTags returns all meta tags for the Place object, including OpenGraphObject fields and place-specific ones.
func (place *Place) metaTags() []struct {
	property string
	content  string
} {
	return []struct {
		property string
		content  string
	}{
		{"og:type", "place"},
		{"og:title", place.Title},
		{"og:url", place.URL},
		{"og:description", place.Description},
		{"og:image", place.Image},
		{"place:location:latitude", fmt.Sprintf("%.4f", place.Latitude)},
		{"place:location:longitude", fmt.Sprintf("%.4f", place.Longitude)},
		{"place:contact_data:street_address", place.StreetAddress},
		{"place:contact_data:locality", place.Locality},
		{"place:contact_data:region", place.Region},
		{"place:contact_data:postal_code", place.PostalCode},
		{"place:contact_data:country_name", place.Country},
	}
}
