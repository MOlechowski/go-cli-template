package greeting

import (
	"testing"

	"github.com/go-cli-template/hello-world-cli/internal/domain/language"
)

func TestGenerateGreeting(t *testing.T) {
	langService := language.NewService()
	service := NewService(langService)

	tests := []struct {
		name     string
		opts     GreetingOptions
		wantMsg  string
		wantErr  bool
	}{
		{
			name: "basic hello world",
			opts: GreetingOptions{
				Language: "en",
			},
			wantMsg: "Hello, World!",
		},
		{
			name: "hello world with emoji",
			opts: GreetingOptions{
				Language:     "en",
				IncludeEmoji: true,
			},
			wantMsg: "👋 Hello, World!",
		},
		{
			name: "personalized greeting",
			opts: GreetingOptions{
				Name:     "Alice",
				Language: "en",
			},
			wantMsg: "Hello, Alice!",
		},
		{
			name: "spanish greeting",
			opts: GreetingOptions{
				Name:     "Carlos",
				Language: "es",
			},
			wantMsg: "¡Hola, Carlos!",
		},
		{
			name: "japanese greeting with emoji",
			opts: GreetingOptions{
				Name:         "Tanaka",
				Language:     "ja",
				IncludeEmoji: true,
			},
			wantMsg: "🇯🇵 こんにちは、Tanakaさん！",
		},
		{
			name: "unknown language falls back to english",
			opts: GreetingOptions{
				Name:     "Test",
				Language: "unknown",
			},
			wantMsg: "Hello, Test!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			greeting, err := service.GenerateGreeting(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateGreeting() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if greeting.Message != tt.wantMsg {
				t.Errorf("GenerateGreeting() message = %v, want %v", greeting.Message, tt.wantMsg)
			}
			if greeting.Language != tt.opts.Language {
				t.Errorf("GenerateGreeting() language = %v, want %v", greeting.Language, tt.opts.Language)
			}
		})
	}
}

func TestGetSupportedLanguages(t *testing.T) {
	langService := language.NewService()
	service := NewService(langService)

	langs := service.GetSupportedLanguages()
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