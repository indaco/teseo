package opengraph

import (
	"context"
	"html/template"
	"io"
	"log"

	"github.com/a-h/templ"
	"github.com/indaco/teseo"
)

// MusicPlaylist represents the Open Graph music playlist metadata.
// For more details about the meaning of the properties see: https://ogp.me/#type_music.playlist
//
// Example usage:
//
// Pure struct usage:
//
//	// Create a music playlist using pure struct
//	musicPlaylist := &opengraph.MusicPlaylist{
//		OpenGraphObject: opengraph.OpenGraphObject{
//			Title:       "Example Playlist Title",
//			URL:         "https://www.example.com/music/playlist/example-playlist",
//			Description: "This is an example playlist description.",
//			Image:       "https://www.example.com/images/playlist.jpg",
//		},
//		SongURLs: []string{"https://www.example.com/musicians/jane-doe", "https://www.example.com/musicians/john-doe"},
//		Duration: "60",
//	}
//
// Factory method usage:
//
//	// Create a music playlist
//	musicPlaylist := opengraph.NewMusicPlaylist(
//		"Example Playlist Title",
//		"https://www.example.com/music/playlist/example-playlist",
//		"This is an example playlist description.",
//		"https://www.example.com/images/playlist.jpg",
//		[]string{"https://www.example.com/musicians/jane-doe", "https://www.example.com/musicians/john-doe"},
//		"60",
//	)
//
// // Rendering the HTML meta tags using templ:
//
//	templ Page() {
//		@musicPlaylist.ToMetaTgs()
//	}
//
// // Rendering the HTML meta tags as `template.HTML` value:
//
//	metaTagsHtml := musicAlbum.ToGoHTMLMetaTgs()
//
// Expected output:
//
//	<meta property="og:type" content="music.playlist"/>
//	<meta property="og:title" content="Example Playlist Title"/>
//	<meta property="og:url" content="https://www.example.com/music/playlist/example-playlist"/>
//	<meta property="og:description" content="This is an example playlist description."/>
//	<meta property="og:image" content="https://www.example.com/images/playlist.jpg"/>
//	<meta property="music:song" content="https://www.example.com/musicians/jane-doe"/>
//	<meta property="music:song" content="https://www.example.com/musicians/john-doe"/>
//	<meta property="music:duration" content="60"/>
type MusicPlaylist struct {
	OpenGraphObject
	SongURLs []string // music:song, URLs to the songs in the playlist
	Duration string   // music:duration, duration of the playlist in seconds
}

// NewMusicPlaylist initializes a MusicPlaylist with the default type "music.playlist".
func NewMusicPlaylist(title, url, description, image string, songURLs []string, duration string) *MusicPlaylist {
	musicPlaylist := &MusicPlaylist{
		OpenGraphObject: OpenGraphObject{
			Title:       title,
			URL:         url,
			Description: description,
			Image:       image,
		},
		SongURLs: songURLs,
		Duration: duration,
	}
	musicPlaylist.ensureDefaults()
	return musicPlaylist
}

// ToMetaTags generates the HTML meta tags for the Open Graph Music Playlist as templ.Component.
func (mp *MusicPlaylist) ToMetaTags() templ.Component {
	mp.ensureDefaults()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, tag := range mp.metaTags() {
			if tag.content != "" {
				if err := teseo.WriteMetaTag(w, tag.property, tag.content); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// ToGoHTMLMetaTags generates the HTML meta tags for the Open Graph Music Playlist as `template.HTML` value for Go's `html/template`.
func (mp *MusicPlaylist) ToGoHTMLMetaTags() (template.HTML, error) {
	// Create the templ component.
	templComponent := mp.ToMetaTags()

	// Render the templ component to a `template.HTML` value.
	html, err := templ.ToGoHTML(context.Background(), templComponent)
	if err != nil {
		log.Fatalf("failed to convert to html: %v", err)
	}

	return html, nil
}

// ensureDefaults sets default values for MusicPlaylist.
func (mp *MusicPlaylist) ensureDefaults() {
	mp.OpenGraphObject.ensureDefaults("music.playlist")
}

// metaTags returns all meta tags for the MusicPlaylist object, including OpenGraphObject fields and music-specific ones.
func (mp *MusicPlaylist) metaTags() []struct {
	property string
	content  string
} {
	tags := []struct {
		property string
		content  string
	}{
		{"og:type", "music.playlist"},
		{"og:title", mp.Title},
		{"og:url", mp.URL},
		{"og:description", mp.Description},
		{"og:image", mp.Image},
		{"music:duration", mp.Duration},
	}

	// Add music:song tags for each song URL
	for _, songURL := range mp.SongURLs {
		if songURL != "" {
			tags = append(tags, struct {
				property string
				content  string
			}{"music:song", songURL})
		}
	}

	return tags
}
