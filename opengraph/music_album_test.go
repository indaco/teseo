package opengraph

import (
	"context"
	"strings"
	"testing"
)

func TestNewMusicAlbum_SetsFieldsAndDefaults(t *testing.T) {
	album := NewMusicAlbum(
		"Greatest Hits",
		"https://example.com/album",
		"A classic compilation",
		"https://example.com/cover.jpg",
		"2024-09-15",
		"Rock",
		[]string{"https://example.com/musicians/jane", "https://example.com/musicians/john"},
	)

	if album.Type != "music.album" {
		t.Errorf("expected Type to be 'music.album', got '%s'", album.Type)
	}
	if album.Genre != "Rock" {
		t.Errorf("expected Genre to be 'Rock'")
	}
	if len(album.Musician) != 2 {
		t.Errorf("expected 2 musicians, got %d", len(album.Musician))
	}
}

func TestMusicAlbum_ensureDefaults(t *testing.T) {
	a := &MusicAlbum{}
	a.ensureDefaults()
	if a.Type != "music.album" {
		t.Errorf("expected default type to be 'music.album', got '%s'", a.Type)
	}
}

func TestMusicAlbum_metaTags(t *testing.T) {
	album := NewMusicAlbum(
		"Greatest Hits",
		"https://example.com/album",
		"A classic compilation",
		"https://example.com/cover.jpg",
		"2024-09-15",
		"Rock",
		[]string{"https://example.com/musicians/jane", "https://example.com/musicians/john"},
	)

	tags := album.metaTags()

	assertTag := func(prop, expected string) {
		for _, tag := range tags {
			if tag.property == prop && tag.content == expected {
				return
			}
		}
		t.Errorf("missing or incorrect tag: %s=%s", prop, expected)
	}

	assertTag("og:type", "music.album")
	assertTag("og:title", "Greatest Hits")
	assertTag("og:url", "https://example.com/album")
	assertTag("og:description", "A classic compilation")
	assertTag("og:image", "https://example.com/cover.jpg")
	assertTag("music:release_date", "2024-09-15")
	assertTag("music:genre", "Rock")
	assertTag("music:musician", "https://example.com/musicians/jane")
	assertTag("music:musician", "https://example.com/musicians/john")
}

func TestMusicAlbum_metaTags_EmptyValues(t *testing.T) {
	album := &MusicAlbum{
		OpenGraphObject: OpenGraphObject{Title: "Empty Fields"},
		Musician:        []string{"", "https://example.com/musicians/john"},
	}

	tags := album.metaTags()

	musicianFound := false
	for _, tag := range tags {
		if tag.property == "music:musician" {
			if tag.content == "" {
				t.Errorf("empty musician tag should not be included")
			} else {
				musicianFound = true
			}
		}
	}

	if !musicianFound {
		t.Errorf("valid musician tag was not found")
	}
}

func TestMusicAlbum_ToMetaTags_WriteError(t *testing.T) {
	album := NewMusicAlbum(
		"Render Test Album",
		"https://example.com/render",
		"Render test description",
		"https://example.com/cover.png",
		"2024-10-01",
		"Indie",
		[]string{"https://example.com/musicians/alex"},
	)
	album.ensureDefaults()

	writer := &failingWriter{}

	err := album.ToMetaTags().Render(context.Background(), writer)
	if err == nil {
		t.Fatal("expected error but got nil")
	}

	expected := "failed to write og:type meta tag"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestMusicAlbum_ToGoHTMLMetaTags_Render(t *testing.T) {
	album := NewMusicAlbum(
		"Render Test Album",
		"https://example.com/render",
		"Render test description",
		"https://example.com/cover.png",
		"2024-10-01",
		"Indie",
		[]string{"https://example.com/musicians/alex"},
	)

	html, err := album.ToGoHTMLMetaTags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := string(html)
	if !strings.Contains(out, `property="og:title"`) {
		t.Errorf("expected 'og:title' in rendered HTML")
	}
	if !strings.Contains(out, "Render Test Album") {
		t.Errorf("expected album title in rendered HTML")
	}
	if !strings.Contains(out, `property="music:musician"`) {
		t.Errorf("expected 'music:musician' in rendered HTML")
	}
}
