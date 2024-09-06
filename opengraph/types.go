package opengraph

// OpenGraphObject represents common Open Graph metadata.
type OpenGraphObject struct {
	Type        string // og:type, the type of the object
	Title       string // og:title, the title of the object
	URL         string // og:url, the canonical URL of the object
	Description string // og:description, a brief description of the object
	Image       string // og:image, URL to the image of the object
}

// ensureDefaults sets default values for OpenGraphObject if they are not already set.
func (og *OpenGraphObject) ensureDefaults(defaultType string) {
	if og.Type == "" {
		og.Type = defaultType
	}
}
