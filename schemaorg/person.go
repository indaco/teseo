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

// Person represents a Schema.org Person object.
//
// Example usage:
//
// Pure struct usage:
//
// 	person := &schemaorg.Person{
// 		Name:  "Jane Doe",
// 		Email: "jane.doe@example.com",
// 		JobTitle: "Software Engineer",
// 		WorksFor: &schemaorg.Organization{Name: "Example Company"},
// 	}
//
// Factory method usage:
//
// 	person := schemaorg.NewPerson(
// 		"Jane Doe",
// 		"jane.doe@example.com",
// 		"Software Engineer",
// 		&schemaorg.Organization{Name: "Example Company"},
// 	)
//
// // Rendering JSON-LD using templ:
//
//	templ Page() {
//		@person.ToJsonLd()
//	}
//
// // Rendering JSON-LD as `template.HTML` value:
//
//	jsonLdHtml := person.ToGoHTMLJsonLd()
//
// Expected output:
//
// 	{
// 		"@context": "https://schema.org",
// 		"@type": "Person",
// 		"name": "Jane Doe",
// 		"email": "jane.doe@example.com",
// 		"jobTitle": "Software Engineer",
// 		"worksFor": {"@type": "Organization", "name": "Example Company"}
// 	}

// PostalAddress represents a Schema.org PostalAddress object
type PostalAddress struct {
	Type            string `json:"@type"`
	StreetAddress   string `json:"streetAddress,omitempty"`
	AddressLocality string `json:"addressLocality,omitempty"`
	AddressRegion   string `json:"addressRegion,omitempty"`
	PostalCode      string `json:"postalCode,omitempty"`
	AddressCountry  string `json:"addressCountry,omitempty"`
}

// ensureDefaults sets default values for PostalAddress if they are not already set.
func (addr *PostalAddress) ensureDefaults() {
	if addr.Type == "" {
		addr.Type = "PostalAddress"
	}
}

// NewPerson initializes a Person with default context and type.
func NewPerson(name string, url string, email string, image *ImageObject, jobTitle string, worksFor *Organization, sameAs []string, gender string, birthDate string, nationality string, telephone string, address *PostalAddress, affiliation *Organization) *Person {
	person := &Person{
		Name:        name,
		URL:         url,
		Email:       email,
		Image:       image,
		JobTitle:    jobTitle,
		WorksFor:    worksFor,
		SameAs:      sameAs,
		Gender:      gender,
		BirthDate:   birthDate,
		Nationality: nationality,
		Telephone:   telephone,
		Address:     address,
		Affiliation: affiliation,
	}
	person.ensureDefaults()
	return person
}

// ToJsonLd converts the Person struct to a JSON-LD `templ.Component`.
func (p *Person) ToJsonLd() templ.Component {
	p.ensureDefaults()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		id := fmt.Sprintf("%s-%s", "person", teseo.GenerateUniqueKey())
		return templ.JSONScript(id, p).WithType("application/ld+json").Render(ctx, w)
	})
}

// ToGoHTMLJsonLd renders the Person struct as `template.HTML` value for Go's `html/template`.
func (p *Person) ToGoHTMLJsonLd() (template.HTML, error) {
	// Create the templ component.
	templComponent := p.ToJsonLd()

	// Render the templ component to a `template.HTML` value.
	html, err := templ.ToGoHTML(context.Background(), templComponent)
	if err != nil {
		log.Fatalf("failed to convert to html: %v", err)
	}

	return html, nil
}

// ensureDefaults sets default values for Person and its nested objects if they are not already set.
func (p *Person) ensureDefaults() {
	if p.Context == "" {
		p.Context = "https://schema.org"
	}

	if p.Type == "" {
		p.Type = "Person"
	}

	if p.Image != nil {
		p.Image.ensureDefaults()
	}

	if p.WorksFor != nil {
		p.WorksFor.ensureDefaults()
	}

	if p.Address != nil {
		p.Address.ensureDefaults()
	}

	if p.Affiliation != nil {
		p.Affiliation.ensureDefaults()
	}
}
