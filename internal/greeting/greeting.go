package greeting

import (
	"fmt"
	"time"
)

// Greeting represents a greeting message
type Greeting struct {
	Message   string    `json:"message"`
	Language  string    `json:"language,omitempty"`
	Emoji     string    `json:"emoji,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// Options configures how a greeting is generated
type Options struct {
	Name         string
	Language     string
	IncludeEmoji bool
}

// Language translations
var translations = map[string]struct {
	template string
	hello    string
	emoji    string
}{
	"en": {"Hello there, %s! 🎉", "Hello, World!", "👋"},
	"es": {"¡Hola, %s!", "¡Hola, Mundo!", "👋"},
	"fr": {"Bonjour, %s!", "Bonjour le monde!", "👋"},
	"de": {"Hallo, %s!", "Hallo, Welt!", "👋"},
	"ja": {"こんにちは、%sさん！", "こんにちは、世界！", "🇯🇵"},
	"zh": {"你好，%s！", "你好，世界！", "🇨🇳"},
}

// Generate creates a greeting based on the given options
func Generate(opts Options) *Greeting {
	greeting := &Greeting{
		Timestamp: time.Now(),
		Language:  opts.Language,
	}

	// Get language data, default to English if not found
	lang, ok := translations[opts.Language]
	if !ok {
		lang = translations["en"]
		greeting.Language = "en"
	}

	// Generate the message
	if opts.Name != "" {
		greeting.Message = fmt.Sprintf(lang.template, opts.Name)
	} else {
		greeting.Message = lang.hello
	}

	// Add emoji if requested
	if opts.IncludeEmoji {
		greeting.Emoji = lang.emoji
		greeting.Message = fmt.Sprintf("%s %s", greeting.Emoji, greeting.Message)
	}

	return greeting
}

// GetSupportedLanguages returns all supported language codes
func GetSupportedLanguages() []string {
	languages := make([]string, 0, len(translations))
	for code := range translations {
		languages = append(languages, code)
	}
	return languages
}
