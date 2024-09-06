package opengraph

import (
	"context"
	"html/template"
	"io"
	"log"

	"github.com/a-h/templ"
	"github.com/indaco/teseo"
)

// Video represents the Open Graph video metadata.
//
// Example usage:
//
// Pure struct usage:
//
//	// Create a video using pure struct
//	video := &opengraph.Video{
//		OpenGraphObject: opengraph.OpenGraphObject{
//			Title:       "Example Video",
//			URL:         "https://www.example.com/video/example-video",
//			Description: "This is an example video description.",
//			Image:       "https://www.example.com/images/video.jpg",
//		},
//		Duration: "300", // Duration in seconds
//		ActorURLs: []string{
//			"https://www.example.com/actors/jane-doe",
//			"https://www.example.com/actors/john-doe",
//		},
//		DirectorURL: "https://www.example.com/directors/jane-director",
//		ReleaseDate: "2024-09-15",
//	}
//
// Factory method usage:
//
//	// Create a video using the factory method
//	video := opengraph.NewVideo(
//		"Example Video",
//		"https://www.example.com/video/example-video",
//		"This is an example video description.",
//		"https://www.example.com/images/video.jpg",
//		"300", // Duration in seconds
//		[]string{"https://www.example.com/actors/jane-doe", "https://www.example.com/actors/john-doe"},
//		"https://www.example.com/directors/jane-director",
//		"2024-09-15",
//	)
//
// // Rendering the HTML meta tags using templ:
//
//	templ Page() {
//		@video.ToMetaTgs()
//	}
//
// // Rendering the HTML meta tags as `template.HTML` value:
//
//	metaTagsHtml := video.ToGoHTMLMetaTgs()
//
// Expected output:
//
//	<meta property="og:type" content="video.movie"/>
//	<meta property="og:title" content="Example Video"/>
//	<meta property="og:url" content="https://www.example.com/video/example-video"/>
//	<meta property="og:description" content="This is an example video description."/>
//	<meta property="og:image" content="https://www.example.com/images/video.jpg"/>
//	<meta property="video:duration" content="300"/>
//	<meta property="video:actor" content="https://www.example.com/actors/jane-doe"/>
//	<meta property="video:actor" content="https://www.example.com/actors/john-doe"/>
//	<meta property="video:director" content="https://www.example.com/directors/jane-director"/>
//	<meta property="video:release_date" content="2024-09-15"/>
type Video struct {
	OpenGraphObject
	Duration    string   // video:duration, duration of the video in seconds
	ActorURLs   []string // video:actor, URLs to the actors in the video
	DirectorURL string   // video:director, URL to the director of the video
	ReleaseDate string   // video:release_date, the release date of the video
}

// NewVideo initializes a Video with the default type "video.movie".
func NewVideo(title, url, description, image, duration string, actorURLs []string, directorURL, releaseDate string) *Video {
	video := &Video{
		OpenGraphObject: OpenGraphObject{
			Title:       title,
			URL:         url,
			Description: description,
			Image:       image,
		},
		Duration:    duration,
		ActorURLs:   actorURLs,
		DirectorURL: directorURL,
		ReleaseDate: releaseDate,
	}
	video.ensureDefaults()
	return video
}

// ToMetaTags generates the HTML meta tags for the Open Graph Video using templ.Component.
func (video *Video) ToMetaTags() templ.Component {
	video.ensureDefaults()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, tag := range video.metaTags() {
			if tag.content != "" {
				if err := teseo.WriteMetaTag(w, tag.property, tag.content); err != nil {
					return err
				}
			}
		}

		// Write video:actor meta tags for each actor URL
		for _, actorURL := range video.ActorURLs {
			if actorURL != "" {
				if err := teseo.WriteMetaTag(w, "video:actor", actorURL); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// ToGoHTMLMetaTags generates the HTML meta tags for the Open Graph Video as `template.HTML` value for Go's `html/template`.
func (video *Video) ToGoHTMLMetaTags() (template.HTML, error) {
	// Create the templ component.
	templComponent := video.ToMetaTags()

	// Render the templ component to a `template.HTML` value.
	html, err := templ.ToGoHTML(context.Background(), templComponent)
	if err != nil {
		log.Fatalf("failed to convert to html: %v", err)
	}

	return html, nil
}

// ensureDefaults sets default values for Video.
func (video *Video) ensureDefaults() {
	video.OpenGraphObject.ensureDefaults("video.movie")
}

// metaTags returns all meta tags for the Video object, including OpenGraphObject fields and video-specific ones.
func (video *Video) metaTags() []struct {
	property string
	content  string
} {
	return []struct {
		property string
		content  string
	}{
		{"og:type", "video.movie"},
		{"og:title", video.Title},
		{"og:url", video.URL},
		{"og:description", video.Description},
		{"og:image", video.Image},
		{"video:duration", video.Duration},
		{"video:director", video.DirectorURL},
		{"video:release_date", video.ReleaseDate},
	}
}
