package schemaorg

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"slices"
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

func TestSiteNavigationElement_EnsureDefaults(t *testing.T) {
	sne := &SiteNavigationElement{
		ItemList: &ItemList{},
	}
	sne.ensureDefaults()

	if sne.Context != "https://schema.org" {
		t.Errorf("expected context to be schema.org, got %s", sne.Context)
	}
	if sne.Type != "SiteNavigationElement" {
		t.Errorf("expected type to be SiteNavigationElement, got %s", sne.Type)
	}
	if sne.Position != 1 {
		t.Errorf("expected default position to be 1, got %d", sne.Position)
	}
	if sne.ItemList.Context != "https://schema.org" {
		t.Errorf("expected itemList context to be schema.org, got %s", sne.ItemList.Context)
	}
	if sne.ItemList.Type != "ItemList" {
		t.Errorf("expected itemList type to be ItemList, got %s", sne.ItemList.Type)
	}
}

func TestNewSiteNavigationElement_SetsDefaults(t *testing.T) {
	itemList := &ItemList{ItemListElement: []ItemListElement{{Name: "Home", URL: "https://example.com", Position: 1}}}
	sne := NewSiteNavigationElement("Main Nav", "https://example.com", 2, "main", itemList)

	if sne.Context != "https://schema.org" {
		t.Errorf("expected context to be schema.org, got %s", sne.Context)
	}
	if sne.Type != "SiteNavigationElement" {
		t.Errorf("expected type to be SiteNavigationElement, got %s", sne.Type)
	}
	if sne.Position != 2 {
		t.Errorf("expected position 2, got %d", sne.Position)
	}
	if sne.ItemList == nil || sne.ItemList.Type != "ItemList" {
		t.Errorf("expected ItemList with type ItemList")
	}
}

func TestNewItemList_SetsDefaults(t *testing.T) {
	items := []ItemListElement{{Name: "Home", URL: "https://example.com", Position: 1}}
	list := NewItemList(items)

	if list.Context != "https://schema.org" {
		t.Errorf("expected context to be schema.org")
	}
	if list.Type != "ItemList" {
		t.Errorf("expected type to be ItemList")
	}
	if len(list.ItemListElement) != 1 {
		t.Errorf("expected one item, got %d", len(list.ItemListElement))
	}
}

func TestSiteNavigationElement_ToGoHTMLJsonLd(t *testing.T) {
	sne := NewSiteNavigationElementWithItemList("Main Nav", "https://example.com", []ItemListElement{
		NewItemListElement("Home", "https://example.com", 1),
	})

	html, err := sne.ToGoHTMLJsonLd()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if html == "" {
		t.Errorf("expected non-empty HTML output")
	}
}

func TestSiteNavigationElement_Validate(t *testing.T) {
	tests := []struct {
		name     string
		input    *SiteNavigationElement
		expected []string
	}{
		{
			name: "valid input",
			input: NewSiteNavigationElementWithItemList("Nav", "https://example.com", []ItemListElement{
				NewItemListElement("Home", "https://example.com", 1),
			}),
			expected: nil,
		},
		{
			name:  "missing name and url",
			input: &SiteNavigationElement{},
			expected: []string{
				"missing recommended field: name",
				"missing recommended field: url",
				"ItemList should contain at least one item",
			},
		},
		{
			name: "missing fields in item list element",
			input: &SiteNavigationElement{
				Name: "Nav", URL: "https://example.com",
				ItemList: &ItemList{
					ItemListElement: []ItemListElement{
						{}, // all fields missing
					},
				},
			},
			expected: []string{
				"missing name in ItemListElement at position 1",
				"missing url in ItemListElement at position 1",
				"missing position in ItemListElement at index 0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			warnings := tt.input.Validate()
			if len(warnings) != len(tt.expected) {
				t.Errorf("expected %d warnings, got %d: %v", len(tt.expected), len(warnings), warnings)
				return
			}
			for _, expected := range tt.expected {
				found := slices.Contains(warnings, expected)
				if !found {
					t.Errorf("expected warning %q not found in %v", expected, warnings)
				}
			}
		})
	}
}

func TestToSitemapFile_Errors(t *testing.T) {
	sne := &SiteNavigationElement{}

	// ItemList is nil
	err := sne.ToSitemapFile("dummy.xml")
	if err == nil || err.Error() != "ItemList is nil, cannot generate sitemap" {
		t.Errorf("expected error for nil ItemList, got %v", err)
	}

	// XML marshal error (simulate by injecting invalid data if needed)

	// File write error (read-only path)
	err = sne.ToSitemapFile("/tmp/tmpoai36nm9") // from previous execution
	if err == nil {
		t.Errorf("expected write error, got nil")
	}
}

func TestFromSitemapFile_Errors(t *testing.T) {
	sne := &SiteNavigationElement{}

	// File not found
	err := sne.FromSitemapFile("/nonexistent/path/to/file.xml") // from previous execution
	if err == nil {
		t.Errorf("expected error for nonexistent file")
	}

	// Invalid XML
	tempFile, err := os.CreateTemp("", "invalid-*.xml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	_, err = tempFile.WriteString("<<< invalid xml >>>")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tempFile.Close()

	err = sne.FromSitemapFile(tempFile.Name())
	if err == nil {
		t.Errorf("expected error for invalid XML")
	}
}

func TestToSitemapFile_WriteFileError(t *testing.T) {
	original := writeFile
	defer func() { writeFile = original }()

	writeFile = func(name string, data []byte, perm os.FileMode) error {
		return fmt.Errorf("mock write error")
	}

	sne := &SiteNavigationElement{
		ItemList: &ItemList{
			ItemListElement: []ItemListElement{
				{URL: "https://example.com"},
			},
		},
	}
	err := sne.ToSitemapFile("dummy.xml")
	if err == nil || err.Error() != "error writing XML file: mock write error" {
		t.Errorf("expected mock write error, got %v", err)
	}
}

func TestToSitemapFile_MarshalIndentError(t *testing.T) {
	// Backup and restore original marshalIndent
	originalMarshal := marshalIndent
	defer func() { marshalIndent = originalMarshal }()

	// Simulate marshal error
	marshalIndent = func(v any, prefix, indent string) ([]byte, error) {
		return nil, fmt.Errorf("simulated marshal error")
	}

	sne := &SiteNavigationElement{
		ItemList: &ItemList{
			ItemListElement: []ItemListElement{
				{URL: "https://example.com"},
			},
		},
	}

	err := sne.ToSitemapFile("dummy.xml")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if err.Error() != "error marshaling sitemap to XML: simulated marshal error" {
		t.Errorf("unexpected error: %v", err)
	}
}

type faultyReader struct{}

func (faultyReader) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("simulated read error")
}
func (faultyReader) Close() error { return nil }

func TestFromSitemapFile_ReadError(t *testing.T) {
	originalOpen := openFile
	defer func() { openFile = originalOpen }()

	openFile = func(name string) (io.ReadCloser, error) {
		return faultyReader{}, nil
	}

	sne := &SiteNavigationElement{}
	err := sne.FromSitemapFile("dummy.xml")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expected := "could not read XML file: simulated read error"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}
