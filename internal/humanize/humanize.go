package humanize

import (
	"fmt"
)

// Input is the request to the synonymaizer pipeline.
type Input struct {
	Text     string
	Language Language // optional, auto-detected if empty or LangAuto
}

// Result is what the pipeline returns. Output is the rewritten text with
// synonyms applied; the rest is telemetry for the UI.
type Result struct {
	Output string

	Before Analysis
	After  Analysis

	DetectedLanguage Language

	SynonymsApplied int
	PunctFixes      int
	SentencesSplit  int

	Notes []string
}

// Process runs the synonymaizer pipeline. It is safe to call from many
// goroutines (no shared mutable state beyond the regex cache, which is
// protected by a mutex inside replace.go).
func Process(in Input) Result {
	text := in.Text
	res := Result{Before: Analyse(text)}

	// 1. Detect or use specified language
	lang := in.Language
	if lang == "" || lang == LangAuto {
		lang = DetectLanguage(text)
	}
	res.DetectedLanguage = lang

	// 2. Pick the dictionary pack for the detected language
	p := pickPack(lang)

	// 3. Apply synonym replacements to elevate vocabulary.
	var n int
	text, n = replaceMap(text, p.AICliches)
	res.SynonymsApplied = n

	// 4. Adaptive: only split monster sentences (> 28 words). English-style logic
	// works reasonably for other Latin scripts; less effective for Russian but harmless.
	text, res.SentencesSplit = softenLongSentences(text, 28)

	// 5. Punctuation / dash / quote normalisation.
	text, res.PunctFixes = normalizePunctuation(text)
	text = normalizePunctuationForLang(text, lang)

	// 6. Mop up orphan commas / periods left behind by deletions.
	text = fixOrphanPunctuation(text)

	// 7. Collapse runs of spaces and trim line endings.
	text = collapseSpaces(text)

	// 8. Re-capitalise sentence starts.
	text = capitalizeSentences(text)

	res.Output = text
	res.After = Analyse(text)

	res.Notes = buildNotes(res)
	return res
}

// ProcessRun is the synonymaizer pipeline for a single DOCX text run (one
// <w:t> node). Unlike Process, it applies ONLY synonym replacement and does
// NOT touch structure:
//
//   - Leading/trailing whitespace is preserved exactly (no collapseSpaces).
//     A DOCX paragraph is often split across several runs — e.g. "The ",
//     "good", " man" — and trimming a run's edge spaces would glue words
//     together ("Theexcellentman").
//   - Sentence capitalisation is NOT forced (no capitalizeSentences). The
//     original document already has correct capitalisation; forcing the
//     first letter of every run uppercase would wrongly capitalise mid-
//     sentence runs. Word-level case is still preserved by replaceMap's
//     preserveCase (so "Good" → "Excellent", "good" → "excellent").
//   - Long sentences are NOT split and punctuation is NOT normalised: a run
//     is a fragment, not a sentence, so those operations would corrupt it.
//
// The result: the run's text is identical to the original except that
// individual words are swapped for synonyms — exact structural fidelity.
func ProcessRun(in Input) Result {
	text := in.Text
	res := Result{Before: Analyse(text)}

	lang := in.Language
	if lang == "" || lang == LangAuto {
		lang = DetectLanguage(text)
	}
	res.DetectedLanguage = lang

	p := pickPack(lang)

	var n int
	text, n = replaceMap(text, p.AICliches)
	res.SynonymsApplied = n

	res.Output = text
	res.After = Analyse(text)
	res.Notes = buildNotes(res)
	return res
}

func buildNotes(r Result) []string {
	var notes []string

	// Language detection note
	if r.DetectedLanguage != "" {
		langName := languageDisplayName(r.DetectedLanguage)
		notes = append(notes, fmt.Sprintf("Language: %s", langName))
	}

	if r.SynonymsApplied > 0 {
		notes = append(notes, fmt.Sprintf("Synonyms applied: %d", r.SynonymsApplied))
	}
	if r.SentencesSplit > 0 {
		notes = append(notes, fmt.Sprintf("Long sentences split: %d", r.SentencesSplit))
	}
	if r.PunctFixes > 0 {
		notes = append(notes, fmt.Sprintf("Punctuation fixes: %d", r.PunctFixes))
	}
	// Track word-count increase from using longer synonyms.
	wordsDelta := r.After.Words - r.Before.Words
	if wordsDelta > 0 {
		notes = append(notes, fmt.Sprintf("Word count: %d → %d (+%d)", r.Before.Words, r.After.Words, wordsDelta))
	} else if wordsDelta < 0 {
		notes = append(notes, fmt.Sprintf("Word count: %d → %d (%d)", r.Before.Words, r.After.Words, wordsDelta))
	}
	if r.SynonymsApplied == 0 {
		notes = append(notes, "No synonyms matched — vocabulary already optimised.")
	}
	return notes
}

func languageDisplayName(lang Language) string {
	for _, info := range AllLanguages() {
		if info.Code == lang {
			return info.Native
		}
	}
	return string(lang)
}
