package schemaorg

// Common type definitions used across multiple JSON-LD entities

// ContactPoint represents a Schema.org ContactPoint object
// For more details about the meaning of the properties see: https://schema.org/ContactPoint
type ContactPoint struct {
	Type              string `json:"@type"`
	Telephone         string `json:"telephone,omitempty"`
	ContactType       string `json:"contactType,omitempty"`
	AreaServed        string `json:"areaServed,omitempty"`
	AvailableLanguage string `json:"availableLanguage,omitempty"`
}

// ImageObject represents a Schema.org ImageObject object
// For more details about the meaning of the properties see: https://schema.org/ImageObject
type ImageObject struct {
	Type string `json:"@type"`
	URL  string `json:"url,omitempty"`
}

// ensureDefaults sets default values for ImageObject if they are not already set.
func (img *ImageObject) ensureDefaults() {
	if img.Type == "" {
		img.Type = "ImageObject"
	}
}

// ListItem represents a Schema.org ListItem object
// For more details about the meaning of the properties see: https://schema.org/ListItem
type ListItem struct {
	Type     string `json:"@type"`
	Position int    `json:"position,omitempty"`
	Name     string `json:"name,omitempty"`
	Item     string `json:"item,omitempty"`
}
