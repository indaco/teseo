package schemaorg

import (
	"context"
	"fmt"
	"html"
	"io"
	"strings"

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
		return templ.JSONScript(teseo.GenerateUniqueKey(), p).WithType("application/ld+json").Render(ctx, w)
	})
}

// ToGoHTMLJsonLd renders the Person struct as a string for Go's `html/template`.
func (p *Person) ToGoHTMLJsonLd() string {
	p.ensureDefaults()

	var sb strings.Builder
	sb.WriteString(`<script type="application/ld+json">`)
	sb.WriteString("\n{\n")
	sb.WriteString(fmt.Sprintf(`  "@context": "%s",`, html.EscapeString(p.Context)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(`  "@type": "%s",`, html.EscapeString(p.Type)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(`  "name": "%s",`, html.EscapeString(p.Name)))
	sb.WriteString("\n")
	if p.URL != "" {
		sb.WriteString(fmt.Sprintf(`  "url": "%s",`, html.EscapeString(p.URL)))
		sb.WriteString("\n")
	}
	if p.Email != "" {
		sb.WriteString(fmt.Sprintf(`  "email": "%s",`, html.EscapeString(p.Email)))
		sb.WriteString("\n")
	}
	if p.Image != nil {
		sb.WriteString(`  "image": {`)
		sb.WriteString(fmt.Sprintf(`"@type": "ImageObject", "url": "%s"`, html.EscapeString(p.Image.URL)))
		sb.WriteString("},\n")
	}
	if p.JobTitle != "" {
		sb.WriteString(fmt.Sprintf(`  "jobTitle": "%s",`, html.EscapeString(p.JobTitle)))
		sb.WriteString("\n")
	}
	if p.WorksFor != nil {
		sb.WriteString(`  "worksFor": {`)
		sb.WriteString(fmt.Sprintf(`"@type": "Organization", "name": "%s", "url": "%s"`, html.EscapeString(p.WorksFor.Name), html.EscapeString(p.WorksFor.URL)))
		sb.WriteString("},\n")
	}
	if len(p.SameAs) > 0 {
		sb.WriteString(`  "sameAs": [`)
		for i, sameAs := range p.SameAs {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf(`"%s"`, html.EscapeString(sameAs)))
		}
		sb.WriteString("],\n")
	}
	if p.Gender != "" {
		sb.WriteString(fmt.Sprintf(`  "gender": "%s",`, html.EscapeString(p.Gender)))
		sb.WriteString("\n")
	}
	if p.BirthDate != "" {
		sb.WriteString(fmt.Sprintf(`  "birthDate": "%s",`, html.EscapeString(p.BirthDate)))
		sb.WriteString("\n")
	}
	if p.Nationality != "" {
		sb.WriteString(fmt.Sprintf(`  "nationality": "%s",`, html.EscapeString(p.Nationality)))
		sb.WriteString("\n")
	}
	if p.Telephone != "" {
		sb.WriteString(fmt.Sprintf(`  "telephone": "%s",`, html.EscapeString(p.Telephone)))
		sb.WriteString("\n")
	}
	if p.Address != nil {
		sb.WriteString(`  "address": {`)
		sb.WriteString(fmt.Sprintf(`"@type": "PostalAddress", "addressLocality": "%s", "addressCountry": "%s"`, html.EscapeString(p.Address.AddressLocality), html.EscapeString(p.Address.AddressCountry)))
		sb.WriteString("},\n")
	}
	if p.Affiliation != nil {
		sb.WriteString(`  "affiliation": {`)
		sb.WriteString(fmt.Sprintf(`"@type": "Organization", "name": "%s"`, html.EscapeString(p.Affiliation.Name)))
		sb.WriteString("},\n")
	}

	sb.WriteString("}\n</script>")
	return sb.String()
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
