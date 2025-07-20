package language

// Language represents a supported language
type Language struct {
	Code            string
	Name            string
	GreetingTemplate string // Template with %s for name
	DefaultGreeting  string // Default greeting without name
	Emoji           string
}