package schemaorg

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"reflect"
	"slices"
	"strings"
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
var sampleSiteNav = &SiteNavigationElementList{
	Context: "https://schema.org",
	Type:    "ItemList",
	ItemListElement: []SiteNavigationElement{
		{Type: "SiteNavigationElement", Name: "http://www.example.com/", Description: "", URL: "http://www.example.com/", Position: 1},
		{Type: "SiteNavigationElement", Name: "http://www.example.com/about", Description: "", URL: "http://www.example.com/about", Position: 2},
	},
}

func TestSiteNavigationElementList_EnsureDefaults(t *testing.T) {
	sne := &SiteNavigationElementList{
		ItemListElement: []SiteNavigationElement{},
	}
	sne.ensureDefaults()

	if sne.Context != "https://schema.org" {
		t.Errorf("expected context to be schema.org, got %s", sne.Context)
	}
	if sne.Type != "ItemList" {
		t.Errorf("expected type to be ItemList, got %s", sne.Type)
	}

	if len(sne.ItemListElement) != 0 {
		t.Errorf("expected no item, got %d", len(sne.ItemListElement))
	}
}

func TestNewSiteNavigationElementList(t *testing.T) {
	items := []SiteNavigationElement{{Position: 1, Name: "Home", URL: "https://example.com"}}
	snel := NewSiteNavigationElementList("", items)

	if snel.Context != "https://schema.org" {
		t.Errorf("expected context to be schema.org")
	}
	if snel.Type != "ItemList" {
		t.Errorf("expected type to be ItemList")
	}

}

func TestSiteNavigationElement_Validate(t *testing.T) {
	tests := []struct {
		name     string
		input    *SiteNavigationElementList
		expected []string
	}{
		{
			name: "valid input",
			input: NewSiteNavigationElementList("Nav", []SiteNavigationElement{
				NewSimpleSiteNavigationElement(1, "Home", "https://example.com"),
			}),
			expected: nil,
		},
		{
			name:  "empty item list",
			input: &SiteNavigationElementList{},
			expected: []string{
				"ItemList should contain at least one item",
			},
		},
		{
			name: "missing fields in item list element",
			input: &SiteNavigationElementList{
				Identifier: "Nav",
				ItemListElement: []SiteNavigationElement{
					{}, // all fields missing
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
			_, warnings := tt.input.Validate()
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

func TestNewSiteNavigationElement_SetsDefaults(t *testing.T) {
	element := &SiteNavigationElement{Position: 1, Name: "Home", URL: "https://example.com"}
	element.ensureDefaults()

	if element.Type != "SiteNavigationElement" {
		t.Errorf("expected type to be SiteNavigationElement, got %s", element.Type)
	}
}

func TestNewSiteNavigationElementsFromLinks(t *testing.T) {
	links := []NavigationLink{
		{Name: "Home", URL: "https://example.com", Description: "Homepage"},
		{Name: "About", URL: "https://example.com/about", Description: "About us"},
	}

	expected := []SiteNavigationElement{
		{
			Type:        "SiteNavigationElement",
			Position:    1,
			Name:        "Home",
			Description: "Homepage",
			URL:         "https://example.com",
		},
		{
			Type:        "SiteNavigationElement",
			Position:    2,
			Name:        "About",
			Description: "About us",
			URL:         "https://example.com/about",
		},
	}

	got := NewSiteNavigationElementsFromLinks(links)

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("unexpected SiteNavigationElement slice\nExpected: %+v\nGot: %+v", expected, got)
	}
}

func TestSiteNavigationElementList_ToJsonLd_WritesWarnings(t *testing.T) {
	// Save original stderr and restore after test
	origStderr := os.Stderr
	defer func() { os.Stderr = origStderr }()

	// Create pipe to capture stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stderr = w

	// Create input with missing fields to trigger warnings
	list := &SiteNavigationElementList{
		Identifier: "test",
		ItemListElement: []SiteNavigationElement{
			{}, // All fields missing
		},
	}

	// Run ToJsonLd to capture stderr output
	component := list.ToJsonLd()

	// Close writer and read stderr content
	if err := w.Close(); err != nil {
		t.Fatalf("failed to close pipe writer: %v", err)
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("failed to read stderr: %v", err)
	}

	output := buf.String()

	// Check expected messages
	expectedWarnings := []string{
		"missing name in ItemListElement at position 1",
		"missing url in ItemListElement at position 1",
		"missing position in ItemListElement at index 0",
	}

	for _, warning := range expectedWarnings {
		if !strings.Contains(output, warning) {
			t.Errorf("expected warning %q in stderr output, got:\n%s", warning, output)
		}
	}

	// Basic check that a templ.Component was returned
	if component == nil {
		t.Errorf("expected non-nil component from ToJsonLd")
	}
}

func TestSiteNavigationElement_ToGoHTMLJsonLd(t *testing.T) {
	sne := NewSiteNavigationElementList("Main Nav", []SiteNavigationElement{
		NewSimpleSiteNavigationElement(1, "Home", "https://example.com"),
	})

	html, err := sne.ToGoHTMLJsonLd()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if html == "" {
		t.Errorf("expected non-empty HTML output")
	}
}

// TestFromSitemapFile tests the FromSitemapFile function
func TestFromSitemapFile(t *testing.T) {
	// Create a temporary file with sample XML content
	tempFile, err := os.CreateTemp("", "sitemap-*.xml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() {
		if err := os.Remove(tempFile.Name()); err != nil {
			t.Logf("warning: failed to remove temp file: %v", err)
		}
	}()

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
	var siteNav SiteNavigationElementList

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

func TestFromSitemapFile_Errors(t *testing.T) {
	sne := &SiteNavigationElementList{}

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
	defer func() {
		if err := os.Remove(tempFile.Name()); err != nil {
			t.Logf("warning: failed to remove temp file: %v", err)
		}
	}()

	_, err = tempFile.WriteString("<<< invalid xml >>>")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	// Close writer
	if err := tempFile.Close(); err != nil {
		t.Fatalf("failed to close pipe writer: %v", err)
	}

	err = sne.FromSitemapFile(tempFile.Name())
	if err == nil {
		t.Errorf("expected error for invalid XML")
	}
}

type brokenCloser struct {
	io.Reader
}

func (b *brokenCloser) Close() error {
	return fmt.Errorf("simulated close error")
}

func TestFromSitemapFile_CloseError(t *testing.T) {
	// Backup and restore original openFile
	origOpenFile := openFile
	defer func() { openFile = origOpenFile }()

	// Minimal valid XML
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
	<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
		<url><loc>https://example.com</loc></url>
	</urlset>`

	// override openFile inside the test function
	openFile = func(name string) (io.ReadCloser, error) {
		return &brokenCloser{Reader: strings.NewReader(xmlData)}, nil
	}

	var sne SiteNavigationElementList
	err := sne.FromSitemapFile("dummy.xml")

	if err == nil || !strings.Contains(err.Error(), "failed to close file") {
		t.Errorf("expected close error, got: %v", err)
	}
}

// TestToSitemapFile tests the ToSitemapFile function
func TestToSitemapFile(t *testing.T) {
	// Create a temporary file to write the sitemap
	tempFile, err := os.CreateTemp("", "sitemap-*.xml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() {
		if err := os.Remove(tempFile.Name()); err != nil {
			t.Logf("warning: failed to remove temp file: %v", err)
		}
	}()

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

	// Compare the generated XML with the expected output by unmarshalling both
	var expected, actual XMLSitemap
	if err := xml.Unmarshal([]byte(sampleSitemapXML), &expected); err != nil {
		t.Fatalf("invalid expected XML: %v", err)
	}
	if err := xml.Unmarshal(output, &actual); err != nil {
		t.Fatalf("invalid generated XML: %v", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("sitemap XML mismatch\nExpected: %+v\nGot: %+v", expected, actual)
	}
}

func TestToSitemapFile_Errors(t *testing.T) {
	sne := &SiteNavigationElementList{}

	// ItemList is nil
	err := sne.ToSitemapFile("dummy.xml")
	expected := "failed to generate sitemap XML: item list is nil, cannot generate sitemap"
	if err == nil || err.Error() != expected {
		t.Errorf("expected error %q, got %v", expected, err)
	}

	// XML marshal error (simulate by injecting invalid data if needed)

	// File write error (read-only path)
	err = sne.ToSitemapFile("/tmp/tmpoai36nm9") // from previous execution
	if err == nil {
		t.Errorf("expected write error, got nil")
	}
}

func TestToSitemapFile_WriteFileError(t *testing.T) {
	original := writeFile
	defer func() { writeFile = original }()

	writeFile = func(name string, data []byte, perm os.FileMode) error {
		return fmt.Errorf("mock write error")
	}

	sne := &SiteNavigationElementList{
		ItemListElement: []SiteNavigationElement{
			{URL: "https://example.com"},
		},
	}
	err := sne.ToSitemapFile("dummy.xml")
	expected := `failed to write sitemap file "dummy.xml": mock write error`
	if err == nil || err.Error() != expected {
		t.Errorf("expected %q, got %v", expected, err)
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

	sne := &SiteNavigationElementList{
		ItemListElement: []SiteNavigationElement{
			{URL: "https://example.com"},
		},
	}

	err := sne.ToSitemapFile("dummy.xml")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	expected := "failed to generate sitemap XML: error marshaling sitemap XML: simulated marshal error"
	if err.Error() != expected {
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

	sne := &SiteNavigationElementList{}
	err := sne.FromSitemapFile("dummy.xml")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	expected := "could not read XML file: simulated read error"
	if err.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, err.Error())
	}
}
