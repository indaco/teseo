package opengraph

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/a-h/templ"
	"github.com/indaco/teseo"
)

// MusicAlbum represents the Open Graph music album metadata.
//
// Example usage:
//
// Pure struct usage:
//
//	// Create a music album using pure struct
//	musicAlbum := &opengraph.MusicAlbum{
//		OpenGraphObject: opengraph.OpenGraphObject{
//			Title:       "Example Album Title",
//			URL:         "https://www.example.com/music/album/example-album",
//			Description: "This is an example album description.",
//			Image:       "https://www.example.com/images/album.jpg",
//		},
//		Musician:    []string{"https://www.example.com/musicians/jane-doe", "https://www.example.com/musicians/john-doe"},
//		ReleaseDate: "2024-09-15",
//		Genre:       "Rock",
//	}
//
// Factory method usage:
//
//	// Create a music album
//	musicAlbum := opengraph.NewMusicAlbum(
//		"Example Album Title",
//		"https://www.example.com/music/album/example-album",
//		"This is an example album description.",
//		"https://www.example.com/images/album.jpg",
//		"2024-09-15",
//		"Rock",
//		[]string{"https://www.example.com/musicians/jane-doe", "https://www.example.com/musicians/john-doe"},
//	)
//
// // Rendering the HTML meta tags using templ:
//
//	templ Page() {
//		@musicAlbum.ToMetaTgs()
//	}
//
// // Rendering the HTML meta tags as `template.HTML` value:
//
//	metaTagsHtml := musicAlbum.ToGoHTMLMetaTgs()
//
// Expected output:
//
//	<meta property="og:type" content="music.album"/>
//	<meta property="og:title" content="Example Album Title"/>
//	<meta property="og:url" content="https://www.example.com/music/album/example-album"/>
//	<meta property="og:description" content="This is an example album description."/>
//	<meta property="og:image" content="https://www.example.com/images/album.jpg"/>
//	<meta property="music:release_date" content="2024-09-15"/>
//	<meta property="music:genre" content="Rock"/>
//	<meta property="music:musician" content="https://www.example.com/musicians/jane-doe"/>
//	<meta property="music:musician" content="https://www.example.com/musicians/john-doe"/>
type MusicAlbum struct {
	OpenGraphObject
	Musician    []string // music:musician, URLs to the musicians in the album
	ReleaseDate string   // music:release_date, the release date of the album
	Genre       string   // music:genre, genre of the album
}

// NewMusicAlbum initializes a MusicAlbum with the default type "music.album".
func NewMusicAlbum(title, url, description, image, releaseDate, genre string, musician []string) *MusicAlbum {
	musicAlbum := &MusicAlbum{
		OpenGraphObject: OpenGraphObject{
			Title:       title,
			URL:         url,
			Description: description,
			Image:       image,
		},
		Musician:    musician,
		ReleaseDate: releaseDate,
		Genre:       genre,
	}
	musicAlbum.ensureDefaults()
	return musicAlbum
}

// ToMetaTags generates the HTML meta tags for the Open Graph Music Album as templ.Component.
func (ma *MusicAlbum) ToMetaTags() templ.Component {
	ma.ensureDefaults()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, tag := range ma.metaTags() {
			if tag.content != "" {
				if err := teseo.WriteMetaTag(w, tag.property, tag.content); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// ToGoHTMLMetaTags generates the HTML meta tags for the Open Graph Music Album as a string for Go's `html/template`.
func (ma *MusicAlbum) ToGoHTMLMetaTags() string {
	ma.ensureDefaults()

	var sb strings.Builder

	for _, tag := range ma.metaTags() {
		if tag.content != "" {
			sb.WriteString(fmt.Sprintf(`<meta property="%s" content="%s"/>`, tag.property, tag.content))
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// ensureDefaults sets default values for MusicAlbum.
func (ma *MusicAlbum) ensureDefaults() {
	ma.OpenGraphObject.ensureDefaults("music.album")
}

// metaTags returns all meta tags for the MusicAlbum object, including OpenGraphObject fields and music-specific ones.
func (ma *MusicAlbum) metaTags() []struct {
	property string
	content  string
} {
	tags := []struct {
		property string
		content  string
	}{
		{"og:type", "music.album"},
		{"og:title", ma.Title},
		{"og:url", ma.URL},
		{"og:description", ma.Description},
		{"og:image", ma.Image},
		{"music:release_date", ma.ReleaseDate},
		{"music:genre", ma.Genre},
	}

	// Add music:musician tags
	for _, musician := range ma.Musician {
		if musician != "" {
			tags = append(tags, struct {
				property string
				content  string
			}{"music:musician", musician})
		}
	}

	return tags
}
