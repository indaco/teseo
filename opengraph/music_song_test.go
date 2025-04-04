package opengraph

import (
	"strings"
	"testing"
)

func TestNewMusicSong_SetsFieldsAndDefaults(t *testing.T) {
	song := NewMusicSong(
		"Test Song",
		"https://example.com/song",
		"Example description",
		"https://example.com/image.jpg",
		"300",
		"https://example.com/album",
		[]string{"https://example.com/musicians/jane", "https://example.com/musicians/john"},
	)

	if song.Type != "music.song" {
		t.Errorf("expected Type to be 'music.song', got '%s'", song.Type)
	}
	if song.AlbumURL != "https://example.com/album" {
		t.Errorf("expected AlbumURL to be set")
	}
	if len(song.MusicianURLs) != 2 {
		t.Errorf("expected 2 musicians, got %d", len(song.MusicianURLs))
	}
}

func TestMusicSong_ensureDefaults(t *testing.T) {
	ms := &MusicSong{}
	ms.ensureDefaults()
	if ms.Type != "music.song" {
		t.Errorf("expected default type to be 'music.song'")
	}
}

func TestMusicSong_metaTags(t *testing.T) {
	song := NewMusicSong(
		"My Song",
		"https://example.com/song",
		"Some song description",
		"https://example.com/image.jpg",
		"210",
		"https://example.com/album",
		[]string{"https://example.com/m/jane", "https://example.com/m/john"},
	)

	tags := song.metaTags()

	assertTag := func(prop, expected string) {
		for _, tag := range tags {
			if tag.property == prop && tag.content == expected {
				return
			}
		}
		t.Errorf("missing or incorrect tag: %s=%s", prop, expected)
	}

	assertTag("og:type", "music.song")
	assertTag("og:title", "My Song")
	assertTag("og:url", "https://example.com/song")
	assertTag("og:description", "Some song description")
	assertTag("og:image", "https://example.com/image.jpg")
	assertTag("music:duration", "210")
	assertTag("music:album", "https://example.com/album")
	assertTag("music:musician", "https://example.com/m/jane")
	assertTag("music:musician", "https://example.com/m/john")
}

func TestMusicSong_metaTags_EmptyMusicianFiltered(t *testing.T) {
	ms := &MusicSong{
		OpenGraphObject: OpenGraphObject{Title: "No Empty Musician"},
		MusicianURLs:    []string{"", "https://example.com/m/john"},
	}

	tags := ms.metaTags()

	musicianFound := false
	for _, tag := range tags {
		if tag.property == "music:musician" {
			if tag.content == "" {
				t.Errorf("empty musician should not be rendered")
			} else {
				musicianFound = true
			}
		}
	}
	if !musicianFound {
		t.Errorf("valid musician not rendered")
	}
}

func TestMusicSong_ToGoHTMLMetaTags_Render(t *testing.T) {
	song := NewMusicSong(
		"Final Tune",
		"https://example.com/final",
		"The closing track",
		"https://example.com/final.jpg",
		"240",
		"https://example.com/album/final",
		[]string{"https://example.com/musician"},
	)

	html, err := song.ToGoHTMLMetaTags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := string(html)
	if !strings.Contains(out, `property="og:title"`) {
		t.Errorf("expected 'og:title' in rendered HTML")
	}
	if !strings.Contains(out, `Final Tune`) {
		t.Errorf("expected song title in rendered HTML")
	}
	if !strings.Contains(out, `property="music:musician"`) {
		t.Errorf("expected 'music:musician' tag in rendered HTML")
	}
}
