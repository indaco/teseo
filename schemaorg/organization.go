package schemaorg

// Organization represents a Schema.org Organization object.
// For more details about the meaning of the properties, see:https://schema.org/Organization
//
// Example usage:
//
// Pure struct usage:
//
// 	organization := &schemaorg.Organization{
// 		Name:        "Example Organization",
// 		URL:         "https://www.example.com",
// 		Logo:        "https://www.example.com/logo.jpg",
// 		Description: "This is an example organization.",
// 	}
//
// Factory method usage:
//
// 	organization := schemaorg.NewOrganization(
// 		"Example Organization",
// 		"https://www.example.com",
// 		"https://www.example.com/logo.jpg",
// 		"This is an example organization",
// 	)
//
// // Rendering JSON-LD using templ:
//
//	templ Page() {
//		@organization.ToJsonLd()
//	}
//
// // Rendering JSON-LD as `template.HTML` value:
//
//	jsonLdHtml := organization.ToGoHTMLJsonLd()
//
// Expected output:
//
// 	{
// 		"@context": "https://schema.org",
// 		"@type": "Organization",
// 		"name": "Example Organization",
// 		"url": "https://www.example.com",
// 		"logo": "https://www.example.com/logo.jpg",
// 		"description": "This is an example organization"
// 	}

// NewOrganization initializes an Organization with default context and type.
func NewOrganization(name string, url string, logoURL string, contactPoints []ContactPoint, sameAs []string) *Organization {
	org := &Organization{
		Name: name,
		URL:  url,
		Logo: &ImageObject{
			Type: "ImageObject",
			URL:  logoURL,
		},
		ContactPoints: contactPoints,
		SameAs:        sameAs,
	}
	org.ensureDefaults()
	return org
}
