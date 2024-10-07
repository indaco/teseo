package opengraph

import (
	"context"
	"html/template"
	"io"
	"log"

	"github.com/a-h/templ"
	"github.com/indaco/teseo"
)

// Business represents the Open Graph business metadata.
// For more details about the meaning of the properties see: https://ogp.me/#metadata
//
// Example usage:
//
// Pure struct usage:
//
//	// Create a business using pure struct
//	business := &opengraph.Business{
//		OpenGraphObject: opengraph.OpenGraphObject{
//			Title:       "Example Business",
//			URL:         "https://www.example.com/business",
//			Description: "This is an example business description.",
//			Image:       "https://www.example.com/images/business.jpg",
//		},
//		StreetAddress: "123 Main St",
//		Locality:      "Anytown",
//		Region:        "CA",
//		PostalCode:    "12345",
//		Country:       "USA",
//		Email:         "info@example.com",
//		PhoneNumber:   "+1-800-555-1234",
//		Website:       "https://www.example.com",
//	}
//
// Factory method usage:
//
//	// Create a business
//	business := opengraph.NewBusiness(
//		"Example Business",
//		"https://www.example.com/business",
//		"This is an example business description.",
//		"https://www.example.com/images/business.jpg",
//		"123 Main St",
//		"Anytown",
//		"CA",
//		"12345",
//		"USA",
//		"info@example.com",
//		"+1-800-555-1234",
//		"https://www.example.com",
//	)
//
// // Rendering the HTML meta tags using templ:
//
//	templ Page() {
//		@business.ToMetaTgs()
//	}
//
// // Rendering the HTML meta tags as `template.HTML` value:
//
//	metaTagsHtml := business.ToGoHTMLMetaTgs()
//
// Expected output:
//
//	<meta property="og:type" content="business.business"/>
//	<meta property="og:title" content="Example Business"/>
//	<meta property="og:url" content="https://www.example.com/business"/>
//	<meta property="og:description" content="This is an example business description."/>
//	<meta property="og:image" content="https://www.example.com/images/business.jpg"/>
//	<meta property="business:contact_data:street_address" content="123 Main St"/>
//	<meta property="business:contact_data:locality" content="Anytown"/>
//	<meta property="business:contact_data:region" content="CA"/>
//	<meta property="business:contact_data:postal_code" content="12345"/>
//	<meta property="business:contact_data:country_name" content="USA"/>
//	<meta property="business:contact_data:email" content="info@example.com"/>
//	<meta property="business:contact_data:phone_number" content="+1-800-555-1234"/>
//	<meta property="business:contact_data:website" content="https://www.example.com"/>
type Business struct {
	OpenGraphObject
	StreetAddress string // business:contact_data:street_address, street address of the business
	Locality      string // business:contact_data:locality, locality or city of the business
	Region        string // business:contact_data:region, region or state of the business
	PostalCode    string // business:contact_data:postal_code, postal code of the business
	Country       string // business:contact_data:country_name, country of the business
	Email         string // business:contact_data:email, email address of the business
	PhoneNumber   string // business:contact_data:phone_number, phone number of the business
	Website       string // business:contact_data:website, website URL of the business
}

// NewBusiness initializes a Business with the default type "business.business".
func NewBusiness(title, url, description, image, streetAddress, locality, region, postalCode, country, email, phoneNumber, website string) *Business {
	business := &Business{
		OpenGraphObject: OpenGraphObject{
			Title:       title,
			URL:         url,
			Description: description,
			Image:       image,
		},
		StreetAddress: streetAddress,
		Locality:      locality,
		Region:        region,
		PostalCode:    postalCode,
		Country:       country,
		Email:         email,
		PhoneNumber:   phoneNumber,
		Website:       website,
	}
	business.ensureDefaults()
	return business
}

// ToMetaTags generates the HTML meta tags for the Open Graph Business as templ.Component.
func (bus *Business) ToMetaTags() templ.Component {
	bus.ensureDefaults()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, tag := range bus.metaTags() {
			if tag.content != "" {
				if err := teseo.WriteMetaTag(w, tag.property, tag.content); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// ToGoHTMLMetaTags generates the HTML meta tags for the Open Graph Business as `template.HTML` value for Go's `html/template`.
func (bus *Business) ToGoHTMLMetaTags() (template.HTML, error) {
	// Create the templ component.
	templComponent := bus.ToMetaTags()

	// Render the templ component to a `template.HTML` value.
	html, err := templ.ToGoHTML(context.Background(), templComponent)
	if err != nil {
		log.Fatalf("failed to convert to html: %v", err)
	}

	return html, nil
}

// ensureDefaults sets default values for Business.
func (bus *Business) ensureDefaults() {
	bus.OpenGraphObject.ensureDefaults("business.business")
}

// metaTags returns all meta tags for the Business object, including OpenGraphObject fields and business-specific ones.
func (bus *Business) metaTags() []struct {
	property string
	content  string
} {
	return []struct {
		property string
		content  string
	}{
		{"og:type", "business.business"},
		{"og:title", bus.Title},
		{"og:url", bus.URL},
		{"og:description", bus.Description},
		{"og:image", bus.Image},
		{"business:contact_data:street_address", bus.StreetAddress},
		{"business:contact_data:locality", bus.Locality},
		{"business:contact_data:region", bus.Region},
		{"business:contact_data:postal_code", bus.PostalCode},
		{"business:contact_data:country_name", bus.Country},
		{"business:contact_data:email", bus.Email},
		{"business:contact_data:phone_number", bus.PhoneNumber},
		{"business:contact_data:website", bus.Website},
	}
}
