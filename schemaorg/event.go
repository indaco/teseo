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

// Event represents a Schema.org Event object.
//
// Example usage:
//
// Pure struct usage:
//
//	event := &schemaorg.Event{
//		Name:        "Example Event",
//		StartDate:   "2024-09-20T19:00:00",
//		EndDate:     "2024-09-20T23:00:00",
//		Location:    &schemaorg.Place{Name: "Example Venue", Address: "123 Main St"},
//		Description: "This is an example event.",
//	}
//
// Factory method usage:
//
//	event := schemaorg.NewEvent(
//		"Example Event",
//		"2024-09-20T19:00:00",
//		"2024-09-20T23:00:00",
//		&schemaorg.Place{Name: "Example Venue", Address: "123 Main St"},
//		"This is an example event",
//	)
//
// // Rendering JSON-LD using templ:
//
//	templ Page() {
//		@event.ToJsonLd()
//	}
//
// // Rendering JSON-LD as `template.HTML` value:
//
//	jsonLdHtml := event.ToGoHTMLJsonLd()
//
// Expected output:
//
//	{
//		"@context": "https://schema.org",
//		"@type": "Event",
//		"name": "Example Event",
//		"startDate": "2024-09-20T19:00:00",
//		"endDate": "2024-09-20T23:00:00",
//		"location": {"@type": "Place", "name": "Example Venue", "address": "123 Main St"},
//		"description": "This is an example event"
//	}
type Event struct {
	Context             string        `json:"@context"`
	Type                string        `json:"@type"`
	Name                string        `json:"name,omitempty"`
	Description         string        `json:"description,omitempty"`
	StartDate           string        `json:"startDate,omitempty"`
	EndDate             string        `json:"endDate,omitempty"`
	Location            *Place        `json:"location,omitempty"`
	Organizer           *Organization `json:"organizer,omitempty"`
	Performer           *Person       `json:"performer,omitempty"`
	Image               []string      `json:"image,omitempty"`
	EventStatus         string        `json:"eventStatus,omitempty"`
	EventAttendanceMode string        `json:"eventAttendanceMode,omitempty"`
	Offers              *Offer        `json:"offers,omitempty"`
}

// Place represents a Schema.org Place object
type Place struct {
	Context string          `json:"@context"`
	Type    string          `json:"@type"`
	Name    string          `json:"name,omitempty"`
	Address *PostalAddress  `json:"address,omitempty"`
	Geo     *GeoCoordinates `json:"geo,omitempty"`
}

// NewEvent initializes an Event with default context and type.
func NewEvent(name, description, startDate, endDate string, location *Place, organizer *Organization, performer *Person, images []string, eventStatus, eventAttendanceMode string, offers *Offer) *Event {
	event := &Event{
		Name:                name,
		Description:         description,
		StartDate:           startDate,
		EndDate:             endDate,
		Location:            location,
		Organizer:           organizer,
		Performer:           performer,
		Image:               images,
		EventStatus:         eventStatus,
		EventAttendanceMode: eventAttendanceMode,
		Offers:              offers,
	}
	event.ensureDefaults()
	return event
}

// ToJsonLd converts the Event struct to a JSON-LD `templ.Component`.
func (e *Event) ToJsonLd() templ.Component {
	e.ensureDefaults()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		return templ.JSONScript(teseo.GenerateUniqueKey(), e).WithType("application/ld+json").Render(ctx, w)
	})
}

// ToGoHTMLJsonLd renders the Event struct as a string for Go's `html/template`.
func (e *Event) ToGoHTMLJsonLd() string {
	e.ensureDefaults()

	var sb strings.Builder
	sb.WriteString(`<script type="application/ld+json">`)
	sb.WriteString("\n{\n")
	sb.WriteString(fmt.Sprintf(`  "@context": "%s",`, html.EscapeString(e.Context)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(`  "@type": "%s",`, html.EscapeString(e.Type)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(`  "name": "%s",`, html.EscapeString(e.Name)))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(`  "description": "%s",`, html.EscapeString(e.Description)))
	sb.WriteString("\n")

	if e.StartDate != "" {
		sb.WriteString(fmt.Sprintf(`  "startDate": "%s",`, html.EscapeString(e.StartDate)))
		sb.WriteString("\n")
	}
	if e.EndDate != "" {
		sb.WriteString(fmt.Sprintf(`  "endDate": "%s",`, html.EscapeString(e.EndDate)))
		sb.WriteString("\n")
	}

	// EventStatus
	if e.EventStatus != "" {
		sb.WriteString(fmt.Sprintf(`  "eventStatus": "%s",`, html.EscapeString(e.EventStatus)))
		sb.WriteString("\n")
	}

	// EventAttendanceMode
	if e.EventAttendanceMode != "" {
		sb.WriteString(fmt.Sprintf(`  "eventAttendanceMode": "%s",`, html.EscapeString(e.EventAttendanceMode)))
		sb.WriteString("\n")
	}

	// Location
	if e.Location != nil {
		sb.WriteString(`  "location": {`)
		sb.WriteString(fmt.Sprintf(`"@type": "Place", "name": "%s"`, html.EscapeString(e.Location.Name)))
		if e.Location.Address != nil {
			sb.WriteString(`, "address": {`)
			sb.WriteString(fmt.Sprintf(`"@type": "PostalAddress", "addressLocality": "%s", "addressCountry": "%s"`, html.EscapeString(e.Location.Address.AddressLocality), html.EscapeString(e.Location.Address.AddressCountry)))
			sb.WriteString("}")
		}
		if e.Location.Geo != nil {
			sb.WriteString(`, "geo": {`)
			sb.WriteString(fmt.Sprintf(`"@type": "GeoCoordinates", "latitude": %f, "longitude": %f`, e.Location.Geo.Latitude, e.Location.Geo.Longitude))
			sb.WriteString("}")
		}
		sb.WriteString("},\n")
	}

	// Organizer
	if e.Organizer != nil {
		sb.WriteString(`  "organizer": {`)
		sb.WriteString(fmt.Sprintf(`"@type": "Organization", "name": "%s"`, html.EscapeString(e.Organizer.Name)))
		sb.WriteString("},\n")
	}

	// Performer
	if e.Performer != nil {
		sb.WriteString(`  "performer": {`)
		sb.WriteString(fmt.Sprintf(`"@type": "Person", "name": "%s"`, html.EscapeString(e.Performer.Name)))
		sb.WriteString("},\n")
	}

	// Images
	if len(e.Image) > 0 {
		sb.WriteString(`  "image": [`)
		for i, img := range e.Image {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf(`"%s"`, html.EscapeString(img)))
		}
		sb.WriteString("],\n")
	}

	// Offers
	if e.Offers != nil {
		sb.WriteString(`  "offers": {`)
		sb.WriteString(fmt.Sprintf(`"@type": "Offer", "price": "%s", "priceCurrency": "%s"`, html.EscapeString(e.Offers.Price), html.EscapeString(e.Offers.PriceCurrency)))
		sb.WriteString("},\n")
	}

	sb.WriteString("}\n</script>")
	return sb.String()
}

// ensureDefaults sets default values for Event and its nested objects if they are not already set.
func (e *Event) ensureDefaults() {
	if e.Context == "" {
		e.Context = "https://schema.org"
	}
	if e.Type == "" {
		e.Type = "Event"
	}

	if e.Location != nil {
		e.Location.ensureDefaults()
	}

	if e.Organizer != nil {
		e.Organizer.ensureDefaults()
	}

	if e.Performer != nil {
		e.Performer.ensureDefaults()
	}

	if e.Offers != nil {
		e.Offers.ensureDefaults()
	}
}

// ensureDefaults sets default values for Place if they are not already set.
func (p *Place) ensureDefaults() {
	if p.Context == "" {
		p.Context = "https://schema.org"
	}
	if p.Type == "" {
		p.Type = "Place"
	}

	if p.Address != nil {
		p.Address.ensureDefaults()
	}

	if p.Geo != nil {
		p.Geo.ensureDefaults()
	}
}
