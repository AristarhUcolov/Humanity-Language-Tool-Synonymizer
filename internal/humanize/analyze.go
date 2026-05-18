package humanize

import (
	"strings"
	"unicode"
)

// Analysis summarises shape-level signals about a piece of input text. It is
// not a classifier — it is a budget of how AI-flavoured the text looks, used
// to drive adaptive decisions in the pipeline (and surfaced in the UI).
type Analysis struct {
	Chars         int
	Words         int
	Sentences     int
	AvgWordLen    float64
	AvgSentLen    float64

	ContractionDensity float64 // contractions per 100 words
	PassiveDensity     float64 // past-participle-by hits per 100 words (rough)
	ClicheHits         int     // count of AI-cliche signal words present
	AIScore            int     // 0–100 heuristic — higher = more AI-like
}

// Analyse runs cheap heuristics over the text.
func Analyse(text string) Analysis {
	a := Analysis{Chars: len(text)}
	if text == "" {
		return a
	}

	lower := strings.ToLower(text)
	words := strings.Fields(text)
	a.Words = len(words)

	var totalLen int
	contr := 0
	for _, w := range words {
		totalLen += len([]rune(w))
		if strings.ContainsRune(w, '\'') || strings.ContainsRune(w, '’') {
			contr++
		}
	}
	if a.Words > 0 {
		a.AvgWordLen = float64(totalLen) / float64(a.Words)
		a.ContractionDensity = 100.0 * float64(contr) / float64(a.Words)
	}

	a.Sentences = countSentences(text)
	if a.Sentences > 0 {
		a.AvgSentLen = float64(a.Words) / float64(a.Sentences)
	}

	// Passive-ish marker: "X was Yed by" / "is being Yed". A rough probe.
	a.PassiveDensity = 100.0 * float64(countAny(lower, passiveProbes)) / max(1.0, float64(a.Words))

	a.ClicheHits = countAny(lower, aiSignalWords)

	a.AIScore = computeAIScore(a)
	return a
}

func countSentences(text string) int {
	n := 0
	for _, r := range text {
		if r == '.' || r == '!' || r == '?' || r == '…' {
			n++
		}
	}
	if n == 0 && strings.TrimSpace(text) != "" {
		return 1
	}
	return n
}

func countAny(lower string, needles []string) int {
	n := 0
	for _, s := range needles {
		n += strings.Count(lower, s)
	}
	return n
}

func computeAIScore(a Analysis) int {
	if a.Words < 5 {
		return 0
	}
	// Each component contributes 0–40, capped, then summed and clamped.
	cliche := minF(40, float64(a.ClicheHits)*8)
	density := 0.0
	if a.Words > 0 {
		density = 100.0 * float64(a.ClicheHits) / float64(a.Words) // hits per 100 words
	}
	densityPts := minF(30, density*5)
	longSent := 0.0
	if a.AvgSentLen > 22 {
		longSent = minF(15, (a.AvgSentLen-22)*1.5)
	}
	longWord := 0.0
	if a.AvgWordLen > 5.5 {
		longWord = minF(15, (a.AvgWordLen-5.5)*5)
	}
	score := cliche + densityPts + longSent + longWord
	if score > 100 {
		score = 100
	}
	return int(score)
}

func minF(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// CapitalizedRatio returns the fraction of words that start with a capital
// letter (rough formality proxy used elsewhere).
func CapitalizedRatio(text string) float64 {
	words := strings.Fields(text)
	if len(words) == 0 {
		return 0
	}
	caps := 0
	for _, w := range words {
		r := []rune(w)
		if len(r) > 0 && unicode.IsUpper(r[0]) {
			caps++
		}
	}
	return float64(caps) / float64(len(words))
}

// passiveProbes is a small set of substrings that suggest passive voice or
// AI-flavoured constructions. Used only as a heuristic, not a classifier.
var passiveProbes = []string{
	" was made by ", " is being made", " was used to ", " is used to ",
	" can be seen", " is considered", " are considered", " was conducted",
	" is conducted", " has been ", " have been ",
}

// aiSignalWords is the short list of strong AI tells we count for scoring.
// Keep this list small and high-precision; the full rewrite list lives in
// dict_en.go.
var aiSignalWords = []string{
	"furthermore", "moreover", "additionally", "in conclusion", "delve",
	"dive into", "navigating", "tapestry", "landscape of", "realm of",
	"in the realm of", "in today's", "it is important to note",
	"it's important to note", "plays a crucial role", "plays a key role",
	"plays a vital role", "plays a pivotal role", "cannot be overstated",
	"a testament to", "ever-evolving", "ever-changing", "cutting-edge",
	"state-of-the-art", "game-changer", "game-changing", "paradigm shift",
	"unlock the", "harness the power", "leverage the power",
	"a wealth of", "a plethora of", "a myriad of", "myriad", "multitude",
	"groundbreaking", "revolutionary", "transformative", "seamless",
	"holistic", "robust", "as we all know", "needless to say",
	"at the end of the day", "when it comes to", "in terms of",
	"with regards to", "let's dive", "let us dive",
}
