package opengraph

import (
	"context"
	"html/template"
	"io"
	"log"

	"github.com/a-h/templ"
	"github.com/indaco/teseo"
)

// MusicRadioStation represents the Open Graph music radio station metadata.
// For more details about the meaning of the properties see: https://ogp.me/#type_music.radio_station
//
// Example usage:
//
// Pure struct usage:
//
//	// Create a music radio station using pure struct
//	musicRadioStation := &opengraph.MusicRadioStation{
//		OpenGraphObject: opengraph.OpenGraphObject{
//			Title:       "Example Radio Station",
//			URL:         "https://www.example.com/music/radio/example-radio",
//			Description: "This is an example radio station description.",
//			Image:       "https://www.example.com/images/radio.jpg",
//		},
//	}
//
// Factory method usage:
//
//	// Create a music radio station using the factory method
//	musicRadioStation := opengraph.NewMusicRadioStation(
//		"Example Radio Station",
//		"https://www.example.com/music/radio/example-radio",
//		"This is an example radio station description.",
//		"https://www.example.com/images/radio.jpg",
//	)
//
// // Rendering the HTML meta tags using templ:
//
//	templ Page() {
//		@musicRadioStation.ToMetaTgs()
//	}
//
// // Rendering the HTML meta tags as `template.HTML` value:
//
//	metaTagsHtml := musicRadioStation.ToGoHTMLMetaTgs()
//
// Expected output:
//
//	<meta property="og:type" content="music.radio_station"/>
//	<meta property="og:title" content="Example Radio Station"/>
//	<meta property="og:url" content="https://www.example.com/music/radio/example-radio"/>
//	<meta property="og:description" content="This is an example radio station description."/>
//	<meta property="og:image" content="https://www.example.com/images/radio.jpg"/>
type MusicRadioStation struct {
	OpenGraphObject
}

// NewMusicRadioStation initializes a MusicRadioStation with the default type "music.radio_station".
func NewMusicRadioStation(title, url, description, image string) *MusicRadioStation {
	musicRadioStation := &MusicRadioStation{
		OpenGraphObject: OpenGraphObject{
			Title:       title,
			URL:         url,
			Description: description,
			Image:       image,
		},
	}
	musicRadioStation.ensureDefaults()
	return musicRadioStation
}

// ToMetaTags generates the HTML meta tags for the Open Graph Music Radio Station as templ.Component.
func (mrs *MusicRadioStation) ToMetaTags() templ.Component {
	mrs.ensureDefaults()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, tag := range mrs.metaTags() {
			if tag.content != "" {
				if err := teseo.WriteMetaTag(w, tag.property, tag.content); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// ToGoHTMLMetaTags generates the HTML meta tags for the Open Graph Music Radio Station as `template.HTML` value for Go's `html/template`.
func (mrs *MusicRadioStation) ToGoHTMLMetaTags() (template.HTML, error) {
	// Create the templ component.
	templComponent := mrs.ToMetaTags()

	// Render the templ component to a `template.HTML` value.
	html, err := templ.ToGoHTML(context.Background(), templComponent)
	if err != nil {
		log.Fatalf("failed to convert to html: %v", err)
	}

	return html, nil
}

// ensureDefaults sets default values for MusicRadioStation.
func (mrs *MusicRadioStation) ensureDefaults() {
	mrs.OpenGraphObject.ensureDefaults("music.radio_station")
}

// metaTags returns all meta tags for the MusicRadioStation object, including OpenGraphObject fields.
func (mrs *MusicRadioStation) metaTags() []struct {
	property string
	content  string
} {
	return []struct {
		property string
		content  string
	}{
		{"og:type", "music.radio_station"},
		{"og:title", mrs.Title},
		{"og:url", mrs.URL},
		{"og:description", mrs.Description},
		{"og:image", mrs.Image},
	}
}
