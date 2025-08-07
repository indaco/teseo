package schemaorg

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"os"

	"github.com/a-h/templ"
	"github.com/indaco/teseo"
)

var (
	writeFile                                              = os.WriteFile
	marshalIndent                                          = xml.MarshalIndent
	openFile      func(name string) (io.ReadCloser, error) = func(name string) (io.ReadCloser, error) {
		return os.Open(name)
	}
)

// SiteNavigationElementList represents a Schema.org `ItemList` object whose
// `itemListElement` consists of `SiteNavigationElement` entries.
//
// It is typically used to describe a website's main navigation structure in
// structured data (JSON-LD), and can also be converted into a sitemap XML file.
//
// See: https://schema.org/SiteNavigationElement
//
// Example (using helper functions):
//
//	snl := schemaorg.NewSiteNavigationElementList("main", []schemaorg.SiteNavigationElement{
//		schemaorg.NewSiteNavigationElement(1, "Home", "Go to homepage", "https://www.example.com"),
//		schemaorg.NewSiteNavigationElement(2, "About", "About us", "https://www.example.com/about"),
//	})
//
//	// Render as JSON-LD with templ
//	templ Page() {
//		@snl.ToJsonLd()
//	}
//
//	// Render as Go template.HTML
//	html, err := snl.ToGoHTMLJsonLd()
//
//	// Generate a sitemap XML file
//	err := snl.ToSitemapFile("public/sitemap.xml")
//
//	// Load from a sitemap XML file
//	var snl schemaorg.SiteNavigationElementList
//	err := snl.FromSitemapFile("public/sitemap.xml")
//
// JSON-LD Output:
//
//	{
//	  "@context": "https://schema.org",
//	  "@type": "ItemList",
//	  "identifier": "main",
//	  "itemListElement": [
//	    {
//	      "@type": "SiteNavigationElement",
//	      "position": 1,
//	      "name": "Home",
//	      "description": "Go to homepage",
//	      "url": "https://www.example.com"
//	    },
//	    {
//	      "@type": "SiteNavigationElement",
//	      "position": 2,
//	      "name": "About",
//	      "description": "About us",
//	      "url": "https://www.example.com/about"
//	    }
//	  ]
//	}
//
// Sitemap XML Output:
//
//	<?xml version="1.0" encoding="UTF-8"?>
//	<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
//	  <url>
//	    <loc>https://www.example.com</loc>
//	    <priority>0.5</priority>
//	  </url>
//	  <url>
//	    <loc>https://www.example.com/about</loc>
//	    <priority>0.5</priority>
//	  </url>
//	</urlset>

// --------------------------
// Schema.org Types
// --------------------------

// SiteNavigationElement represents a Schema.org SiteNavigationElement.
type SiteNavigationElement struct {
	Type        string `json:"@type"`
	Position    int    `json:"position,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}

// SiteNavigationElementList represents an ItemList of SiteNavigationElement.
type SiteNavigationElementList struct {
	Context         string                  `json:"@context"`
	Type            string                  `json:"@type"`
	Identifier      string                  `json:"identifier,omitempty"`
	ItemListElement []SiteNavigationElement `json:"itemListElement,omitempty"`
}

// --------------------------
// XML Sitemap Types
// --------------------------

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

// NavigationLink represents a structured navigation item with optional description.
type NavigationLink struct {
	Name        string
	URL         string
	Description string
}

// --------------------------
// Constructors & Defaults
// -------------------------

// NewSiteNavigationElementList creates a SiteNavigationElementList (ItemList)
// with default context/type and an optional identifier.
func NewSiteNavigationElementList(
	identifier string,
	items []SiteNavigationElement,
) *SiteNavigationElementList {
	return &SiteNavigationElementList{
		Context:         "https://schema.org",
		Type:            "ItemList",
		Identifier:      identifier,
		ItemListElement: items,
	}
}

// NewSiteNavigationElement creates a SiteNavigationElement with full control over all fields.
func NewSiteNavigationElement(position int, name, description, url string) SiteNavigationElement {
	return SiteNavigationElement{
		Type:        "SiteNavigationElement",
		Position:    position,
		Name:        name,
		Description: description,
		URL:         url,
	}
}

// NewSimpleSiteNavigationElement creates a SiteNavigationElement with no description.
func NewSimpleSiteNavigationElement(position int, name, url string) SiteNavigationElement {
	return NewSiteNavigationElement(position, name, "", url)
}

// NewSiteNavigationElementsFromLinks creates SiteNavigationElement items from a slice
// of NavigationLink structs. The position is assigned incrementally starting at 1.
func NewSiteNavigationElementsFromLinks(links []NavigationLink) []SiteNavigationElement {
	var result []SiteNavigationElement
	for i, link := range links {
		result = append(result, NewSiteNavigationElement(i+1, link.Name, link.Description, link.URL))
	}
	return result
}

// ensureDefaults initializes a SiteNavigationElementList with default context and type.
func (snl *SiteNavigationElementList) ensureDefaults() {
	if snl.Context == "" {
		snl.Context = "https://schema.org"
	}

	if snl.Type == "" {
		snl.Type = "ItemList"
	}

}

func (sne *SiteNavigationElement) ensureDefaults() {
	if sne.Type == "" {
		sne.Type = "SiteNavigationElement"
	}
}

// --------------------------
// Validation
// --------------------------

func (snl *SiteNavigationElementList) Validate() (bool, []string) {
	var warnings []string

	if len(snl.ItemListElement) == 0 {
		warnings = append(warnings, "ItemList should contain at least one item")
	} else {
		for i, item := range snl.ItemListElement {
			if item.Name == "" {
				warnings = append(warnings, fmt.Sprintf("missing name in ItemListElement at position %d", i+1))
			}
			if item.URL == "" {
				warnings = append(warnings, fmt.Sprintf("missing url in ItemListElement at position %d", i+1))
			}
			if item.Position == 0 {
				warnings = append(warnings, fmt.Sprintf("missing position in ItemListElement at index %d", i))
			}
		}
	}

	return len(warnings) == 0, warnings
}

// --------------------------
// JSON-LD Rendering
// --------------------------

// ToJsonLd converts the SiteNavigationElement struct to a JSON-LD `templ.Component`.
func (snl *SiteNavigationElementList) ToJsonLd() templ.Component {
	snl.ensureDefaults()

	ok, msgs := snl.Validate()
	if !ok {
		for _, msg := range msgs {
			fmt.Fprintln(os.Stderr, "⚠️", msg)
		}
	}

	for i := range snl.ItemListElement {
		snl.ItemListElement[i].ensureDefaults()
	}

	id := "siteNavItemList-" + teseo.GenerateUniqueKey()
	if snl.Identifier != "" {
		id = "siteNavItemList-" + snl.Identifier
	}

	return templ.JSONScript(id, snl).WithType("application/ld+json")
}

// ToGoHTMLJsonLd renders the SiteNavigationElement struct as `template.HTML` value for Go's `html/template`.
func (itemList *SiteNavigationElementList) ToGoHTMLJsonLd() (template.HTML, error) {
	return teseo.RenderToHTML(itemList.ToJsonLd())
}

// --------------------------
// Sitemap File Handling
// --------------------------

// ToSitemapBytes returns the XML sitemap content as a byte slice.
func (itemList *SiteNavigationElementList) ToSitemapBytes() ([]byte, error) {
	if itemList.ItemListElement == nil {
		return nil, fmt.Errorf("item list is nil, cannot generate sitemap")
	}

	sitemap := XMLSitemap{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
	}

	for _, item := range itemList.ItemListElement {
		sitemap.Urls = append(sitemap.Urls, XMLSitemapUrl{
			Loc:      item.URL,
			Priority: "0.5",
		})
	}

	data, err := marshalIndent(sitemap, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshaling sitemap XML: %w", err)
	}

	return append([]byte(xml.Header), data...), nil
}

// ToSitemapFile generates a sitemap XML file and writes it to the specified path.
func (itemList *SiteNavigationElementList) ToSitemapFile(filename string) error {
	data, err := itemList.ToSitemapBytes()
	if err != nil {
		return fmt.Errorf("failed to generate sitemap XML: %w", err)
	}

	if err := writeFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write sitemap file %q: %w", filename, err)
	}

	return nil
}

// FromSitemapFile parses a sitemap XML file and populates the SiteNavigationElement struct.
func (itemList *SiteNavigationElementList) FromSitemapFile(filename string) (err error) {
	// Open the XML file
	xmlFile, err := openFile(filename)
	if err != nil {
		return fmt.Errorf("failed to open sitemap XML file %q: %w", filename, err)
	}
	defer func() {
		if cerr := xmlFile.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close file: %w", cerr)
		}
	}()

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
	itemList.ensureDefaults()

	for i, url := range sitemap.Urls {
		item := SiteNavigationElement{
			Type:        "SiteNavigationElement",
			URL:         url.Loc,
			Name:        url.Loc,
			Description: "",
			Position:    i + 1,
		}
		itemList.ItemListElement = append(itemList.ItemListElement, item)
	}

	return nil
}
