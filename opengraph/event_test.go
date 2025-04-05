package opengraph

import (
	"context"
	"strings"
	"testing"
)

func TestNewEvent_SetsFieldsAndDefaults(t *testing.T) {
	e := NewEvent(
		"Event Title",
		"https://example.com/event",
		"An event desc",
		"https://example.com/event.jpg",
		"2024-01-01T09:00:00Z",
		"2024-01-01T17:00:00Z",
		"Convention Center",
	)

	if e.Type != "event" {
		t.Errorf("expected Type to be 'event', got '%s'", e.Type)
	}
	if e.StartDate != "2024-01-01T09:00:00Z" {
		t.Errorf("expected StartDate to be set")
	}
}

func TestEvent_ensureDefaults(t *testing.T) {
	e := &Event{}
	e.ensureDefaults()
	if e.Type != "event" {
		t.Errorf("ensureDefaults did not set type to event")
	}
}

func TestEvent_metaTags(t *testing.T) {
	e := NewEvent(
		"Tech Conf",
		"https://conf.com",
		"Join us for a day of tech",
		"https://conf.com/img.jpg",
		"2024-05-10T08:00:00Z",
		"2024-05-10T17:00:00Z",
		"Main Hall",
	)

	tags := e.metaTags()

	assertTag := func(prop, expected string) {
		for _, tag := range tags {
			if tag.property == prop && tag.content == expected {
				return
			}
		}
		t.Errorf("missing or incorrect tag: %s=%s", prop, expected)
	}

	assertTag("og:type", "event")
	assertTag("og:title", "Tech Conf")
	assertTag("og:url", "https://conf.com")
	assertTag("og:description", "Join us for a day of tech")
	assertTag("og:image", "https://conf.com/img.jpg")
	assertTag("event:start_date", "2024-05-10T08:00:00Z")
	assertTag("event:end_date", "2024-05-10T17:00:00Z")
	assertTag("event:location", "Main Hall")
}

func TestEvent_metaTags_EmptyValues(t *testing.T) {
	e := &Event{
		OpenGraphObject: OpenGraphObject{Title: "Minimal"},
		StartDate:       "",
		EndDate:         "2024-12-31T23:59:59Z",
		Location:        "",
	}

	tags := e.metaTags()

	startSeen := false
	endSeen := false
	locSeen := false

	for _, tag := range tags {
		if tag.property == "event:start_date" {
			startSeen = true
			if tag.content != "" {
				t.Errorf("expected empty start_date")
			}
		}
		if tag.property == "event:end_date" {
			endSeen = true
			if tag.content != "2024-12-31T23:59:59Z" {
				t.Errorf("incorrect end_date content")
			}
		}
		if tag.property == "event:location" {
			locSeen = true
			if tag.content != "" {
				t.Errorf("expected empty location")
			}
		}
	}

	if !startSeen || !endSeen || !locSeen {
		t.Errorf("some tags were missing from metaTags")
	}
}

func TestEvent_ToMetaTags_WriteError(t *testing.T) {
	e := NewEvent(
		"Dev Meetup",
		"https://dev.com/meetup",
		"A developer event",
		"https://dev.com/img.png",
		"2024-07-01T09:00:00Z",
		"2024-07-01T17:00:00Z",
		"City Hub",
	)
	e.ensureDefaults()

	writer := &failingWriter{}

	err := e.ToMetaTags().Render(context.Background(), writer)
	if err == nil {
		t.Fatal("expected error but got nil")
	}

	expected := "failed to write og:type meta tag"
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestEvent_ToGoHTMLMetaTags_Render(t *testing.T) {
	e := NewEvent(
		"Dev Meetup",
		"https://dev.com/meetup",
		"A developer event",
		"https://dev.com/img.png",
		"2024-07-01T09:00:00Z",
		"2024-07-01T17:00:00Z",
		"City Hub",
	)

	html, err := e.ToGoHTMLMetaTags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := string(html)

	if !strings.Contains(out, `property="og:title"`) {
		t.Errorf("expected 'og:title' in output HTML")
	}
	if !strings.Contains(out, `property="event:start_date"`) {
		t.Errorf("expected 'event:start_date' in output HTML")
	}
	if !strings.Contains(out, "Dev Meetup") {
		t.Errorf("expected content value to include event title")
	}
}
