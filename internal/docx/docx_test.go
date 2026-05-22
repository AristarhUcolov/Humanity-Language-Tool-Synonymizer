package docx

import (
	"strings"
	"testing"
)

// upper is a trivial transform used to verify which text the pipeline touches.
func upper(s string) string { return strings.ToUpper(s) }

func TestTransformParagraphSingleRun(t *testing.T) {
	para := `<w:p><w:r><w:t>hello</w:t></w:r></w:p>`
	got := string(transformParagraph([]byte(para), upper))
	if !strings.Contains(got, "<w:t>HELLO</w:t>") {
		t.Errorf("single run not transformed: %s", got)
	}
}

func TestTransformParagraphGroupsSameFormatting(t *testing.T) {
	// Two runs, no rPr → identical formatting → grouped.
	// "give" + " up" joined becomes "GIVE UP"; first run holds the result.
	para := `<w:p><w:r><w:t>give</w:t></w:r><w:r><w:t xml:space="preserve"> up</w:t></w:r></w:p>`
	joined := ""
	got := string(transformParagraph([]byte(para), func(s string) string {
		joined = s
		return strings.ToUpper(s)
	}))
	if joined != "give up" {
		t.Errorf("runs not joined for transform: got %q, want %q", joined, "give up")
	}
	if !strings.Contains(got, "GIVE UP") {
		t.Errorf("grouped result missing: %s", got)
	}
}

func TestTransformParagraphKeepsDifferentFormattingSeparate(t *testing.T) {
	// Run 2 is bold (<w:rPr><w:b/></w:rPr>) — must NOT be grouped with run 1.
	para := `<w:p>` +
		`<w:r><w:t>plain</w:t></w:r>` +
		`<w:r><w:rPr><w:b/></w:rPr><w:t>bold</w:t></w:r>` +
		`</w:p>`
	var seen []string
	transformParagraph([]byte(para), func(s string) string {
		seen = append(seen, s)
		return s
	})
	// Two separate transform calls — never joined into one.
	for _, s := range seen {
		if s == "plainbold" {
			t.Errorf("runs with different formatting were wrongly merged")
		}
	}
	if len(seen) != 2 {
		t.Errorf("expected 2 separate transforms, got %d: %v", len(seen), seen)
	}
}

func TestTransformParagraphComplexRunNotGrouped(t *testing.T) {
	// A run containing <w:tab/> is complex and must be handled individually.
	para := `<w:p>` +
		`<w:r><w:t>left</w:t></w:r>` +
		`<w:r><w:tab/><w:t>right</w:t></w:r>` +
		`</w:p>`
	got := string(transformParagraph([]byte(para), upper))
	if !strings.Contains(got, "<w:tab/>") {
		t.Errorf("tab element lost: %s", got)
	}
	if !strings.Contains(got, "LEFT") || !strings.Contains(got, "RIGHT") {
		t.Errorf("complex-run text not transformed: %s", got)
	}
}

func TestIsSimpleRun(t *testing.T) {
	cases := []struct {
		run  string
		want bool
	}{
		{`<w:r><w:t>x</w:t></w:r>`, true},
		{`<w:r><w:rPr><w:b/></w:rPr><w:t>x</w:t></w:r>`, true},
		{`<w:r><w:tab/><w:t>x</w:t></w:r>`, false},
		{`<w:r><w:br/><w:t>x</w:t></w:r>`, false},
		{`<w:r><w:t>a</w:t><w:t>b</w:t></w:r>`, false}, // two <w:t>
	}
	for _, c := range cases {
		if got := isSimpleRun(c.run); got != c.want {
			t.Errorf("isSimpleRun(%q) = %v, want %v", c.run, got, c.want)
		}
	}
}

func TestRunPropsIdentical(t *testing.T) {
	a := `<w:r><w:rPr><w:b/></w:rPr><w:t>x</w:t></w:r>`
	b := `<w:r><w:rPr><w:b/></w:rPr><w:t>y</w:t></w:r>`
	c := `<w:r><w:rPr><w:i/></w:rPr><w:t>z</w:t></w:r>`
	if runProps(a) != runProps(b) {
		t.Errorf("identical rPr not equal")
	}
	if runProps(a) == runProps(c) {
		t.Errorf("different rPr reported equal")
	}
}

func TestSetSimpleRunTextPreservesSpace(t *testing.T) {
	run := `<w:r><w:t>x</w:t></w:r>`
	got := setSimpleRunText(run, " spaced ")
	if !strings.Contains(got, `xml:space="preserve"`) {
		t.Errorf("xml:space not added for edge spaces: %s", got)
	}
}

func TestProcessRejectsNonDocx(t *testing.T) {
	_, err := Process([]byte("not a zip"), upper)
	if err == nil {
		t.Errorf("expected error for non-DOCX input")
	}
}

func TestTransformTextNodesFallback(t *testing.T) {
	// Content without paragraphs falls back to per-<w:t> transform.
	data := []byte(`<doc><w:t>alpha</w:t> mid <w:t>beta</w:t></doc>`)
	got := string(transformContentXML(data, upper))
	if !strings.Contains(got, "ALPHA") || !strings.Contains(got, "BETA") {
		t.Errorf("fallback transform failed: %s", got)
	}
}
