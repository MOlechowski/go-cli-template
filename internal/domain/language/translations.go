package language

// initializeLanguages sets up all supported languages
func initializeLanguages() map[string]Language {
	return map[string]Language{
		"en": {
			Code:             "en",
			Name:             "English",
			GreetingTemplate: "Hello, %s!",
			DefaultGreeting:  "Hello, World!",
			Emoji:            "👋",
		},
		"es": {
			Code:             "es",
			Name:             "Spanish",
			GreetingTemplate: "¡Hola, %s!",
			DefaultGreeting:  "¡Hola, Mundo!",
			Emoji:            "🌍",
		},
		"fr": {
			Code:             "fr",
			Name:             "French",
			GreetingTemplate: "Bonjour, %s!",
			DefaultGreeting:  "Bonjour, le Monde!",
			Emoji:            "🇫🇷",
		},
		"de": {
			Code:             "de",
			Name:             "German",
			GreetingTemplate: "Hallo, %s!",
			DefaultGreeting:  "Hallo, Welt!",
			Emoji:            "🇩🇪",
		},
		"ja": {
			Code:             "ja",
			Name:             "Japanese",
			GreetingTemplate: "こんにちは、%sさん！",
			DefaultGreeting:  "こんにちは、世界！",
			Emoji:            "🇯🇵",
		},
		"zh": {
			Code:             "zh",
			Name:             "Chinese",
			GreetingTemplate: "你好，%s！",
			DefaultGreeting:  "你好，世界！",
			Emoji:            "🇨🇳",
		},
	}
}
