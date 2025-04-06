package schemaorg

import (
	"slices"
	"testing"
)

func TestNewFAQPage_Defaults(t *testing.T) {
	q := NewQuestion("What is Go?", NewAnswer("A programming language"))
	faq := NewFAQPage([]*Question{q})
	if faq.Context != "https://schema.org" {
		t.Errorf("expected default context")
	}
	if faq.Type != "FAQPage" {
		t.Errorf("expected default type")
	}
}

func TestFAQPage_EnsureDefaults(t *testing.T) {
	q := &Question{
		AcceptedAnswer: &Answer{},
	}
	fp := &FAQPage{
		MainEntity: []*Question{q},
	}
	fp.ensureDefaults()

	if fp.Context != "https://schema.org" {
		t.Errorf("expected default context")
	}
	if fp.Type != "FAQPage" {
		t.Errorf("expected default type")
	}
	if q.Type != "Question" {
		t.Errorf("expected default type for Question")
	}
	if q.AcceptedAnswer.Type != "Answer" {
		t.Errorf("expected default type for Answer")
	}
}

func TestFAQPage_ToGoHTMLJsonLd(t *testing.T) {
	fp := NewFAQPage([]*Question{
		NewQuestion("Q", NewAnswer("A")),
	})
	html, err := fp.ToGoHTMLJsonLd()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if html == "" {
		t.Errorf("expected non-empty HTML")
	}
}

func TestFAQPage_Validate(t *testing.T) {
	tests := []struct {
		name     string
		page     *FAQPage
		expected []string
	}{
		{
			name: "valid page",
			page: NewFAQPage([]*Question{
				NewQuestion("Q", NewAnswer("A")),
			}),
			expected: nil,
		},
		{
			name: "missing main entity",
			page: &FAQPage{},
			expected: []string{
				"FAQPage should contain at least one question",
			},
		},
		{
			name: "missing question name",
			page: NewFAQPage([]*Question{
				NewQuestion("", NewAnswer("A")),
			}),
			expected: []string{
				"Question 1 is missing a name",
			},
		},
		{
			name: "missing accepted answer",
			page: NewFAQPage([]*Question{
				{Name: "Q"},
			}),
			expected: []string{
				"Question 1 is missing an accepted answer",
			},
		},
		{
			name: "empty answer text",
			page: NewFAQPage([]*Question{
				NewQuestion("Q", &Answer{Text: ""}),
			}),
			expected: []string{
				"Answer for question 1 is missing text",
			},
		},
		{
			name: "multiple warnings",
			page: NewFAQPage([]*Question{
				{},
			}),
			expected: []string{
				"Question 1 is missing a name",
				"Question 1 is missing an accepted answer",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings := tt.page.Validate()
			if len(warnings) != len(tt.expected) {
				t.Errorf("expected %d warnings, got %d: %v", len(tt.expected), len(warnings), warnings)
				return
			}
			for _, expected := range tt.expected {
				found := slices.Contains(warnings, expected)
				if !found {
					t.Errorf("missing expected warning: %s", expected)
				}
			}
		})
	}
}
