package opengraph

import (
	"context"
	"html/template"
	"io"
	"log"

	"github.com/a-h/templ"
	"github.com/indaco/teseo"
)

// VideoMovie represents the Open Graph video movie metadata.
// For more details about the meaning of the properties see: https://ogp.me/#type_video.movie
//
// Example usage:
//
// Pure struct usage:
//
//	// Create a video movie using pure struct
//	videoMovie := &opengraph.VideoMovie{
//		OpenGraphObject: opengraph.OpenGraphObject{
//			Title:       "Example Movie",
//			URL:         "https://www.example.com/video/movie/example-movie",
//			Description: "This is an example movie description.",
//			Image:       "https://www.example.com/images/movie.jpg",
//		},
//		Duration:    "7200", // Duration in seconds (2 hours)
//		ActorURLs:   []string{"https://www.example.com/actors/jane-doe", "https://www.example.com/actors/john-doe"},
//		DirectorURL: "https://www.example.com/directors/jane-director",
//		ReleaseDate: "2024-09-15",
//	}
//
// Factory method usage:
//
//	// Create a video movie using the factory method
//	videoMovie := opengraph.NewVideoMovie(
//		"Example Movie",
//		"https://www.example.com/video/movie/example-movie",
//		"This is an example movie description.",
//		"https://www.example.com/images/movie.jpg",
//		"7200", // Duration in seconds (2 hours)
//		[]string{"https://www.example.com/actors/jane-doe", "https://www.example.com/actors/john-doe"},
//		"https://www.example.com/directors/jane-director",
//		"2024-09-15",
//	)
//
// // Rendering the HTML meta tags using templ:
//
//	templ Page() {
//		@videoMovie.ToMetaTgs()
//	}
//
// // Rendering the HTML meta tags as `template.HTML` value:
//
//	metaTagsHtml := videoMovie.ToGoHTMLMetaTgs()
//
// Expected output:
//
//	<meta property="og:type" content="video.movie"/>
//	<meta property="og:title" content="Example Movie"/>
//	<meta property="og:url" content="https://www.example.com/video/movie/example-movie"/>
//	<meta property="og:description" content="This is an example movie description."/>
//	<meta property="og:image" content="https://www.example.com/images/movie.jpg"/>
//	<meta property="video:duration" content="7200"/>
//	<meta property="video:actor" content="https://www.example.com/actors/jane-doe"/>
//	<meta property="video:actor" content="https://www.example.com/actors/john-doe"/>
//	<meta property="video:director" content="https://www.example.com/directors/jane-director"/>
//	<meta property="video:release_date" content="2024-09-15"/>
type VideoMovie struct {
	OpenGraphObject
	Duration    string   // video:duration, duration of the movie in seconds
	ActorURLs   []string // video:actor, URLs to the actors in the movie
	DirectorURL string   // video:director, URL to the director of the movie
	ReleaseDate string   // video:release_date, the release date of the movie
}

// NewVideoMovie initializes a VideoMovie with the default type "video.movie".
func NewVideoMovie(title, url, description, image, duration string, actorURLs []string, directorURL, releaseDate string) *VideoMovie {
	videoMovie := &VideoMovie{
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
	videoMovie.ensureDefaults()
	return videoMovie
}

// ToMetaTags generates the HTML meta tags for the Open Graph Video Movie as templ.Component.
func (vm *VideoMovie) ToMetaTags() templ.Component {
	vm.ensureDefaults()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, tag := range vm.metaTags() {
			if tag.content != "" {
				if err := teseo.WriteMetaTag(w, tag.property, tag.content); err != nil {
					return err
				}
			}
		}

		// Write video:actor meta tags for each actor URL
		for _, actorURL := range vm.ActorURLs {
			if actorURL != "" {
				if err := teseo.WriteMetaTag(w, "video:actor", actorURL); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// ToGoHTMLMetaTags generates the HTML meta tags for the Open Graph Video Movie as `template.HTML` value for Go's `html/template`.
func (vm *VideoMovie) ToGoHTMLMetaTags() (template.HTML, error) {
	// Create the templ component.
	templComponent := vm.ToMetaTags()

	// Render the templ component to a `template.HTML` value.
	html, err := templ.ToGoHTML(context.Background(), templComponent)
	if err != nil {
		log.Fatalf("failed to convert to html: %v", err)
	}

	return html, nil
}

// ensureDefaults sets default values for VideoMovie.
func (vm *VideoMovie) ensureDefaults() {
	vm.OpenGraphObject.ensureDefaults("video.movie")
}

// metaTags returns all meta tags for the VideoMovie object, including OpenGraphObject fields and video movie-specific ones.
func (vm *VideoMovie) metaTags() []struct {
	property string
	content  string
} {
	return []struct {
		property string
		content  string
	}{
		{"og:type", "video.movie"},
		{"og:title", vm.Title},
		{"og:url", vm.URL},
		{"og:description", vm.Description},
		{"og:image", vm.Image},
		{"video:duration", vm.Duration},
		{"video:director", vm.DirectorURL},
		{"video:release_date", vm.ReleaseDate},
	}
}
