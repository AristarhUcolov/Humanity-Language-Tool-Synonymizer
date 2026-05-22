package humanize

import (
	"strings"
	"testing"
)

func TestDetectLanguage(t *testing.T) {
	cases := []struct {
		text string
		want Language
	}{
		{"The good man works hard every day.", LangEnglish},
		{"Хороший человек работает каждый день и читает книги.", LangRussian},
		{"Omul bun lucrează în fiecare zi și citește cărți.", LangMoldovan},
		{"", LangEnglish},
		{"12345 67890", LangEnglish},
	}
	for _, c := range cases {
		got := DetectLanguage(c.text)
		if got != c.want {
			t.Errorf("DetectLanguage(%q) = %q, want %q", c.text, got, c.want)
		}
	}
}

func TestProcessAppliesSynonyms(t *testing.T) {
	res := Process(Input{Text: "The good man works.", Language: LangEnglish})
	if res.SynonymsApplied == 0 {
		t.Fatalf("expected synonyms to be applied, got 0")
	}
	if res.Output == "The good man works." {
		t.Errorf("output unchanged: %q", res.Output)
	}
	if !strings.HasPrefix(res.Output, "The ") {
		t.Errorf("first word capitalisation lost: %q", res.Output)
	}
}

func TestProcessRunPreservesEdgeSpaces(t *testing.T) {
	// A DOCX run fragment with leading/trailing spaces must keep them.
	cases := []string{" good ", "good ", " good", "  good  "}
	for _, in := range cases {
		res := ProcessRun(Input{Text: in, Language: LangEnglish})
		if !strings.HasPrefix(res.Output, strings.Repeat(" ", leadingSpaces(in))) {
			t.Errorf("ProcessRun(%q): lost leading space → %q", in, res.Output)
		}
		if !strings.HasSuffix(res.Output, strings.Repeat(" ", trailingSpaces(in))) {
			t.Errorf("ProcessRun(%q): lost trailing space → %q", in, res.Output)
		}
	}
}

func TestProcessRunNoCapitalization(t *testing.T) {
	// ProcessRun must NOT force-capitalise a mid-sentence fragment.
	res := ProcessRun(Input{Text: "man works", Language: LangEnglish})
	if len(res.Output) > 0 && res.Output[0] >= 'A' && res.Output[0] <= 'Z' {
		t.Errorf("ProcessRun wrongly capitalised fragment: %q", res.Output)
	}
}

func TestProcessRunPreservesWordCase(t *testing.T) {
	// Capitalised input word → capitalised synonym.
	res := ProcessRun(Input{Text: "Good", Language: LangEnglish})
	if len(res.Output) == 0 || res.Output[0] < 'A' || res.Output[0] > 'Z' {
		t.Errorf("ProcessRun lost word capitalisation: %q", res.Output)
	}
}

func TestReplaceMapCasePreserving(t *testing.T) {
	m := map[string]string{"good": "excellent"}
	got, n := replaceMap("Good and good", m)
	if n != 2 {
		t.Errorf("expected 2 replacements, got %d", n)
	}
	if !strings.Contains(got, "Excellent") || !strings.Contains(got, "excellent") {
		t.Errorf("case not preserved: %q", got)
	}
}

func TestReplaceMapWordBoundary(t *testing.T) {
	// "good" must not match inside "goodness".
	m := map[string]string{"good": "excellent"}
	got, n := replaceMap("goodness", m)
	if n != 0 || got != "goodness" {
		t.Errorf("matched inside word: %q (n=%d)", got, n)
	}
}

func TestProcessRunRussian(t *testing.T) {
	res := ProcessRun(Input{Text: "хороший", Language: LangRussian})
	if res.Output == "хороший" {
		t.Errorf("Russian synonym not applied: %q", res.Output)
	}
}

func TestProcessRunMoldovan(t *testing.T) {
	res := ProcessRun(Input{Text: "bun", Language: LangMoldovan})
	if res.Output == "bun" {
		t.Errorf("Moldovan synonym not applied: %q", res.Output)
	}
}

func TestAnalyseCounts(t *testing.T) {
	a := Analyse("Hello world. This is a test.")
	if a.Words != 6 {
		t.Errorf("Words = %d, want 6", a.Words)
	}
	if a.Sentences != 2 {
		t.Errorf("Sentences = %d, want 2", a.Sentences)
	}
}

// helpers
func leadingSpaces(s string) int {
	n := 0
	for _, r := range s {
		if r != ' ' {
			break
		}
		n++
	}
	return n
}

func trailingSpaces(s string) int {
	n := 0
	for i := len(s) - 1; i >= 0 && s[i] == ' '; i-- {
		n++
	}
	return n
}
