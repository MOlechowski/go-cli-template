package language

// Service handles language-related operations
type Service struct {
	languages map[string]Language
}

// NewService creates a new language service
func NewService() *Service {
	return &Service{
		languages: initializeLanguages(),
	}
}

// GetGreetingTemplate returns the greeting template for a language
func (s *Service) GetGreetingTemplate(code string) string {
	if lang, ok := s.languages[code]; ok {
		return lang.GreetingTemplate
	}
	return s.languages["en"].GreetingTemplate
}

// GetDefaultGreeting returns the default greeting for a language
func (s *Service) GetDefaultGreeting(code string) string {
	if lang, ok := s.languages[code]; ok {
		return lang.DefaultGreeting
	}
	return s.languages["en"].DefaultGreeting
}

// GetGreetingEmoji returns the greeting emoji for a language
func (s *Service) GetGreetingEmoji(code string) string {
	if lang, ok := s.languages[code]; ok {
		return lang.Emoji
	}
	return s.languages["en"].Emoji
}

// GetSupportedLanguages returns all supported language codes
func (s *Service) GetSupportedLanguages() []string {
	codes := make([]string, 0, len(s.languages))
	for code := range s.languages {
		codes = append(codes, code)
	}
	return codes
}

// IsSupported checks if a language is supported
func (s *Service) IsSupported(code string) bool {
	_, ok := s.languages[code]
	return ok
}
