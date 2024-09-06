package handlers

import (
	"net/http"

	"github.com/indaco/teseo/_demos/pages"
	"github.com/indaco/teseo/_demos/types"
	"github.com/indaco/teseo/schemaorg"
)

func HandleCompany(w http.ResponseWriter, r *http.Request) {
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
	}

	err := pages.CompanyPage(headerItems).Render(r.Context(), w)
	if err != nil {
		return
	}
}
