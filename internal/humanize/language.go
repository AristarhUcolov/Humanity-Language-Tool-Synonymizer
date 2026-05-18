package humanize

import (
	"strings"
	"unicode"
)

// Language represents a supported language for synonym replacement.
type Language string

const (
	LangEnglish  Language = "en"
	LangRussian  Language = "ru"
	LangMoldovan Language = "mo"
	LangAuto     Language = "auto"
)

// AllLanguages returns the list of supported languages with display names.
func AllLanguages() []LanguageInfo {
	return []LanguageInfo{
		{Code: LangAuto, Name: "Auto-detect", Native: "Auto"},
		{Code: LangEnglish, Name: "English", Native: "English"},
		{Code: LangRussian, Name: "Russian", Native: "Русский"},
		{Code: LangMoldovan, Name: "Moldovan", Native: "Moldovenească"},
	}
}

// LanguageInfo describes a supported language.
type LanguageInfo struct {
	Code   Language `json:"code"`
	Name   string   `json:"name"`
	Native string   `json:"native"`
}

// DetectLanguage analyzes text and returns the most likely language.
// Uses script detection and character frequency analysis.
func DetectLanguage(text string) Language {
	if text == "" {
		return LangEnglish
	}

	cyrillicCount := 0
	latinCount := 0
	diacriticCount := 0
	totalLetters := 0

	moldovanIndicators := 0
	russianIndicators := 0

	lower := strings.ToLower(text)

	for _, r := range text {
		if !unicode.IsLetter(r) {
			continue
		}
		totalLetters++

		if isCyrillic(r) {
			cyrillicCount++
		} else if unicode.Is(unicode.Latin, r) {
			latinCount++
			if isMoldovaDiacritic(r) {
				diacriticCount++
			}
		}
	}

	if totalLetters == 0 {
		return LangEnglish
	}

	// Count language-specific markers
	for _, marker := range russianMarkers {
		russianIndicators += strings.Count(lower, marker)
	}
	for _, marker := range moldovanMarkers {
		moldovanIndicators += strings.Count(lower, marker)
	}

	cyrillicRatio := float64(cyrillicCount) / float64(totalLetters)
	diacriticRatio := float64(diacriticCount) / float64(totalLetters)

	// Decision logic
	if cyrillicRatio > 0.5 {
		return LangRussian
	}
	if diacriticRatio > 0.02 || moldovanIndicators > russianIndicators {
		// Moldovan/Romanian uses ă, â, î, ș, ț
		return LangMoldovan
	}
	return LangEnglish
}

// isCyrillic returns true if rune is in Cyrillic script.
func isCyrillic(r rune) bool {
	return unicode.Is(unicode.Cyrillic, r)
}

// isMoldovaDiacritic returns true for letters specific to Moldovan/Romanian.
func isMoldovaDiacritic(r rune) bool {
	switch r {
	case 'ă', 'Ă', 'â', 'Â', 'î', 'Î', 'ș', 'Ș', 'ț', 'Ț', 'ş', 'Ş', 'ţ', 'Ţ':
		return true
	}
	return false
}

// russianMarkers are common Russian function words used for detection.
var russianMarkers = []string{
	" и ", " в ", " на ", " не ", " что ", " это ", " с ", " по ",
	" к ", " о ", " от ", " для ", " как ", " но ", " или ", " если ",
	" же ", " бы ", " ли ", " только ",
}

// moldovanMarkers are common Romanian/Moldovan function words.
var moldovanMarkers = []string{
	" și ", " de ", " la ", " în ", " pe ", " cu ", " un ", " o ",
	" este ", " sunt ", " nu ", " ce ", " care ", " pentru ", " dar ",
	" sau ", " dacă ", " când ", " unde ", " cum ",
}

// pickPack returns the dictionary pack for the given language.
// Falls back to English if the requested language has no dictionary.
func pickPack(lang Language) dictPack {
	switch lang {
	case LangRussian:
		return synPackRu
	case LangMoldovan:
		return synPackMo
	default:
		return synPack
	}
}
