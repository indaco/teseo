package schemaorg

import (
	"context"
	"fmt"
	"html/template"
	"log"

	"github.com/a-h/templ"
	"github.com/indaco/teseo"
)

// Person represents a Schema.org Person object.
// For more details about the meaning of the properties see: https://schema.org/Person
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

// Person represents a Schema.org Person object
// For more details about the meaning of the properties see: https://schema.org/Person
type Person struct {
	Context     string         `json:"@context"`
	Type        string         `json:"@type"`
	Name        string         `json:"name,omitempty"`
	URL         string         `json:"url,omitempty"`
	Email       string         `json:"email,omitempty"`
	Image       *ImageObject   `json:"image,omitempty"`
	JobTitle    string         `json:"jobTitle,omitempty"`
	WorksFor    *Organization  `json:"worksFor,omitempty"`
	SameAs      []string       `json:"sameAs,omitempty"`
	Gender      string         `json:"gender,omitempty"`
	BirthDate   string         `json:"birthDate,omitempty"`
	Nationality string         `json:"nationality,omitempty"`
	Telephone   string         `json:"telephone,omitempty"`
	Address     *PostalAddress `json:"address,omitempty"`
	Affiliation *Organization  `json:"affiliation,omitempty"`
}

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

// Validate checks for recommended or required fields in Person.
func (p *Person) Validate() []string {
	var warnings []string

	if p.Name == "" {
		warnings = append(warnings, "missing recommended field: name")
	}

	if p.Email == "" {
		warnings = append(warnings, "missing recommended field: email")
	}

	if p.JobTitle == "" {
		warnings = append(warnings, "missing recommended field: jobTitle")
	}

	return warnings
}

// ToJsonLd converts the Person struct to a JSON-LD `templ.Component`.
func (p *Person) ToJsonLd() templ.Component {
	p.ensureDefaults()
	id := fmt.Sprintf("%s-%s", "person", teseo.GenerateUniqueKey())
	return templ.JSONScript(id, p).WithType("application/ld+json")
}

// ToGoHTMLJsonLd renders the Person struct as `template.HTML` value for Go's `html/template`.
func (p *Person) ToGoHTMLJsonLd() (template.HTML, error) {
	return teseo.RenderToHTML(p.ToJsonLd())
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
