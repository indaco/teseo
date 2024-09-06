package opengraph

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/a-h/templ"
	"github.com/indaco/teseo"
)

// MusicSong represents the Open Graph music song metadata.
//
// Example usage:
//
// Pure struct usage:
//
//	// Create a music song using pure struct
//	musicSong := &opengraph.MusicSong{
//		OpenGraphObject: opengraph.OpenGraphObject{
//			Title:       "Example Song Title",
//			URL:         "https://www.example.com/music/song/example-song",
//			Description: "This is an example song description.",
//			Image:       "https://www.example.com/images/song.jpg",
//		},
//		Duration: "240", // Duration in seconds
//		AlbumURL: "https://www.example.com/music/album/example-album",
//		MusicianURLs: []string{
//			"https://www.example.com/musicians/jane-doe",
//			"https://www.example.com/musicians/john-doe",
//		},
//	}
//
// Factory method usage:
//
//	// Create a music song using the factory method
//	musicSong := opengraph.NewMusicSong(
//		"Example Song Title",
//		"https://www.example.com/music/song/example-song",
//		"This is an example song description.",
//		"https://www.example.com/images/song.jpg",
//		"240", // Duration in seconds
//		"https://www.example.com/music/album/example-album",
//		[]string{"https://www.example.com/musicians/jane-doe", "https://www.example.com/musicians/john-doe"},
//	)
//
// // Rendering the HTML meta tags using templ:
//
//	templ Page() {
//		@musicSong.ToMetaTgs()
//	}
//
// // Rendering the HTML meta tags as `template.HTML` value:
//
//	metaTagsHtml := musicSong.ToGoHTMLMetaTgs()
//
// Expected output:
//
//	<meta property="og:type" content="music.song"/>
//	<meta property="og:title" content="Example Song Title"/>
//	<meta property="og:url" content="https://www.example.com/music/song/example-song"/>
//	<meta property="og:description" content="This is an example song description."/>
//	<meta property="og:image" content="https://www.example.com/images/song.jpg"/>
//	<meta property="music:duration" content="240"/>
//	<meta property="music:album" content="https://www.example.com/music/album/example-album"/>
//	<meta property="music:musician" content="https://www.example.com/musicians/jane-doe"/>
//	<meta property="music:musician" content="https://www.example.com/musicians/john-doe"/>
type MusicSong struct {
	OpenGraphObject
	Duration     string   // music:duration, duration of the song in seconds
	AlbumURL     string   // music:album, URL to the album
	MusicianURLs []string // music:musician, URLs to the musicians
}

// NewMusicSong initializes a MusicSong with the default type "music.song".
func NewMusicSong(title, url, description, image, duration, albumURL string, musicianURLs []string) *MusicSong {
	musicSong := &MusicSong{
		OpenGraphObject: OpenGraphObject{
			Title:       title,
			URL:         url,
			Description: description,
			Image:       image,
		},
		Duration:     duration,
		AlbumURL:     albumURL,
		MusicianURLs: musicianURLs,
	}
	musicSong.ensureDefaults()
	return musicSong
}

// ToMetaTags generates the HTML meta tags for the Open Graph Music Song as templ.Component.
func (ms *MusicSong) ToMetaTags() templ.Component {
	ms.ensureDefaults()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, tag := range ms.metaTags() {
			if tag.content != "" {
				if err := teseo.WriteMetaTag(w, tag.property, tag.content); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// ToGoHTMLMetaTags generates the HTML meta tags for the Open Graph Music Song as a string for Go's `html/template`.
func (ms *MusicSong) ToGoHTMLMetaTags() string {
	ms.ensureDefaults()

	var sb strings.Builder

	for _, tag := range ms.metaTags() {
		if tag.content != "" {
			sb.WriteString(fmt.Sprintf(`<meta property="%s" content="%s"/>`, tag.property, tag.content))
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// ensureDefaults sets default values for MusicSong.
func (ms *MusicSong) ensureDefaults() {
	ms.OpenGraphObject.ensureDefaults("music.song")
}

// metaTags returns all meta tags for the MusicSong object, including OpenGraphObject fields and music-specific ones.
func (ms *MusicSong) metaTags() []struct {
	property string
	content  string
} {
	tags := []struct {
		property string
		content  string
	}{
		{"og:type", "music.song"},
		{"og:title", ms.Title},
		{"og:url", ms.URL},
		{"og:description", ms.Description},
		{"og:image", ms.Image},
		{"music:duration", ms.Duration},
		{"music:album", ms.AlbumURL},
	}

	// Add music:musician tags for each musician URL
	for _, musicianURL := range ms.MusicianURLs {
		if musicianURL != "" {
			tags = append(tags, struct {
				property string
				content  string
			}{"music:musician", musicianURL})
		}
	}

	return tags
}
