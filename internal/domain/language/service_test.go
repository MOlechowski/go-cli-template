package language

import (
	"testing"
)

func TestGetGreetingTemplate(t *testing.T) {
	service := NewService()

	tests := []struct {
		code     string
		expected string
	}{
		{"en", "Hello, %s!"},
		{"es", "¡Hola, %s!"},
		{"fr", "Bonjour, %s!"},
		{"unknown", "Hello, %s!"}, // Should default to English
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			result := service.GetGreetingTemplate(tt.code)
			if result != tt.expected {
				t.Errorf("GetGreetingTemplate(%s) = %v, want %v", tt.code, result, tt.expected)
			}
		})
	}
}

func TestGetDefaultGreeting(t *testing.T) {
	service := NewService()

	tests := []struct {
		code     string
		expected string
	}{
		{"en", "Hello, World!"},
		{"es", "¡Hola, Mundo!"},
		{"ja", "こんにちは、世界！"},
		{"unknown", "Hello, World!"}, // Should default to English
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			result := service.GetDefaultGreeting(tt.code)
			if result != tt.expected {
				t.Errorf("GetDefaultGreeting(%s) = %v, want %v", tt.code, result, tt.expected)
			}
		})
	}
}

func TestIsSupported(t *testing.T) {
	service := NewService()

	tests := []struct {
		code     string
		expected bool
	}{
		{"en", true},
		{"es", true},
		{"fr", true},
		{"de", true},
		{"ja", true},
		{"zh", true},
		{"unknown", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			result := service.IsSupported(tt.code)
			if result != tt.expected {
				t.Errorf("IsSupported(%s) = %v, want %v", tt.code, result, tt.expected)
			}
		})
	}
}