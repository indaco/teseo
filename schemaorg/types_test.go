package schemaorg

import (
	"encoding/json"
	"strings"
	"testing"
)

// Struct that uses the custom StringList type
type MyStruct struct {
	AreaServed StringList `json:"areaServed,omitempty"`
}

func TestStringList_IsZero(t *testing.T) {
	tests := []struct {
		name  string
		input StringList
		want  bool
	}{
		{"nil slice", nil, true},
		{"empty slice", StringList{}, true},
		{"single empty string", StringList{""}, true},
		{"single non-empty", StringList{"x"}, false},
		{"multiple values", StringList{"x", "y"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.IsZero()
			if got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringList_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    StringList
		expected string
	}{
		{"nil slice", nil, `null`},
		{"empty slice", StringList{}, `null`},
		{"slice with empty string", StringList{""}, `null`},
		{"single value", StringList{"foo"}, `"foo"`},
		{"multiple values", StringList{"foo", "bar"}, `["foo","bar"]`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("MarshalJSON() error = %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("MarshalJSON() = %s, want %s", data, tt.expected)
			}
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    MyStruct
		expected string
	}{

		{
			name:     "Non-empty slice",
			input:    MyStruct{AreaServed: StringList{"hello", "world"}},
			expected: `{"areaServed":["hello","world"]}`,
		},
		{
			name:     "String",
			input:    MyStruct{AreaServed: StringList{"hello"}},
			expected: `{"areaServed":"hello"}`,
		},
		{
			name:     "Empty slice (should be omitted)",
			input:    MyStruct{AreaServed: StringList{}},
			expected: `{}`,
		},
		{
			name:     "Nil slice (should be omitted)",
			input:    MyStruct{},
			expected: `{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.input)
			if err != nil {
				t.Fatalf("unexpected error marshaling: %v", err)
			}

			if string(data) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, data)
			}
		})
	}
}

func TestStringList_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected StringList
		wantErr  bool
	}{
		{"null input", `null`, nil, false},
		{"empty string", `""`, StringList{""}, false},
		{"trimmed string", `"  foo  "`, StringList{"foo"}, false},
		{"single non-empty", `"bar"`, StringList{"bar"}, false},
		{"empty array", `[]`, StringList{}, false},
		{"array of strings", `["a", "b", "c"]`, StringList{"a", "b", "c"}, false},
		{"array with spaces", `["  x  ", " y "]`, StringList{"x", "y"}, false},
		{"array with empty values", `["", " "]`, StringList{}, false},
		{"invalid type", `123`, nil, true},
		{"invalid json", `{`, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got StringList
			err := json.Unmarshal([]byte(tt.input), &got)
			if (err != nil) != tt.wantErr {
				t.Fatalf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !equalSlices(got, tt.expected) {
				t.Errorf("UnmarshalJSON() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStringList_ToSlice(t *testing.T) {
	tests := []struct {
		name  string
		input StringList
		want  []string
	}{
		{"nil", nil, nil},
		{"empty", StringList{}, nil},
		{"single empty string", StringList{""}, nil},
		{"non-empty", StringList{"a", "b"}, []string{"a", "b"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.ToSlice()
			if !equalSlices(got, tt.want) {
				t.Errorf("ToSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringList_String(t *testing.T) {
	tests := []struct {
		name  string
		input StringList
		want  string
	}{
		{"nil", nil, ""},
		{"empty", StringList{}, ""},
		{"single", StringList{"foo"}, "foo"},
		{"multi", StringList{"a", "b", "c"}, "a, b, c"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.String()
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func equalSlices(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if strings.TrimSpace(v) != strings.TrimSpace(b[i]) {
			return false
		}
	}
	return true
}
