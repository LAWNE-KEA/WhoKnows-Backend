package monitoring

import (
	"html"
	"reflect"
	"strings"
	"unicode/utf8"
)

// SanitizeValue sanitizes the input string by performing the following operations:
//  1. Removes control characters (non-printable ASCII characters 0-31 and delete character 127).
//  2. Normalizes whitespace by trimming leading and trailing spaces and replacing multiple spaces with a single space.
//  3. Escapes HTML and special characters (backslash, double quote, single quote) to prevent injection attacks.
//
// Parameters:
//  - value: The input string to be sanitized.
//
// Returns:
//  - A sanitized string with control characters removed, whitespace normalized, and special characters escaped.

func SanitizeValue(value string) string {
	if !utf8.ValidString(value) {
		value = strings.ToValidUTF8(value, "")
	}

	var builder strings.Builder
	for _, r := range value {
		if r >= 32 && r != 127 { // Skip non-printable ASCII (0-31) and delete (127)
			builder.WriteRune(r)
		}
	}

	value = strings.TrimSpace(value)
	value = strings.Join(strings.Fields(value), " ")
	value = html.EscapeString(value)
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "\"", "\\\"")
	value = strings.ReplaceAll(value, "'", "\\'")
	return value
}

// SanitizeFields sanitizes a map of fields, returning a new map with sanitized string values.
func SanitizeFields(fields map[string]interface{}) map[string]interface{} {
	sanitized := make(map[string]interface{})
	for key, value := range fields {
		if str, ok := value.(string); ok {
			sanitized[key] = SanitizeValue(str)
		} else {
			sanitized[key] = value
		}
	}
	return sanitized
}

func SanitizeStruct(obj interface{}) {
	val := reflect.ValueOf(obj)

	// Ensure obj is a pointer to a struct
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		panic("SanitizeStruct requires a pointer to a struct")
	}

	val = val.Elem() // Dereference the pointer
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}
		sanitizeField(field)
	}
}

func sanitizeField(field reflect.Value) {
	// Sanitize string fields
	if field.Kind() == reflect.String {
		field.SetString(SanitizeValue(field.String()))
	}

	// Optionally handle nested structs
	if field.Kind() == reflect.Struct {
		SanitizeStruct(field.Addr().Interface())
	}

	// Optionally handle slices of strings or structs
	if field.Kind() == reflect.Slice {
		sanitizeSlice(field)
	}
}

func sanitizeSlice(field reflect.Value) {
	if field.IsNil() {
		return
	}

	sliceElemType := field.Type().Elem()
	if sliceElemType.Kind() == reflect.String {
		for j := 0; j < field.Len(); j++ {
			field.Index(j).SetString(SanitizeValue(field.Index(j).String()))
		}
	} else if sliceElemType.Kind() == reflect.Struct {
		for j := 0; j < field.Len(); j++ {
			SanitizeStruct(field.Index(j).Addr().Interface())
		}
	}
}