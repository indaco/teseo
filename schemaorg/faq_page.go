package schemaorg

import (
	"fmt"
	"html/template"

	"github.com/a-h/templ"
	"github.com/indaco/teseo"
)

// FAQPage represents a Schema.org FAQPage object.
// For more details about the meaning of the properties see:https://schema.org/FAQPage
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

func (fp *FAQPage) Validate() []string {
	var warnings []string
	if len(fp.MainEntity) == 0 {
		warnings = append(warnings, "FAQPage should contain at least one question")
	}
	for i, q := range fp.MainEntity {
		if q.Name == "" {
			warnings = append(warnings, fmt.Sprintf("Question %d is missing a name", i+1))
		}
		if q.AcceptedAnswer == nil {
			warnings = append(warnings, fmt.Sprintf("Question %d is missing an accepted answer", i+1))
		} else if q.AcceptedAnswer.Text == "" {
			warnings = append(warnings, fmt.Sprintf("Answer for question %d is missing text", i+1))
		}
	}
	return warnings
}

// ToJsonLd converts the FAQPage struct to a JSON-LD `templ.Component`.
func (fp *FAQPage) ToJsonLd() templ.Component {
	fp.ensureDefaults()
	id := fmt.Sprintf("%s-%s", "faqpage", teseo.GenerateUniqueKey())
	return templ.JSONScript(id, fp).WithType("application/ld+json")
}

// ToGoHTMLJsonLd renders the FAQPage struct as`template.HTML` value for Go's `html/template`.
func (fp *FAQPage) ToGoHTMLJsonLd() (template.HTML, error) {
	return teseo.RenderToHTML(fp.ToJsonLd())
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
