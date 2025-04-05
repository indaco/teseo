package opengraph

import (
	"context"
	"strings"
	"testing"
)

func TestNewMusicPlaylist_SetsFieldsAndDefaults(t *testing.T) {
	pl := NewMusicPlaylist(
		"Chill Mix",
		"https://example.com/playlist",
		"Relaxing tunes",
		"https://example.com/cover.jpg",
		[]string{"https://example.com/songs/1", "https://example.com/songs/2"},
		"3600",
	)

	if pl.Type != "music.playlist" {
		t.Errorf("expected Type to be 'music.playlist', got '%s'", pl.Type)
	}
	if len(pl.SongURLs) != 2 {
		t.Errorf("expected 2 songs, got %d", len(pl.SongURLs))
	}
	if pl.Duration != "3600" {
		t.Errorf("expected duration '3600', got '%s'", pl.Duration)
	}
}

func TestMusicPlaylist_ensureDefaults(t *testing.T) {
	pl := &MusicPlaylist{}
	pl.ensureDefaults()
	if pl.Type != "music.playlist" {
		t.Errorf("ensureDefaults did not set type to music.playlist")
	}
}

func TestMusicPlaylist_metaTags(t *testing.T) {
	pl := NewMusicPlaylist(
		"Workout Hits",
		"https://example.com/playlist",
		"Pump-up jams",
		"https://example.com/img.jpg",
		[]string{"https://example.com/songs/track1", "https://example.com/songs/track2"},
		"1800",
	)

	tags := pl.metaTags()

	assertTag := func(prop, expected string) {
		for _, tag := range tags {
			if tag.property == prop && tag.content == expected {
				return
			}
		}
		t.Errorf("missing or incorrect tag: %s=%s", prop, expected)
	}

	assertTag("og:type", "music.playlist")
	assertTag("og:title", "Workout Hits")
	assertTag("og:url", "https://example.com/playlist")
	assertTag("og:description", "Pump-up jams")
	assertTag("og:image", "https://example.com/img.jpg")
	assertTag("music:duration", "1800")
	assertTag("music:song", "https://example.com/songs/track1")
	assertTag("music:song", "https://example.com/songs/track2")
}

func TestMusicPlaylist_metaTags_EmptyValues(t *testing.T) {
	pl := &MusicPlaylist{
		OpenGraphObject: OpenGraphObject{Title: "Only Title"},
		SongURLs:        []string{"", "https://valid.com/song"},
	}

	tags := pl.metaTags()

	validSeen := false
	for _, tag := range tags {
		if tag.property == "music:song" {
			if tag.content == "" {
				t.Errorf("empty song URL should be skipped")
			} else {
				validSeen = true
			}
		}
	}
	if !validSeen {
		t.Errorf("valid song URL was not included")
	}
}

func TestMusicPlaylist_ToMetaTags_WriteError(t *testing.T) {
	pl := NewMusicPlaylist(
		"Focus Playlist",
		"https://example.com/focus",
		"Focus music",
		"https://example.com/focus.jpg",
		[]string{"https://example.com/songs/alpha"},
		"2400",
	)
	pl.ensureDefaults()

	writer := &failingWriter{}

	err := pl.ToMetaTags().Render(context.Background(), writer)
	if err == nil {
		t.Fatal("expected error but got nil")
	}

	expected := "failed to write og:type meta tag"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestMusicPlaylist_ToGoHTMLMetaTags_Render(t *testing.T) {
	pl := NewMusicPlaylist(
		"Focus Playlist",
		"https://example.com/focus",
		"Focus music",
		"https://example.com/focus.jpg",
		[]string{"https://example.com/songs/alpha"},
		"2400",
	)

	html, err := pl.ToGoHTMLMetaTags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := string(html)
	if !strings.Contains(out, `property="og:title"`) {
		t.Errorf("expected 'og:title' in rendered HTML")
	}
	if !strings.Contains(out, "Focus Playlist") {
		t.Errorf("expected playlist title in rendered HTML")
	}
	if !strings.Contains(out, `property="music:song"`) {
		t.Errorf("expected 'music:song' in rendered HTML")
	}
}
