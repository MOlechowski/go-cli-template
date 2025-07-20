package greeting

import (
	"fmt"
	"time"

	"github.com/go-cli-template/hello-world-cli/internal/domain/language"
)

// Service handles greeting generation business logic
type Service struct {
	languageService *language.Service
}

// NewService creates a new greeting service
func NewService(langService *language.Service) *Service {
	return &Service{
		languageService: langService,
	}
}

// GenerateGreeting creates a greeting based on options
func (s *Service) GenerateGreeting(opts Options) (*Greeting, error) {
	greeting := &Greeting{
		Timestamp: time.Now(),
		Language:  opts.Language,
	}

	// Get the appropriate greeting template
	template := s.languageService.GetGreetingTemplate(opts.Language)

	// Generate the message
	if opts.Name != "" {
		greeting.Message = fmt.Sprintf(template, opts.Name)
	} else {
		greeting.Message = s.languageService.GetDefaultGreeting(opts.Language)
	}

	// Add emoji if requested
	if opts.IncludeEmoji {
		greeting.Emoji = s.languageService.GetGreetingEmoji(opts.Language)
		greeting.Message = fmt.Sprintf("%s %s", greeting.Emoji, greeting.Message)
	}

	return greeting, nil
}

// GetSupportedLanguages returns all supported language codes
func (s *Service) GetSupportedLanguages() []string {
	return s.languageService.GetSupportedLanguages()
}
