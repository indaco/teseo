package opengraph

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/a-h/templ"
	"github.com/indaco/teseo"
)

// Audio represents the Open Graph audio metadata.
//
// Example usage:
//
// Pure struct usage:
//
//	// Create an audio object using pure struct
//	audio := &opengraph.Audio{
//		OpenGraphObject: opengraph.OpenGraphObject{
//			Title:       "Example Audio Title",
//			URL:         "https://www.example.com/audio/example-audio",
//			Description: "This is an example audio description.",
//			Image:       "https://www.example.com/images/audio.jpg",
//		},
//		Duration:  "300", // Duration in seconds
//		ArtistURL: "https://www.example.com/musicians/jane-doe",
//	}
//
// Factory method usage:
//
//	// Create an audio object using the factory method
//	audio := opengraph.NewAudio(
//		"Example Audio Title",
//		"https://www.example.com/audio/example-audio",
//		"This is an example audio description.",
//		"https://www.example.com/images/audio.jpg",
//		"300", // Duration in seconds
//		"https://www.example.com/musicians/jane-doe",
//	)
//
// // Rendering the HTML meta tags using templ:
//
//	templ Page() {
//		@audio.ToMetaTgs()
//	}
//
// // Rendering the HTML meta tags as `template.HTML` value:
//
//	metaTagsHtml := audio.ToGoHTMLMetaTgs()
//
// Expected output:
//
//	<meta property="og:type" content="music.audio"/>
//	<meta property="og:title" content="Example Audio Title"/>
//	<meta property="og:url" content="https://www.example.com/audio/example-audio"/>
//	<meta property="og:description" content="This is an example audio description."/>
//	<meta property="og:image" content="https://www.example.com/images/audio.jpg"/>
//	<meta property="music:duration" content="300"/>
//	<meta property="music:musician" content="https://www.example.com/musicians/jane-doe"/>
type Audio struct {
	OpenGraphObject
	Duration  string // music:duration, duration of the audio in seconds
	ArtistURL string // music:musician, URL to the musician or artist
}

// NewAudio initializes an Audio with the default type "music.audio".
func NewAudio(title, url, description, image, duration, artistURL string) *Audio {
	audio := &Audio{
		OpenGraphObject: OpenGraphObject{
			Title:       title,
			URL:         url,
			Description: description,
			Image:       image,
		},
		Duration:  duration,
		ArtistURL: artistURL,
	}
	audio.ensureDefaults()
	return audio
}

// ToMetaTags generates the HTML meta tags for the Open Graph Audio as templ.Component.
func (audio *Audio) ToMetaTags() templ.Component {
	audio.ensureDefaults()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, tag := range audio.metaTags() {
			if tag.content != "" {
				if err := teseo.WriteMetaTag(w, tag.property, tag.content); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// ToGoHTMLMetaTags generates the HTML meta tags for the Open Graph Audio as a string for Go's `html/template`.
func (audio *Audio) ToGoHTMLMetaTags() string {
	audio.ensureDefaults()

	var sb strings.Builder

	for _, tag := range audio.metaTags() {
		if tag.content != "" {
			sb.WriteString(fmt.Sprintf(`<meta property="%s" content="%s"/>`, tag.property, tag.content))
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// ensureDefaults sets default values for Audio.
func (audio *Audio) ensureDefaults() {
	audio.OpenGraphObject.ensureDefaults("music.audio")
}

// metaTags returns all meta tags for the Audio object, including OpenGraphObject fields and audio-specific ones.
func (audio *Audio) metaTags() []struct {
	property string
	content  string
} {
	return []struct {
		property string
		content  string
	}{
		{"og:type", "music.audio"},
		{"og:title", audio.Title},
		{"og:url", audio.URL},
		{"og:description", audio.Description},
		{"og:image", audio.Image},
		{"music:duration", audio.Duration},
		{"music:musician", audio.ArtistURL},
	}
}
