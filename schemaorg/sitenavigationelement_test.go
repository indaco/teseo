package schemaorg

import (
	"bytes"
	"os"
	"reflect"
	"testing"
)

// Sample XML data for testing
const sampleSitemapXML = `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>http://www.example.com/</loc>
    <priority>0.5</priority>
  </url>
  <url>
    <loc>http://www.example.com/about</loc>
    <priority>0.5</priority>
  </url>
</urlset>`

// Sample Go struct data for testing
var sampleSiteNav = &SiteNavigationElement{
	Context: "https://schema.org",
	Type:    "SiteNavigationElement",
	ItemList: &ItemList{
		Context: "https://schema.org",
		Type:    "ItemList",
		ItemListElement: []ItemListElement{
			{
				Type:     "SiteNavigationElement",
				URL:      "http://www.example.com/",
				Position: 1,
			},
			{
				Type:     "SiteNavigationElement",
				URL:      "http://www.example.com/about",
				Position: 2,
			},
		},
	},
}

// TestToSitemapFile tests the ToSitemapFile function
func TestToSitemapFile(t *testing.T) {
	// Create a temporary file to write the sitemap
	tempFile, err := os.CreateTemp("", "sitemap-*.xml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up after test

	// Call ToSitemapFile to write to the temp file
	err = sampleSiteNav.ToSitemapFile(tempFile.Name())
	if err != nil {
		t.Fatalf("ToSitemapFile failed: %v", err)
	}

	// Read the file and check the output
	output, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read generated sitemap file: %v", err)
	}

	// Compare the generated XML with the expected output
	if !bytes.Equal(output, []byte(sampleSitemapXML)) {
		t.Errorf("Generated XML does not match expected XML.\nExpected:\n%s\nGot:\n%s", sampleSitemapXML, string(output))
	}
}

// TestFromSitemapFile tests the FromSitemapFile function
func TestFromSitemapFile(t *testing.T) {
	// Create a temporary file with sample XML content
	tempFile, err := os.CreateTemp("", "sitemap-*.xml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up after test

	// Write the sample XML to the temp file
	_, err = tempFile.WriteString(sampleSitemapXML)
	if err != nil {
		t.Fatalf("Failed to write sample sitemap XML: %v", err)
	}

	// Reset the file offset to the beginning for reading
	if _, err := tempFile.Seek(0, 0); err != nil {
		t.Fatalf("Failed to reset file offset: %v", err)
	}

	// Create an empty SiteNavigationElement to load data into
	var siteNav SiteNavigationElement

	// Call FromSitemapFile to populate the struct
	err = siteNav.FromSitemapFile(tempFile.Name())
	if err != nil {
		t.Fatalf("FromSitemapFile failed: %v", err)
	}

	// Compare the loaded struct with the expected data
	if !reflect.DeepEqual(&siteNav, sampleSiteNav) {
		t.Errorf("Loaded SiteNavigationElement does not match expected struct.\nExpected:\n%+v\nGot:\n%+v", sampleSiteNav, &siteNav)
	}
}
