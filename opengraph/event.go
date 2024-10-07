package opengraph

import (
	"context"
	"html/template"
	"io"
	"log"

	"github.com/a-h/templ"
	"github.com/indaco/teseo"
)

// Event represents the Open Graph event metadata.
// For more details about the meaning of the properties see: https://ogp.me/#metadata
//
// Example usage:
//
// Pure struct usage:
//
//	// Create an event using pure struct
//	event := &opengraph.Event{
//		OpenGraphObject: opengraph.OpenGraphObject{
//			Title:       "Example Event Title",
//			URL:         "https://www.example.com/event/example-event",
//			Description: "This is an example event description.",
//			Image:       "https://www.example.com/images/event.jpg",
//		},
//		StartDate: "2024-09-15T09:00:00Z",
//		EndDate:   "2024-09-15T18:00:00Z",
//		Location:  "Anytown Convention Center",
//	}
//
// Factory method usage:
//
//	// Create an event using the factory method
//	event := opengraph.NewEvent(
//		"Example Event Title",
//		"https://www.example.com/event/example-event",
//		"This is an example event description.",
//		"https://www.example.com/images/event.jpg",
//		"2024-09-15T09:00:00Z",
//		"2024-09-15T18:00:00Z",
//		"Anytown Convention Center",
//	)
//
// // Rendering the HTML meta tags using templ:
//
//	templ Page() {
//		@event.ToMetaTgs()
//	}
//
// // Rendering the HTML meta tags as `template.HTML` value:
//
//	metaTagsHtml := event.ToGoHTMLMetaTgs()
//
// Expected output:
//
//	<meta property="og:type" content="event"/>
//	<meta property="og:title" content="Example Event Title"/>
//	<meta property="og:url" content="https://www.example.com/event/example-event"/>
//	<meta property="og:description" content="This is an example event description."/>
//	<meta property="og:image" content="https://www.example.com/images/event.jpg"/>
//	<meta property="event:start_date" content="2024-09-15T09:00:00Z"/>
//	<meta property="event:end_date" content="2024-09-15T18:00:00Z"/>
//	<meta property="event:location" content="Anytown Convention Center"/>
type Event struct {
	OpenGraphObject
	StartDate string // event:start_date, the start date and time of the event
	EndDate   string // event:end_date, the end date and time of the event
	Location  string // event:location, the location of the event
}

// NewEvent initializes an Event with the default type "event".
func NewEvent(title, url, description, image, startDate, endDate, location string) *Event {
	event := &Event{
		OpenGraphObject: OpenGraphObject{
			Title:       title,
			URL:         url,
			Description: description,
			Image:       image,
		},
		StartDate: startDate,
		EndDate:   endDate,
		Location:  location,
	}
	event.ensureDefaults()
	return event
}

// ToMetaTags generates the HTML meta tags for the Open Graph Event as templ.Component.
func (e *Event) ToMetaTags() templ.Component {
	e.ensureDefaults()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, tag := range e.metaTags() {
			if tag.content != "" {
				if err := teseo.WriteMetaTag(w, tag.property, tag.content); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// ToGoHTMLMetaTags generates the HTML meta tags for the Open Graph Event as `template.HTML` value for Go's `html/template`.
func (e *Event) ToGoHTMLMetaTags() (template.HTML, error) {
	// Create the templ component.
	templComponent := e.ToMetaTags()

	// Render the templ component to a `template.HTML` value.
	html, err := templ.ToGoHTML(context.Background(), templComponent)
	if err != nil {
		log.Fatalf("failed to convert to html: %v", err)
	}

	return html, nil
}

// ensureDefaults sets default values for Event.
func (e *Event) ensureDefaults() {
	e.OpenGraphObject.ensureDefaults("event")
}

// metaTags returns all meta tags for the Event object, including OpenGraphObject fields and event-specific ones.
func (e *Event) metaTags() []struct {
	property string
	content  string
} {
	return []struct {
		property string
		content  string
	}{
		{"og:type", "event"},
		{"og:title", e.Title},
		{"og:url", e.URL},
		{"og:description", e.Description},
		{"og:image", e.Image},
		{"event:start_date", e.StartDate},
		{"event:end_date", e.EndDate},
		{"event:location", e.Location},
	}
}
