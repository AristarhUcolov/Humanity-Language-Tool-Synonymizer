package humanize

import (
	"regexp"
	"strings"
	"unicode"
)

var (
	reEmDashSpaces  = regexp.MustCompile(`\s+—\s+`)
	reEnDashSpaces  = regexp.MustCompile(`\s+–\s+`)
	reDashHyphen    = regexp.MustCompile(`(\p{L})\s+-\s+(\p{L})`) // "word - word" → "word — word"
	reTripleDots    = regexp.MustCompile(`\.{3,}`)
	reMultiSpaces   = regexp.MustCompile(`[ \t]{2,}`)
	reSpaceBefore   = regexp.MustCompile(`\s+([,.;:!?])`)
	reNoSpaceAfter  = regexp.MustCompile(`([,.;:!?])(\p{L})`)
	reMultiNewline  = regexp.MustCompile(`\n{3,}`)
	reSmartQuoteDbl = regexp.MustCompile(`[\x{201C}\x{201D}]`)
	reSmartQuoteSgl = regexp.MustCompile(`[\x{2018}\x{2019}]`)
	reBullet        = regexp.MustCompile(`(?m)^\s*[•●▪◦▫]\s*`)
	reSpaceParen    = regexp.MustCompile(`\(\s+`)         // "( word" → "(word"
	reParenSpace    = regexp.MustCompile(`\s+\)`)         // "word )" → "word)"
	reExclamSeries  = regexp.MustCompile(`!{2,}`)         // "!!!" → "!"
	reQuestSeries   = regexp.MustCompile(`\?{2,}`)        // "???" → "?"
	reMixedExclQst  = regexp.MustCompile(`[!?]{3,}`)      // "?!?!" → "?!"
	rePunctEnd      = regexp.MustCompile(`([.!?…])\s+([a-zа-яă-țА-ЯĂ-Ț])`)
	reCommaNoSpace  = regexp.MustCompile(`(\w),(\w)`)     // "a,b" → "a, b" (not for numbers)
	reNumberComma   = regexp.MustCompile(`(\d),(\d)`)     // keep "1,000"
	// Orphans created by phrase deletions:
	reDoubleComma  = regexp.MustCompile(`,\s*,`)
	reDoublePeriod = regexp.MustCompile(`\.\s*\.+`)
	reCommaPeriod  = regexp.MustCompile(`,\s*\.`)
	rePeriodComma  = regexp.MustCompile(`\.\s*,`)
	reLeadingComma = regexp.MustCompile(`(?m)^\s*,\s*`)
	// Russian/Moldovan specific:
	reCyrillicQuote = regexp.MustCompile(`"([^"]*)"`) // simple double quotes → «ёлочки» for RU/MO
)

// normalizePunctuation cleans up common typography issues:
//   - Double dashes (en/em) with proper spacing
//   - Triple dots → ellipsis
//   - Smart quotes → straight quotes
//   - Bullet chars → dash
//   - Multiple exclamations/questions → single
//   - Missing space after punctuation
//   - Extra space before punctuation
//   - Parenthesis spacing
//   - Multiple newlines collapsed
func normalizePunctuation(text string) (string, int) {
	count := 0

	if m := reEmDashSpaces.FindAllStringIndex(text, -1); m != nil {
		count += len(m)
		text = reEmDashSpaces.ReplaceAllString(text, " — ")
	}
	if m := reEnDashSpaces.FindAllStringIndex(text, -1); m != nil {
		count += len(m)
		text = reEnDashSpaces.ReplaceAllString(text, " — ")
	}
	if m := reDashHyphen.FindAllStringIndex(text, -1); m != nil {
		count += len(m)
		text = reDashHyphen.ReplaceAllString(text, "$1 — $2")
	}
	if m := reTripleDots.FindAllStringIndex(text, -1); m != nil {
		count += len(m)
		text = reTripleDots.ReplaceAllString(text, "…")
	}
	if m := reSmartQuoteDbl.FindAllStringIndex(text, -1); m != nil {
		count += len(m)
		text = reSmartQuoteDbl.ReplaceAllString(text, `"`)
	}
	if m := reSmartQuoteSgl.FindAllStringIndex(text, -1); m != nil {
		count += len(m)
		text = reSmartQuoteSgl.ReplaceAllString(text, `'`)
	}
	if m := reBullet.FindAllStringIndex(text, -1); m != nil {
		count += len(m)
		text = reBullet.ReplaceAllString(text, "- ")
	}
	if m := reExclamSeries.FindAllStringIndex(text, -1); m != nil {
		count += len(m)
		text = reExclamSeries.ReplaceAllString(text, "!")
	}
	if m := reQuestSeries.FindAllStringIndex(text, -1); m != nil {
		count += len(m)
		text = reQuestSeries.ReplaceAllString(text, "?")
	}
	if m := reSpaceParen.FindAllStringIndex(text, -1); m != nil {
		count += len(m)
		text = reSpaceParen.ReplaceAllString(text, "(")
	}
	if m := reParenSpace.FindAllStringIndex(text, -1); m != nil {
		count += len(m)
		text = reParenSpace.ReplaceAllString(text, ")")
	}
	if m := reSpaceBefore.FindAllStringIndex(text, -1); m != nil {
		count += len(m)
		text = reSpaceBefore.ReplaceAllString(text, "$1")
	}
	if m := reNoSpaceAfter.FindAllStringIndex(text, -1); m != nil {
		count += len(m)
		text = reNoSpaceAfter.ReplaceAllString(text, "$1 $2")
	}
	if m := reMultiNewline.FindAllStringIndex(text, -1); m != nil {
		count += len(m)
		text = reMultiNewline.ReplaceAllString(text, "\n\n")
	}
	return text, count
}

// normalizePunctuationForLang applies language-specific punctuation rules.
// Currently applies elochki «» for Russian/Moldovan double quotes.
func normalizePunctuationForLang(text string, lang Language) string {
	if lang == LangRussian || lang == LangMoldovan {
		text = reCyrillicQuote.ReplaceAllString(text, "«$1»")
	}
	return text
}

// fixOrphanPunctuation cleans up garbage left after phrase deletions:
// "foo, , bar" → "foo, bar"; ". ." → "."; leading orphan commas.
func fixOrphanPunctuation(text string) string {
	for i := 0; i < 3; i++ { // converge in a couple of passes
		prev := text
		text = reDoubleComma.ReplaceAllString(text, ",")
		text = reDoublePeriod.ReplaceAllString(text, ".")
		text = reCommaPeriod.ReplaceAllString(text, ".")
		text = rePeriodComma.ReplaceAllString(text, ".")
		text = reLeadingComma.ReplaceAllString(text, "")
		if text == prev {
			break
		}
	}
	return text
}

// capitalizeSentences uppercases the first letter of the text and of every
// sentence that follows a terminal punctuation mark.
// Handles Cyrillic, Latin, and Romanian diacritic characters.
func capitalizeSentences(text string) string {
	runes := []rune(text)
	if len(runes) == 0 {
		return text
	}
	// First non-space letter at the very start.
	for i, r := range runes {
		if unicode.IsSpace(r) {
			continue
		}
		if unicode.IsLetter(r) {
			runes[i] = unicode.ToUpper(r)
		}
		break
	}
	// Sentence boundaries: terminator + optional spaces/quotes + letter.
	terminator := false
	for i, r := range runes {
		if r == '.' || r == '!' || r == '?' || r == '…' {
			terminator = true
			continue
		}
		if !terminator {
			continue
		}
		if unicode.IsSpace(r) || r == '"' || r == '\'' || r == '«' || r == '»' || r == '(' || r == '[' {
			continue
		}
		if unicode.IsLetter(r) {
			runes[i] = unicode.ToUpper(r)
		}
		terminator = false
	}
	return string(runes)
}

func collapseSpaces(text string) string {
	text = reMultiSpaces.ReplaceAllString(text, " ")
	lines := strings.Split(text, "\n")
	for i, l := range lines {
		lines[i] = strings.TrimRight(l, " \t")
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

// softenLongSentences splits sentences longer than maxWords on the first
// coordinating-conjunction comma it can find. Conservative on purpose —
// academic prose is allowed long sentences, only monster runs get broken.
// Supports English, Russian, and Moldovan/Romanian conjunctions.
var reSentenceSplit = regexp.MustCompile(`([.!?…])\s+`)

func softenLongSentences(text string, maxWords int) (string, int) {
	parts := reSentenceSplit.Split(text, -1)
	delims := reSentenceSplit.FindAllStringSubmatch(text, -1)
	splitCount := 0
	// Multi-language conjunctions
	conjunctions := []string{
		// English
		", and ", ", but ", ", which ", ", however ", "; and ", "; but ",
		// Russian
		", и ", ", но ", ", который ", ", которая ", ", которые ", ", однако ", "; и ", "; но ",
		// Moldovan/Romanian
		", și ", ", dar ", ", care ", ", însă ", "; și ", "; dar ",
	}
	for i, sentence := range parts {
		if len(strings.Fields(sentence)) <= maxWords {
			continue
		}
		for _, conn := range conjunctions {
			lower := strings.ToLower(sentence)
			idx := strings.Index(lower, conn)
			if idx < 5 || idx > len(sentence)-10 {
				continue
			}
			left := strings.TrimRight(sentence[:idx], " ,;")
			right := strings.TrimSpace(sentence[idx+len(conn):])
			if right == "" {
				break
			}
			parts[i] = left + ". " + capitalizeFirst(right)
			splitCount++
			break
		}
	}
	var b strings.Builder
	for i, p := range parts {
		b.WriteString(p)
		if i < len(delims) {
			b.WriteString(delims[i][1])
			b.WriteString(" ")
		}
	}
	return b.String(), splitCount
}

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
