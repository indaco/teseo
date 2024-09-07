package schemaorg

import (
	"context"
	"fmt"
	"html/template"
	"log"

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
	id := fmt.Sprintf("%s-%s", "event", teseo.GenerateUniqueKey())
	return templ.JSONScript(id, e).WithType("application/ld+json")
}

// ToGoHTMLJsonLd renders the Event struct as `template.HTML` value for Go's `html/template`.
func (e *Event) ToGoHTMLJsonLd() (template.HTML, error) {
	// Create the templ component.
	templComponent := e.ToJsonLd()

	// Render the templ component to a `template.HTML` value.
	html, err := templ.ToGoHTML(context.Background(), templComponent)
	if err != nil {
		log.Fatalf("failed to convert to html: %v", err)
	}

	return html, nil
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
