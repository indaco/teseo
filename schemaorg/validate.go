package schemaorg

import "log"

// SchemaValidator defines an interface for types that can be validated for schema.org compliance.
type SchemaValidator interface {
	// Validate returns a slice of warning messages for missing recommended or required fields.
	Validate() []string
}

func LogValidationWarnings(v SchemaValidator) {
	for _, warning := range v.Validate() {
		log.Printf("schema warning: %s", warning)
	}
}
