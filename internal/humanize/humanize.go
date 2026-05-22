package humanize

import (
	"fmt"
	"strings"
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
//
// The input is processed line by line: every newline in the original text is
// kept exactly where it was, so paragraph and line structure survive intact.
// Sentence-level work (synonym swaps, long-sentence splitting, punctuation
// normalisation, capitalisation) happens within each line.
func Process(in Input) Result {
	text := in.Text
	res := Result{Before: Analyse(text)}

	// 1. Detect or use specified language (from the whole text).
	lang := in.Language
	if lang == "" || lang == LangAuto {
		lang = DetectLanguage(text)
	}
	res.DetectedLanguage = lang

	// 2. Pick the dictionary pack for the detected language.
	p := pickPack(lang)

	// 3. Process each line independently, preserving every newline.
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue // keep blank lines (paragraph gaps) exactly as-is
		}
		lines[i] = processLine(line, p, lang, &res)
	}
	res.Output = strings.Join(lines, "\n")

	res.After = Analyse(res.Output)
	res.Notes = buildNotes(res)
	return res
}

// processLine runs the sentence-level pipeline on a single line of text and
// accumulates telemetry into res.
func processLine(line string, p dictPack, lang Language, res *Result) string {
	// a. Apply synonym replacements to elevate vocabulary.
	line, n := replaceMap(line, p.AICliches)
	res.SynonymsApplied += n

	// b. Split monster sentences (> 28 words) on a coordinating conjunction.
	line, split := softenLongSentences(line, 28)
	res.SentencesSplit += split

	// c. Punctuation / dash / quote normalisation.
	line, pf := normalizePunctuation(line)
	res.PunctFixes += pf
	line = normalizePunctuationForLang(line, lang)

	// d. Mop up orphan commas / periods left behind by deletions.
	line = fixOrphanPunctuation(line)

	// e. Collapse runs of spaces; trim only the trailing edge (a leading
	//    indent is meaningful structure and is kept).
	line = collapseInlineSpaces(line)

	// f. Re-capitalise sentence starts within the line.
	line = capitalizeSentences(line)
	return line
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
