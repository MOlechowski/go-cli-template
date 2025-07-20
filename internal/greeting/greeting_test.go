package greeting

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name    string
		opts    Options
		wantMsg string
	}{
		{
			name: "basic hello world",
			opts: Options{
				Language: "en",
			},
			wantMsg: "Hello, World!",
		},
		{
			name: "hello world with emoji",
			opts: Options{
				Language:     "en",
				IncludeEmoji: true,
			},
			wantMsg: "👋 Hello, World!",
		},
		{
			name: "personalized greeting",
			opts: Options{
				Name:     "Alice",
				Language: "en",
			},
			wantMsg: "Hello, Alice!",
		},
		{
			name: "spanish greeting",
			opts: Options{
				Name:     "Carlos",
				Language: "es",
			},
			wantMsg: "¡Hola, Carlos!",
		},
		{
			name: "japanese greeting with emoji",
			opts: Options{
				Name:         "Tanaka",
				Language:     "ja",
				IncludeEmoji: true,
			},
			wantMsg: "🇯🇵 こんにちは、Tanakaさん！",
		},
		{
			name: "unknown language falls back to english",
			opts: Options{
				Name:     "Test",
				Language: "unknown",
			},
			wantMsg: "Hello, Test!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			greeting := Generate(tt.opts)
			if greeting.Message != tt.wantMsg {
				t.Errorf("Generate() message = %v, want %v", greeting.Message, tt.wantMsg)
			}
			if tt.opts.Language != "unknown" && greeting.Language != tt.opts.Language {
				t.Errorf("Generate() language = %v, want %v", greeting.Language, tt.opts.Language)
			}
		})
	}
}

func TestGetSupportedLanguages(t *testing.T) {
	langs := GetSupportedLanguages()
	if len(langs) == 0 {
		t.Error("GetSupportedLanguages() returned empty slice")
	}

	// Check that common languages are supported
	expectedLangs := []string{"en", "es", "fr", "de", "ja", "zh"}
	for _, expected := range expectedLangs {
		found := false
		for _, lang := range langs {
			if lang == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected language %s not found in supported languages", expected)
		}
	}
}
