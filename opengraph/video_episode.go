package opengraph

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/a-h/templ"
	"github.com/indaco/teseo"
)

// VideoEpisode represents the Open Graph video episode metadata.
//
// Example usage:
//
// Pure struct usage:
//
//	// Create a video episode using pure struct
//	videoEpisode := &opengraph.VideoEpisode{
//		OpenGraphObject: opengraph.OpenGraphObject{
//			Title:       "Example Video Episode",
//			URL:         "https://www.example.com/video/episode/example-episode",
//			Description: "This is an example video episode description.",
//			Image:       "https://www.example.com/images/episode.jpg",
//		},
//		SeriesURL:   "https://www.example.com/video/series/example-series",
//		Duration:    "1800", // Duration in seconds
//		ActorURLs:   []string{"https://www.example.com/actors/jane-doe", "https://www.example.com/actors/john-doe"},
//		DirectorURL: "https://www.example.com/directors/jane-director",
//		ReleaseDate: "2024-09-15",
//		EpisodeNumber: 1,
//	}
//
// Factory method usage:
//
//	// Create a video episode using the factory method
//	videoEpisode := opengraph.NewVideoEpisode(
//		"Example Video Episode",
//		"https://www.example.com/video/episode/example-episode",
//		"This is an example video episode description.",
//		"https://www.example.com/images/episode.jpg",
//		"1800", // Duration in seconds
//		"https://www.example.com/video/series/example-series",
//		[]string{"https://www.example.com/actors/jane-doe", "https://www.example.com/actors/john-doe"},
//		"https://www.example.com/directors/jane-director",
//		"2024-09-15",
//		1, // Episode number
//	)
//
// // Rendering the HTML meta tags using templ:
//
//	templ Page() {
//		@videoEpisode.ToMetaTgs()
//	}
//
// // Rendering the HTML meta tags as `template.HTML` value:
//
//	metaTagsHtml := videoEpisode.ToGoHTMLMetaTgs()
//
// Expected output:
//
//	<meta property="og:type" content="video.episode"/>
//	<meta property="og:title" content="Example Video Episode"/>
//	<meta property="og:url" content="https://www.example.com/video/episode/example-episode"/>
//	<meta property="og:description" content="This is an example video episode description."/>
//	<meta property="og:image" content="https://www.example.com/images/episode.jpg"/>
//	<meta property="video:duration" content="1800"/>
//	<meta property="video:actor" content="https://www.example.com/actors/jane-doe"/>
//	<meta property="video:actor" content="https://www.example.com/actors/john-doe"/>
//	<meta property="video:director" content="https://www.example.com/directors/jane-director"/>
//	<meta property="video:release_date" content="2024-09-15"/>
//	<meta property="video:series" content="https://www.example.com/video/series/example-series"/>
//	<meta property="video:episode" content="1"/>
type VideoEpisode struct {
	OpenGraphObject
	SeriesURL     string   // video:series, URL to the video series
	Duration      string   // video:duration, duration of the episode in seconds
	ActorURLs     []string // video:actor, URLs to the actors in the episode
	DirectorURL   string   // video:director, URL to the director of the episode
	ReleaseDate   string   // video:release_date, the release date of the episode
	EpisodeNumber int      // video:episode, the episode number in the series
}

// NewVideoEpisode initializes a VideoEpisode with the default type "video.episode".
func NewVideoEpisode(title, url, description, image, duration, seriesURL string, actorURLs []string, directorURL, releaseDate string, episodeNumber int) *VideoEpisode {
	videoEpisode := &VideoEpisode{
		OpenGraphObject: OpenGraphObject{
			Title:       title,
			URL:         url,
			Description: description,
			Image:       image,
		},
		SeriesURL:     seriesURL,
		Duration:      duration,
		ActorURLs:     actorURLs,
		DirectorURL:   directorURL,
		ReleaseDate:   releaseDate,
		EpisodeNumber: episodeNumber,
	}
	videoEpisode.ensureDefaults()
	return videoEpisode
}

// ToMetaTags generates the HTML meta tags for the Open Graph Video Episode as templ.Component.
func (ve *VideoEpisode) ToMetaTags() templ.Component {
	ve.ensureDefaults()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, tag := range ve.metaTags() {
			if tag.content != "" {
				if err := teseo.WriteMetaTag(w, tag.property, tag.content); err != nil {
					return err
				}
			}
		}

		// Write video:actor meta tags for each actor URL
		for _, actorURL := range ve.ActorURLs {
			if actorURL != "" {
				if err := teseo.WriteMetaTag(w, "video:actor", actorURL); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// ToGoHTMLMetaTags generates the HTML meta tags for the Open Graph Video Episode as a string for Go's `html/template`.
func (ve *VideoEpisode) ToGoHTMLMetaTags() string {
	ve.ensureDefaults()

	var sb strings.Builder

	for _, tag := range ve.metaTags() {
		if tag.content != "" {
			sb.WriteString(fmt.Sprintf(`<meta property="%s" content="%s"/>`, tag.property, tag.content))
			sb.WriteString("\n")
		}
	}

	// Write video:actor meta tags for each actor URL
	for _, actorURL := range ve.ActorURLs {
		if actorURL != "" {
			sb.WriteString(fmt.Sprintf(`<meta property="video:actor" content="%s"/>`, actorURL))
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// ensureDefaults sets default values for VideoEpisode.
func (ve *VideoEpisode) ensureDefaults() {
	ve.OpenGraphObject.ensureDefaults("video.episode")
}

// metaTags returns all meta tags for the VideoEpisode object, including OpenGraphObject fields and video episode-specific ones.
func (ve *VideoEpisode) metaTags() []struct {
	property string
	content  string
} {
	return []struct {
		property string
		content  string
	}{
		{"og:type", "video.episode"},
		{"og:title", ve.Title},
		{"og:url", ve.URL},
		{"og:description", ve.Description},
		{"og:image", ve.Image},
		{"video:duration", ve.Duration},
		{"video:director", ve.DirectorURL},
		{"video:release_date", ve.ReleaseDate},
		{"video:series", ve.SeriesURL},
		{"video:episode", fmt.Sprintf("%d", ve.EpisodeNumber)},
	}
}
