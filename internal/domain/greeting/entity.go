package greeting

import "time"

// Greeting represents a greeting message entity
type Greeting struct {
	Message   string    `json:"message"`
	Language  string    `json:"language,omitempty"`
	Emoji     string    `json:"emoji,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// GreetingOptions configures how a greeting is generated
type GreetingOptions struct {
	Name         string
	Language     string
	IncludeEmoji bool
	Format       string // "text" or "json"
}