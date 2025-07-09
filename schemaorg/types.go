package schemaorg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// Common type definitions used across multiple JSON-LD entities

// StringList supports both a single string or a slice of strings in JSON,
// and always trims input strings.
type StringList []string

// IsZero reports whether the value is empty or contains only one empty string.
func (s StringList) IsZero() bool {
	return len(s) == 0 || (len(s) == 1 && s[0] == "")
}

// MarshalJSON encodes the value as either a string or an array of strings.
func (s StringList) MarshalJSON() ([]byte, error) {
	switch len(s) {
	case 0:
		return []byte("null"), nil
	case 1:
		if strings.TrimSpace(s[0]) == "" {
			return []byte("null"), nil
		}
		return json.Marshal(s[0])
	default:
		return json.Marshal([]string(s))
	}
}

// UnmarshalJSON accepts both string and []string, and trims all values.
func (s *StringList) UnmarshalJSON(data []byte) error {
	// Handle null
	if bytes.Equal(data, []byte("null")) {
		*s = nil
		return nil
	}

	// Try single string
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*s = StringList{strings.TrimSpace(single)}
		return nil
	}

	// Try slice of strings
	var multi []string
	if err := json.Unmarshal(data, &multi); err == nil {
		trimmed := make([]string, 0, len(multi))
		for _, v := range multi {
			v = strings.TrimSpace(v)
			if v != "" {
				trimmed = append(trimmed, v)
			}
		}
		*s = trimmed
		return nil
	}

	return fmt.Errorf("StringList: invalid JSON input: %s", string(data))
}

// ToSlice returns a canonical []string value or nil if empty.
func (s StringList) ToSlice() []string {
	if s.IsZero() {
		return nil
	}
	return s
}

// String returns a comma-separated string representation.
func (s StringList) String() string {
	switch len(s) {
	case 0:
		return ""
	case 1:
		return s[0]
	default:
		return strings.Join(s, ", ")
	}
}

// ContactPoint represents a Schema.org ContactPoint object
// For more details about the meaning of the properties see: https://schema.org/ContactPoint
type ContactPoint struct {
	Type              string     `json:"@type"`
	Telephone         string     `json:"telephone,omitempty"`
	ContactType       string     `json:"contactType,omitempty"`
	ContactOption     StringList `json:"contactOption,omitempty"`
	AreaServed        StringList `json:"areaServed,omitempty"`
	AvailableLanguage string     `json:"availableLanguage,omitempty"`
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
