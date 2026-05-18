package humanize

import (
	"regexp"
	"sort"
	"strings"
	"sync"
	"unicode"
)

// replaceMap performs case-insensitive whole-phrase replacements while
// preserving the leading-letter case of each match. Longer phrases win
// over shorter overlapping ones to avoid "important" eating "very important".
func replaceMap(text string, m map[string]string) (string, int) {
	if len(m) == 0 || text == "" {
		return text, 0
	}
	// Sort keys by length desc so longer phrases replace first.
	keys := make([]string, 0, len(m))
	for k := range m {
		if k != "" {
			keys = append(keys, k)
		}
	}
	sort.Slice(keys, func(i, j int) bool { return len(keys[i]) > len(keys[j]) })

	total := 0
	for _, k := range keys {
		re := patternFor(k)
		text = re.ReplaceAllStringFunc(text, func(match string) string {
			// match contains: boundary-prefix (0 or 1 non-letter) + phrase + boundary-suffix (0 or 1 non-letter).
			// We need to extract just the phrase portion to preserve casing, then re-attach the boundaries.
			prefix, body, suffix := splitBoundary(match)
			repl := m[k]
			repl = preserveCase(body, repl)
			total++
			if repl == "" {
				// Deleting: keep one separator if both sides had whitespace, else collapse.
				if isWhitespaceOnly(prefix) && isWhitespaceOnly(suffix) {
					return " "
				}
				// If we ate a comma + space, collapse to nothing (cleanupAfterDelete handles orphans).
				return prefix + suffix
			}
			return prefix + repl + suffix
		})
	}
	return text, total
}

// patternFor returns a cached regexp that matches the phrase with single-char
// non-letter boundaries (or string edges).
func patternFor(phrase string) *regexp.Regexp {
	if cached, ok := loadPattern(phrase); ok {
		return cached
	}
	parts := strings.Fields(phrase)
	for i, p := range parts {
		parts[i] = regexp.QuoteMeta(p)
	}
	body := strings.Join(parts, `\s+`)
	pattern := `(?i)(^|[^\p{L}])(` + body + `)([^\p{L}]|$)`
	re := regexp.MustCompile(pattern)
	storePattern(phrase, re)
	return re
}

var (
	patternMu    sync.RWMutex
	patternCache = map[string]*regexp.Regexp{}
)

func loadPattern(k string) (*regexp.Regexp, bool) {
	patternMu.RLock()
	defer patternMu.RUnlock()
	re, ok := patternCache[k]
	return re, ok
}

func storePattern(k string, re *regexp.Regexp) {
	patternMu.Lock()
	defer patternMu.Unlock()
	patternCache[k] = re
}

// splitBoundary peels at most ONE leading and ONE trailing non-letter char off
// of the match returned by ReplaceAllStringFunc — those correspond exactly to
// the boundary capture groups in our pattern.
func splitBoundary(match string) (prefix, body, suffix string) {
	runes := []rune(match)
	n := len(runes)
	if n == 0 {
		return "", "", ""
	}
	start, end := 0, n
	if !unicode.IsLetter(runes[0]) {
		prefix = string(runes[0])
		start = 1
	}
	if end > start && !unicode.IsLetter(runes[end-1]) {
		suffix = string(runes[end-1])
		end--
	}
	body = string(runes[start:end])
	return
}

func preserveCase(original, replacement string) string {
	if original == "" || replacement == "" {
		return replacement
	}
	firstO := []rune(original)[0]
	if !unicode.IsUpper(firstO) {
		return replacement
	}
	repRunes := []rune(replacement)
	repRunes[0] = unicode.ToUpper(repRunes[0])
	return string(repRunes)
}

func isWhitespaceOnly(s string) bool {
	if s == "" {
		return true
	}
	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}
