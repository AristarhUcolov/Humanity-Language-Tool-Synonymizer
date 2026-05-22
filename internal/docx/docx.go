// Package docx processes Office Open XML (.docx) files by applying a text
// transformation to the document's text while preserving all formatting,
// styles, and structure.
//
// CROSS-RUN MATCHING:
// A DOCX paragraph is split into "runs" (<w:r>), and a single word can land
// in its own run when formatting changes. To let multi-word synonyms match
// across run boundaries, this package groups adjacent runs that share
// identical formatting (<w:rPr>), transforms the joined text, and writes the
// result back — putting the new text in the group's first run and emptying
// the rest. Groups never mix different formatting, so bold/italic/colour are
// preserved exactly.
//
// Runs that are "complex" (contain tabs, breaks, fields, drawings, or more
// than one <w:t>) are never grouped — they are transformed individually so
// nothing structural is disturbed.
//
// FORMATTING PRESERVED: fonts (Times New Roman, Arial, …), sizes, bold,
// italic, underline, strikethrough, colours, highlight, super/subscript,
// effects, paragraph/character styles.
//
// STRUCTURE PRESERVED: paragraphs, tables, lists, headers/footers, sections,
// hyperlinks, fields, comments, tracked changes, footnotes/endnotes, text
// boxes, equations, images, embedded objects.
//
// CONTENT FILES PROCESSED: word/document.xml, word/header*.xml,
// word/footer*.xml, word/footnotes.xml, word/endnotes.xml, word/comments*.xml.
package docx

import (
	"archive/zip"
	"bytes"
	"fmt"
	"html"
	"io"
	"regexp"
	"strings"
)

// Process reads a DOCX from src and writes a new DOCX to the returned byte
// slice with the document text transformed by `transform`.
//
// transform receives a piece of text (a run, or the joined text of a group
// of identically-formatted runs) and returns the rewritten content. All
// formatting, styles, images and layout are preserved.
func Process(src []byte, transform func(string) string) ([]byte, error) {
	if len(src) < 4 || string(src[:4]) != "PK\x03\x04" {
		return nil, fmt.Errorf("not a DOCX (missing ZIP signature)")
	}
	zr, err := zip.NewReader(bytes.NewReader(src), int64(len(src)))
	if err != nil {
		return nil, fmt.Errorf("open zip: %w", err)
	}

	var out bytes.Buffer
	zw := zip.NewWriter(&out)

	for _, f := range zr.File {
		if err := copyEntry(zw, f, transform); err != nil {
			return nil, err
		}
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func copyEntry(zw *zip.Writer, f *zip.File, transform func(string) string) error {
	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("open %s: %w", f.Name, err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return fmt.Errorf("read %s: %w", f.Name, err)
	}

	if shouldProcess(f.Name) {
		data = transformContentXML(data, transform)
	}

	header := &zip.FileHeader{
		Name:     f.Name,
		Method:   f.Method,
		Modified: f.Modified,
	}
	w, err := zw.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("create %s: %w", f.Name, err)
	}
	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("write %s: %w", f.Name, err)
	}
	return nil
}

// shouldProcess returns true for every part of the DOCX that contains
// user-authored text. Skips themes, styles, settings, fonts, templates.
func shouldProcess(name string) bool {
	switch {
	case name == "word/document.xml":
		return true
	case strings.HasPrefix(name, "word/header") && strings.HasSuffix(name, ".xml"):
		return true
	case strings.HasPrefix(name, "word/footer") && strings.HasSuffix(name, ".xml"):
		return true
	case name == "word/footnotes.xml":
		return true
	case name == "word/endnotes.xml":
		return true
	case strings.HasPrefix(name, "word/comments") && strings.HasSuffix(name, ".xml"):
		return true
	case strings.HasPrefix(name, "word/") && strings.Contains(name, "textbox"):
		return true
	}
	return false
}

var (
	// reParagraph matches a <w:p>…</w:p> element. \b after "w:p" stops it from
	// matching <w:pPr>. Paragraphs do not nest, so a non-greedy body is safe.
	reParagraph = regexp.MustCompile(`(?s)<w:p\b[^>]*?>.*?</w:p>`)

	// reRun matches a <w:r>…</w:r> element. \b after "w:r" stops it from
	// matching <w:rPr>/<w:rStyle>. Runs do not nest.
	reRun = regexp.MustCompile(`(?s)<w:r\b[^>]*?>.*?</w:r>`)

	// reRunPr captures the optional <w:rPr>…</w:rPr> properties block.
	reRunPr = regexp.MustCompile(`(?s)<w:rPr>.*?</w:rPr>`)

	// reTextNode matches a <w:t> element: tag, attributes, inner text.
	reTextNode = regexp.MustCompile(`(?s)<(w:t|t)((?:\s[^>]*)?)>(.*?)</(?:w:t|t)>`)

	// complexRunMarkers indicate a run that must NOT be grouped: it carries
	// non-text content whose position matters.
	complexRunMarkers = []string{
		"<w:tab", "<w:br", "<w:cr", "<w:fldChar", "<w:instrText",
		"<w:drawing", "<w:object", "<w:pict", "<w:sym", "<w:noBreakHyphen",
		"<w:softHyphen", "<mc:AlternateContent",
	}
)

// transformContentXML applies transform across the XML. It works paragraph by
// paragraph so multi-word synonyms can match across run boundaries. If the
// part has no paragraphs, it falls back to a plain per-<w:t> transform.
func transformContentXML(data []byte, transform func(string) string) []byte {
	if !bytes.Contains(data, []byte("<w:p")) {
		return transformTextNodes(data, transform)
	}
	return reParagraph.ReplaceAllFunc(data, func(para []byte) []byte {
		return transformParagraph(para, transform)
	})
}

// transformParagraph processes one <w:p> element. It splits the paragraph into
// run and non-run segments, groups adjacent simple runs that share identical
// formatting, and transforms each group's joined text as a unit.
func transformParagraph(para []byte, transform func(string) string) []byte {
	s := string(para)
	runLocs := reRun.FindAllStringIndex(s, -1)
	if len(runLocs) == 0 {
		// No runs — still transform any stray <w:t> nodes.
		return transformTextNodes(para, transform)
	}

	var b strings.Builder
	cursor := 0

	// A pending group of consecutive simple runs with identical rPr.
	type runRef struct{ start, end int }
	var group []runRef
	var groupPr string

	flush := func() {
		if len(group) == 0 {
			return
		}
		if len(group) == 1 {
			// Single run — transform it on its own.
			r := group[0]
			b.WriteString(transformSimpleRun(s[r.start:r.end], transform))
		} else {
			// Join all runs' text, transform together, redistribute.
			var joined strings.Builder
			for _, r := range group {
				joined.WriteString(simpleRunText(s[r.start:r.end]))
			}
			orig := joined.String()
			out := transform(orig)
			if out == orig {
				// No change — emit runs verbatim.
				for _, r := range group {
					b.WriteString(s[r.start:r.end])
				}
			} else {
				// First run gets the whole result; the rest are emptied.
				b.WriteString(setSimpleRunText(s[group[0].start:group[0].end], out))
				for _, r := range group[1:] {
					b.WriteString(setSimpleRunText(s[r.start:r.end], ""))
				}
			}
		}
		group = group[:0]
		groupPr = ""
	}

	for _, loc := range runLocs {
		// Emit the literal segment before this run.
		gap := s[cursor:loc[0]]
		run := s[loc[0]:loc[1]]
		cursor = loc[1]

		// A non-whitespace gap (bookmark, proofErr, …) breaks any group.
		if strings.TrimSpace(gap) != "" {
			flush()
		}

		if isSimpleRun(run) {
			pr := runProps(run)
			if len(group) > 0 && pr != groupPr {
				flush()
			}
			if len(group) == 0 {
				groupPr = pr
			}
			// Defer writing the run; record its location.
			// The gap must be written in order: flush() above already ran if
			// needed, so write the gap now, then stage the run.
			b.WriteString(gap)
			group = append(group, runRef{loc[0], loc[1]})
		} else {
			// Complex run — flush the group, then transform it standalone.
			flush()
			b.WriteString(gap)
			b.Write(transformTextNodes([]byte(run), transform))
		}
	}
	flush()
	b.WriteString(s[cursor:]) // trailing literal (</w:p> and friends)
	return []byte(b.String())
}

// isSimpleRun reports whether a run is safe to group: it has exactly one <w:t>
// and carries no position-sensitive content (tabs, breaks, fields, drawings).
func isSimpleRun(run string) bool {
	for _, m := range complexRunMarkers {
		if strings.Contains(run, m) {
			return false
		}
	}
	return len(reTextNode.FindAllStringIndex(run, -1)) == 1
}

// runProps returns the run's <w:rPr>…</w:rPr> block, or "" if it has none.
// Two runs with byte-identical props are treated as identically formatted.
func runProps(run string) string {
	return reRunPr.FindString(run)
}

// simpleRunText returns the decoded text of a simple run's single <w:t>.
func simpleRunText(run string) string {
	m := reTextNode.FindStringSubmatch(run)
	if len(m) < 4 {
		return ""
	}
	return html.UnescapeString(m[3])
}

// setSimpleRunText replaces the text of a simple run's single <w:t> with
// newText, adding xml:space="preserve" when edge whitespace must survive.
func setSimpleRunText(run, newText string) string {
	return reTextNode.ReplaceAllStringFunc(run, func(node string) string {
		m := reTextNode.FindStringSubmatch(node)
		if len(m) < 4 {
			return node
		}
		tag, attrs := m[1], m[2]
		if !strings.Contains(attrs, "xml:space") && hasEdgeSpace(newText) {
			attrs += ` xml:space="preserve"`
		}
		return "<" + tag + attrs + ">" + html.EscapeString(newText) + "</" + tag + ">"
	})
}

// transformSimpleRun transforms a lone simple run (not part of a group).
func transformSimpleRun(run string, transform func(string) string) string {
	orig := simpleRunText(run)
	if orig == "" {
		return run
	}
	out := transform(orig)
	if out == orig {
		return run
	}
	return setSimpleRunText(run, out)
}

// transformTextNodes applies transform to every <w:t> node independently.
// Used as a fallback for complex runs and parts without paragraphs.
func transformTextNodes(data []byte, transform func(string) string) []byte {
	return reTextNode.ReplaceAllFunc(data, func(node []byte) []byte {
		sub := reTextNode.FindSubmatch(node)
		if len(sub) < 4 {
			return node
		}
		tag := string(sub[1])
		attrs := string(sub[2])
		raw := string(sub[3])
		if raw == "" {
			return node
		}
		text := html.UnescapeString(raw)
		out := transform(text)
		if out == text {
			return node
		}
		if !strings.Contains(attrs, "xml:space") &&
			(hasEdgeSpace(out) || hasEdgeSpace(text)) {
			attrs += ` xml:space="preserve"`
		}
		var b strings.Builder
		b.WriteByte('<')
		b.WriteString(tag)
		b.WriteString(attrs)
		b.WriteByte('>')
		b.WriteString(html.EscapeString(out))
		b.WriteString("</")
		b.WriteString(tag)
		b.WriteByte('>')
		return []byte(b.String())
	})
}

func hasEdgeSpace(s string) bool {
	if s == "" {
		return false
	}
	return s[0] == ' ' || s[0] == '\t' || s[len(s)-1] == ' ' || s[len(s)-1] == '\t'
}

// ExtractPlainText concatenates the text content of word/document.xml,
// up to maxBytes, so callers can run language detection on it.
// Returns an empty string if the DOCX cannot be opened.
func ExtractPlainText(src []byte, maxBytes int) string {
	zr, err := zip.NewReader(bytes.NewReader(src), int64(len(src)))
	if err != nil {
		return ""
	}
	for _, f := range zr.File {
		if f.Name != "word/document.xml" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return ""
		}
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return ""
		}
		var b strings.Builder
		matches := reTextNode.FindAllSubmatch(data, -1)
		for _, m := range matches {
			if len(m) < 4 {
				continue
			}
			b.WriteString(html.UnescapeString(string(m[3])))
			b.WriteByte(' ')
			if maxBytes > 0 && b.Len() >= maxBytes {
				break
			}
		}
		return b.String()
	}
	return ""
}
