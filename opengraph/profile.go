package opengraph

import (
	"context"
	"html/template"
	"io"
	"log"

	"github.com/a-h/templ"
	"github.com/indaco/teseo"
)

// Profile represents the Open Graph profile metadata.
//
// Example usage:
//
// Pure struct usage:
//
//	// Create a profile using pure struct
//	profile := &opengraph.Profile{
//		OpenGraphObject: opengraph.OpenGraphObject{
//			Title:       "John Doe",
//			URL:         "https://www.example.com/profile/johndoe",
//			Description: "This is John Doe's profile.",
//			Image:       "https://www.example.com/images/profile.jpg",
//		},
//		FirstName: "John",
//		LastName:  "Doe",
//		Username:  "johndoe",
//		Gender:    "male",
//	}
//
// Factory method usage:
//
//	// Create a profile
//	profile := opengraph.NewProfile(
//		"John Doe",
//		"John",
//		"Doe",
//		"johndoe",
//		"male",
//		"https://www.example.com/profile/johndoe",
//		"This is John Doe's profile.",
//		"https://www.example.com/images/profile.jpg",
//	)
//
// // Rendering the HTML meta tags using templ:
//
//	templ Page() {
//		@profile.ToMetaTgs()
//	}
//
// // Rendering the HTML meta tags as `template.HTML` value:
//
//	metaTagsHtml := profile.ToGoHTMLMetaTgs()
//
// Expected output:
//
//	<meta property="og:type" content="profile"/>
//	<meta property="og:title" content="John Doe"/>
//	<meta property="profile:first_name" content="John"/>
//	<meta property="profile:last_name" content="Doe"/>
//	<meta property="profile:username" content="johndoe"/>
//	<meta property="profile:gender" content="male"/>
//	<meta property="og:url" content="https://www.example.com/profile/johndoe"/>
//	<meta property="og:description" content="This is John Doe's profile."/>
//	<meta property="og:image" content="https://www.example.com/images/profile.jpg"/>
type Profile struct {
	OpenGraphObject
	FirstName string // profile:first_name, first name
	LastName  string // profile:last_name, last name
	Username  string // profile:username, username
	Gender    string // profile:gender, gender
}

// NewProfile initializes an OpenGraphProfile with the default type "profile".
func NewProfile(title string, firstName string, lastName string, username string, gender string, url string, description string, image string) *Profile {
	profile := &Profile{
		OpenGraphObject: OpenGraphObject{
			Title:       title,
			URL:         url,
			Description: description,
			Image:       image,
		},
		FirstName: firstName,
		LastName:  lastName,
		Username:  username,
		Gender:    gender,
	}
	profile.ensureDefaults()
	return profile
}

// ToMetaTags generates the HTML meta tags for the Open Graph Profile as templ.Component.
func (p *Profile) ToMetaTags() templ.Component {
	p.ensureDefaults()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for _, tag := range p.metaTags() {
			if tag.content != "" {
				if err := teseo.WriteMetaTag(w, tag.property, tag.content); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// ToGoHTMLMetaTags generates the HTML meta tags for the Open Graph Profile as `template.HTML` value for Go's `html/template`.
func (p *Profile) ToGoHTMLMetaTags() (template.HTML, error) {
	// Create the templ component.
	templComponent := p.ToMetaTags()

	// Render the templ component to a `template.HTML` value.
	html, err := templ.ToGoHTML(context.Background(), templComponent)
	if err != nil {
		log.Fatalf("failed to convert to html: %v", err)
	}

	return html, nil
}

// ensureDefaults sets default values for Profile.
func (p *Profile) ensureDefaults() {
	p.OpenGraphObject.ensureDefaults("profile")
}

// metaTags returns all meta tags for the Profile object, including OpenGraphObject fields and profile-specific ones.
func (p *Profile) metaTags() []struct {
	property string
	content  string
} {
	return []struct {
		property string
		content  string
	}{
		{"og:type", "profile"},
		{"og:title", p.Title},
		{"og:url", p.URL},
		{"profile:first_name", p.FirstName},
		{"profile:last_name", p.LastName},
		{"profile:username", p.Username},
		{"profile:gender", p.Gender},
		{"og:description", p.Description},
		{"og:image", p.Image},
	}
}
