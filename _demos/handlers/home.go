package handlers

import (
	"log"
	"net/http"

	"github.com/indaco/teseo/_demos/pages"
	"github.com/indaco/teseo/_demos/types"
	"github.com/indaco/teseo/schemaorg"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	headerItems := &types.SEOItems{
		WebPage: &schemaorg.WebPage{
			URL:           "https://www.example.com",
			Name:          "Home Page",
			Headline:      "Learn More About Us",
			Description:   "This is the home page of Example.com, providing information about out services offering.",
			About:         "Company Information",
			Keywords:      "about us, company, mission, values",
			InLanguage:    "en",
			IsPartOf:      "https://www.example.com",
			LastReviewed:  "2024-09-01",
			PrimaryImage:  "https://www.example.com/images/about-us.jpg",
			DatePublished: "2020-01-01",
			DateModified:  "2024-09-01",
		},
		SiteNavElement: &schemaorg.SiteNavigationElement{
			Name: "Main Navigation",
			URL:  "https://www.example.com",
			ItemList: &schemaorg.ItemList{
				ItemListElement: []schemaorg.ItemListElement{
					{Name: "Home", URL: "https://www.example.com", Position: 1},
					{Name: "About", URL: "https://www.example.com/about", Position: 2},
				},
			},
		},
	}

	err := headerItems.SiteNavElement.ToSitemapFile("./_demos/statics/sitemap.xml")
	if err != nil {
		log.Fatalf("Failed to generate sitemap: %v", err)
	}

	err = pages.HomePage(headerItems).Render(r.Context(), w)
	if err != nil {
		return
	}
}
