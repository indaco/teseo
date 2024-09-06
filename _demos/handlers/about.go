package handlers

import (
	"log"
	"net/http"

	"github.com/indaco/teseo"
	"github.com/indaco/teseo/_demos/pages"
	"github.com/indaco/teseo/_demos/types"
	"github.com/indaco/teseo/opengraph"
	"github.com/indaco/teseo/schemaorg"
	"github.com/indaco/teseo/twittercard"
)

func HandleAbout(w http.ResponseWriter, r *http.Request) {

	pageURL := teseo.GetFullURL(r)
	bcl, err := schemaorg.NewBreadcrumbListFromUrl(pageURL)
	if err != nil {
		log.Fatalf("Error generating breadcrumb list: %v", err)
		return
	}

	// Create a new SiteNavigationElement struct to be populated by reading the sitemap.xml file.
	sne := &schemaorg.SiteNavigationElement{}
	err = sne.FromSitemapFile("./_demos/statics/sitemap.xml")
	if err != nil {
		log.Fatalf("Failed to read sitemap: %v", err)
	}

	headerItems := &types.SEOItems{
		WebPage: &schemaorg.WebPage{
			URL:           "https://www.example.com/about",
			Name:          "About Us",
			Headline:      "Learn More About Our Company",
			Description:   "This is the about us page of Example.com, providing information about our company, our mission, and our values.",
			About:         "Company Information",
			Keywords:      "about us, company, mission, values",
			InLanguage:    "en",
			IsPartOf:      "https://www.example.com",
			LastReviewed:  "2024-09-01",
			PrimaryImage:  "https://www.example.com/images/about-us.jpg",
			DatePublished: "2020-01-01",
			DateModified:  "2024-09-01",
		},
		SiteNavElement: sne,

		Profile: &opengraph.Profile{
			OpenGraphObject: opengraph.OpenGraphObject{
				Title:       "Jane Doe",
				URL:         "https://www.example.com/janedoe",
				Description: "This is Jane Doe's profile on Example.com.",
				Image:       "https://placehold.co/600x400?text=JD",
			},
			FirstName: "Jane",
			LastName:  "Doe",
			Username:  "janedoe123",
			Gender:    "female",
		},

		Breadcrumb: bcl,
		TwitterCard: twittercard.NewCard(
			twittercard.CardSummary,
			"Example Title",
			"This is a description of the content.",
			"https://placehold.co/600x400?text=JD",
			"@example_site",
			"@example_creator",
		),
	}

	err = pages.AboutPage(headerItems).Render(r.Context(), w)
	if err != nil {
		return
	}
}
