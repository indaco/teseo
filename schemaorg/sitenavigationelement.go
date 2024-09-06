package schemaorg

import (
	"context"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"

	"github.com/a-h/templ"
	"github.com/indaco/teseo"
)

// SiteNavigationElement represents a Schema.org SiteNavigationElement object.
//
// Example usage:
//
// Pure struct usage:
//
//	siteNavElement := &schemaorg.SiteNavigationElement{
//		Name: "Main Navigation",
//		URL:  "https://www.example.com",
//		ItemList: &schemaorg.ItemList{
//			ItemListElement: []schemaorg.ItemListElement{
//				{Name: "Home", URL: "https://www.example.com", Position: 1},
//				{Name: "About", URL: "https://www.example.com/about", Position: 2},
//			},
//		},
//	}
//
// Factory method usage:
//
//	siteNavElement := schemaorg.NewSiteNavigationElementWithItemList(
//		"Main Navigation",
//		"https://www.example.com",
//		[]schemaorg.ItemListElement{
//			{Name: "Home", URL: "https://www.example.com", Position: 1},
//			{Name: "About", URL: "https://www.example.com/about", Position: 2},
//		},
//	)
//
// // Rendering JSON-LD using templ:
//
//	templ Page() {
//		@siteNavElement.ToJsonLd()
//	}
//
// // Rendering JSON-LD as `template.HTML` value:
//
//	jsonLdHtml := siteNavElement.ToGoHTMLJsonLd()
//
// Expected output:
//
//	{
//		"@context": "https://schema.org",
//		"@type": "SiteNavigationElement",
//		"name": "Main Navigation",
//		"url": "https://www.example.com",
//		"itemListElement": [
//			{"@type": "ListItem", "position": 1, "name": "Home", "url": "https://www.example.com"},
//			{"@type": "ListItem", "position": 2, "name": "About", "url": "https://www.example.com/about"}
//		]
//	}
//
// Example usage with `ToSitemapFile`:
//
//	// Generate a sitemap XML file
//	siteNavElement := &schemaorg.SiteNavigationElement{
//		Name: "Main Navigation",
//		URL:  "https://www.example.com",
//		ItemList: &schemaorg.ItemList{
//			ItemListElement: []schemaorg.ItemListElement{
//				{Name: "Home", URL: "https://www.example.com", Position: 1},
//				{Name: "About", URL: "https://www.example.com/about", Position: 2},
//			},
//		},
//	}
//	err := siteNavElement.ToSitemapFile("statics/sitemap.xml")
//	if err != nil {
//		log.Fatalf("Failed to generate sitemap: %v", err)
//	}
//
// Example usage with `FromSitemapFile`:
//
//	// Parse a sitemap XML file and populate the SiteNavigationElement struct
//	siteNavElement := &schemaorg.SiteNavigationElement{}
//	err := siteNavElement.FromSitemapFile("statics/sitemap.xml")
//	if err != nil {
//		log.Fatalf("Failed to parse sitemap: %v", err)
//	}
//
// Expected output:
//
//	<?xml version="1.0" encoding="UTF-8"?>
//	<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
//	  <url>
//	    <loc>http://www.example.com/</loc>
//	    <priority>0.5</priority>
//	  </url>
//	  <url>
//	    <loc>http://www.example.com/about</loc>
//	    <priority>0.5</priority>
//	  </url>
//	</urlset>
type SiteNavigationElement struct {
	Context    string    `json:"@context"`
	Type       string    `json:"@type"`
	Name       string    `json:"name,omitempty"`
	URL        string    `json:"url,omitempty"`
	Position   int       `json:"position,omitempty"`
	Identifier string    `json:"identifier,omitempty"`
	ItemList   *ItemList `json:"itemList,omitempty"`
}

// ItemList represents a Schema.org ItemList object
type ItemList struct {
	Context         string            `json:"@context"`
	Type            string            `json:"@type"`
	ItemListElement []ItemListElement `json:"itemListElement"`
}

// ItemListElement represents an individual item in an ItemList
type ItemListElement struct {
	Type     string `json:"@type"`
	Name     string `json:"name,omitempty"`
	URL      string `json:"url,omitempty"`
	Position int    `json:"position,omitempty"`
}

// XMLSitemapUrl represents a single URL entry in the sitemap XML.
type XMLSitemapUrl struct {
	Loc      string `xml:"loc"`
	Priority string `xml:"priority,omitempty"`
}

// XMLSitemap represents the structure of a sitemap XML file.
type XMLSitemap struct {
	XMLName xml.Name        `xml:"urlset"`
	Xmlns   string          `xml:"xmlns,attr"`
	Urls    []XMLSitemapUrl `xml:"url"`
}

// NewSiteNavigationElement initializes a SiteNavigationElement with default context and type.
func NewSiteNavigationElement(name string, url string, position int, identifier string, itemList *ItemList) *SiteNavigationElement {
	sne := &SiteNavigationElement{
		Name:       name,
		URL:        url,
		Position:   position,
		Identifier: identifier,
		ItemList:   itemList,
	}
	sne.ensureDefaults()
	return sne
}

// NewSiteNavigationElementWithItemList initializes a SiteNavigationElement with default context, type, and an ItemList.
func NewSiteNavigationElementWithItemList(name, url string, items []ItemListElement) *SiteNavigationElement {
	// Create a new ItemList
	itemList := &ItemList{
		Context:         "https://schema.org",
		Type:            "ItemList",
		ItemListElement: items,
	}

	// Create a new SiteNavigationElement with the provided ItemList
	siteNavElement := &SiteNavigationElement{
		Name:     name,
		URL:      url,
		ItemList: itemList,
	}

	// Ensure defaults are set
	siteNavElement.ensureDefaults()

	return siteNavElement
}

// NewItemListElement creates a new ItemListElement with default values.
func NewItemListElement(name, url string, position int) ItemListElement {
	return ItemListElement{
		Type:     "ListItem",
		Name:     name,
		URL:      url,
		Position: position,
	}
}

// NewItemList creates a new ItemList with default values.
func NewItemList(elements []ItemListElement) ItemList {
	return ItemList{
		Context:         "https://schema.org",
		Type:            "ItemList",
		ItemListElement: elements,
	}
}

// ToJsonLd converts the SiteNavigationElement struct to a JSON-LD `templ.Component`.
func (sne *SiteNavigationElement) ToJsonLd() templ.Component {
	sne.ensureDefaults()
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		return templ.JSONScript(teseo.GenerateUniqueKey(), sne).WithType("application/ld+json").Render(ctx, w)
	})
}

// ToGoHTMLJsonLd renders the SiteNavigationElement struct as `template.HTML` value for Go's `html/template`.
func (sne *SiteNavigationElement) ToGoHTMLJsonLd() (template.HTML, error) {
	sne.ensureDefaults()

	// Create the templ component.
	templComponent := sne.ToJsonLd()

	// Render the templ component to a `template.HTML` value.
	html, err := templ.ToGoHTML(context.Background(), templComponent)
	if err != nil {
		log.Fatalf("failed to convert to html: %v", err)
	}

	return html, nil
}

// ToSitemapFile generates a sitemap XML file from the SiteNavigationElement struct.
func (s *SiteNavigationElement) ToSitemapFile(filename string) error {
	if s.ItemList == nil {
		return fmt.Errorf("ItemList is nil, cannot generate sitemap")
	}

	// Populate the XML structure with the necessary namespace
	var sitemap XMLSitemap
	sitemap.Xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9" // Set the namespace

	for _, item := range s.ItemList.ItemListElement {
		// Add each item as an XML sitemap URL entry
		url := XMLSitemapUrl{
			Loc:      item.URL,
			Priority: "0.5", // Example priority, can be adjusted or made dynamic
		}
		sitemap.Urls = append(sitemap.Urls, url)
	}

	// Marshal the sitemap struct to XML
	xmlData, err := xml.MarshalIndent(sitemap, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling sitemap to XML: %v", err)
	}

	// Add the XML header
	xmlData = append([]byte(xml.Header), xmlData...)

	// Write the XML data to a file
	err = os.WriteFile(filename, xmlData, 0644)
	if err != nil {
		return fmt.Errorf("error writing XML file: %v", err)
	}

	return nil
}

// FromSitemapFile parses a sitemap XML file and populates the SiteNavigationElement struct.
func (s *SiteNavigationElement) FromSitemapFile(filename string) error {
	// Open the XML file
	xmlFile, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open XML file: %v", err)
	}
	defer xmlFile.Close()

	// Read the file content
	byteValue, err := io.ReadAll(xmlFile)
	if err != nil {
		return fmt.Errorf("could not read XML file: %v", err)
	}

	// Parse the XML content
	var sitemap XMLSitemap
	err = xml.Unmarshal(byteValue, &sitemap)
	if err != nil {
		return fmt.Errorf("could not unmarshal XML content: %v", err)
	}

	// Populate the SiteNavigationElement struct from the parsed XML
	s.Context = "https://schema.org"
	s.Type = "SiteNavigationElement"
	s.ItemList = &ItemList{
		Context: "https://schema.org",
		Type:    "ItemList",
	}

	for i, url := range sitemap.Urls {
		// Add each URL as an ItemListElement in the ItemList
		item := ItemListElement{
			Type:     "SiteNavigationElement",
			URL:      url.Loc,
			Position: i + 1, // Assign position incrementally
		}
		s.ItemList.ItemListElement = append(s.ItemList.ItemListElement, item)
	}

	return nil
}

// makeSiteNavigationElement initializes a SiteNavigationElement with default context and type.
func (sne *SiteNavigationElement) ensureDefaults() {
	if sne.Context == "" {
		sne.Context = "https://schema.org"
	}

	if sne.Type == "" {
		sne.Type = "SiteNavigationElement"
	}

	if sne.Position == 0 {
		sne.Position = 1
	}

	if sne.ItemList != nil && len(sne.ItemList.ItemListElement) > 0 {
		sne.ItemList = &ItemList{
			Context:         "https://schema.org",
			Type:            "ItemList",
			ItemListElement: sne.ItemList.ItemListElement,
		}
	}
}
