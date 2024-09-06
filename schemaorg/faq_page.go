package schemaorg

import (
	"context"
	"io"
	"log"

	"github.com/a-h/templ"
	"github.com/indaco/teseo"
)

// FAQPage represents a Schema.org FAQPage object.
//
// Example usage:
//
// Pure struct usage:
//
//	faqPage := &schemaorg.FAQPage{
//		MainEntity: []schemaorg.Question{
//			{Question: "What is Schema.org?", Answer: &schemaorg.Answer{Answer: "Schema.org is a structured data vocabulary."}},
//		},
//	}
//
// Factory method usage:
//
//	faqPage := schemaorg.NewFAQPage(
//		[]schemaorg.Question{
//			{Question: "What is Schema.org?", Answer: &schemaorg.Answer{Answer: "Schema.org is a structured data vocabulary."}},
//		},
//	)
//
// // Rendering JSON-LD using templ:
//
//	templ Page() {
//		@faqPage.ToJsonLd()
//	}
//
// // Rendering JSON-LD as `template.HTML` value:
//
//	jsonLdHtml := faqPage.ToGoHTMLJsonLd()
//
// Expected output:
//
//	{
//		"@context": "https://schema.org",
//		"@type": "FAQPage",
//		"mainEntity": [
//			{
//				"@type": "Question",
//				"name": "What is Schema.org?",
//				"acceptedAnswer": {"@type": "Answer", "text": "Schema.org is a structured data vocabulary."}
//			}
//		]
//	}
type FAQPage struct {
	Context    string      `json:"@context"`
	Type       string      `json:"@type"`
	MainEntity []*Question `json:"mainEntity,omitempty"`
}

// Question represents a Schema.org Question object
type Question struct {
	Type           string  `json:"@type"`
	Name           string  `json:"name,omitempty"`
	AcceptedAnswer *Answer `json:"acceptedAnswer,omitempty"`
}

// Answer represents a Schema.org Answer object
type Answer struct {
	Type string `json:"@type"`
	Text string `json:"text,omitempty"`
}

// NewFAQPage initializes an FAQPage with default context and type.
func NewFAQPage(questions []*Question) *FAQPage {
	faqPage := &FAQPage{
		Context:    "https://schema.org",
		Type:       "FAQPage",
		MainEntity: questions,
	}
	return faqPage
}

// NewQuestion initializes a Question with default type.
func NewQuestion(name string, answer *Answer) *Question {
	question := &Question{
		Type:           "Question",
		Name:           name,
		AcceptedAnswer: answer,
	}
	return question
}

// NewAnswer initializes an Answer with default type.
func NewAnswer(text string) *Answer {
	answer := &Answer{
		Type: "Answer",
		Text: text,
	}
	return answer
}

// ToJsonLd converts the FAQPage struct to a JSON-LD `templ.Component`.
func (fp *FAQPage) ToJsonLd() templ.Component {
	fp.ensureDefaults()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		return templ.JSONScript(teseo.GenerateUniqueKey(), fp).WithType("application/ld+json").Render(ctx, w)
	})
}

// ToGoHTMLJsonLd renders the FAQPage struct as a string for Go's `html/template`.
func (fp *FAQPage) ToGoHTMLJsonLd() (string, error) {
	fp.ensureDefaults()

	// Create the templ component.
	templComponent := fp.ToJsonLd()

	// Render the templ component to a `template.HTML` value.
	html, err := templ.ToGoHTML(context.Background(), templComponent)
	if err != nil {
		log.Fatalf("failed to convert to html: %v", err)
	}

	return string(html), nil
}

// ensureDefaults sets default values for FAQPage, Question, and Answer if they are not already set.
func (fp *FAQPage) ensureDefaults() {
	if fp.Context == "" {
		fp.Context = "https://schema.org"
	}

	if fp.Type == "" {
		fp.Type = "FAQPage"
	}

	for _, q := range fp.MainEntity {
		q.ensureDefaults()
	}
}

func (q *Question) ensureDefaults() {
	if q.Type == "" {
		q.Type = "Question"
	}

	if q.AcceptedAnswer != nil {
		q.AcceptedAnswer.ensureDefaults()
	}
}

func (a *Answer) ensureDefaults() {
	if a.Type == "" {
		a.Type = "Answer"
	}
}
