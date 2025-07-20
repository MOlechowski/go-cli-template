package json

import (
	"encoding/json"
	"fmt"
)

// PrettyPrintJSON converts an object to pretty-printed JSON
func PrettyPrintJSON(v interface{}) (string, error) {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(bytes), nil
}

// PrintJSON converts an object to JSON (non-pretty)
func PrintJSON(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(bytes), nil
}
