package docx

import (
	"bytes"
	"compress/zlib"
	"strings"
	"testing"
)

func makePDF(contentStream []byte, compressed bool) []byte {
	var stream []byte
	filter := ""
	if compressed {
		var buf bytes.Buffer
		zw := zlib.NewWriter(&buf)
		zw.Write(contentStream)
		zw.Close()
		stream = buf.Bytes()
		filter = "/Filter/FlateDecode"
	} else {
		stream = contentStream
	}
	var pdf bytes.Buffer
	pdf.WriteString("%PDF-1.4\n")
	pdf.WriteString("1 0 obj<</Type/Catalog>>endobj\n")
	pdf.WriteString("4 0 obj<<" + filter + ">>stream\n")
	pdf.Write(stream)
	pdf.WriteString("\nendstream endobj\n")
	pdf.WriteString("trailer<</Root 1 0 R>>\n%%EOF")
	return pdf.Bytes()
}

func TestExtractPDFTextLiteral(t *testing.T) {
	pdf := makePDF([]byte("BT (Hello world) Tj ET"), false)
	got, err := ExtractPDFText(pdf)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !strings.Contains(got, "Hello world") {
		t.Errorf("literal text not extracted: %q", got)
	}
}

func TestExtractPDFTextCompressed(t *testing.T) {
	pdf := makePDF([]byte("BT (Compressed content) Tj ET"), true)
	got, err := ExtractPDFText(pdf)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if !strings.Contains(got, "Compressed content") {
		t.Errorf("compressed text not extracted: %q", got)
	}
}

func TestExtractPDFTextNewline(t *testing.T) {
	pdf := makePDF([]byte("BT (Line one) Tj T* (Line two) Tj ET"), false)
	got, _ := ExtractPDFText(pdf)
	if !strings.Contains(got, "Line one") || !strings.Contains(got, "Line two") {
		t.Errorf("multi-line text lost: %q", got)
	}
	if !strings.Contains(got, "\n") {
		t.Errorf("T* did not produce a line break: %q", got)
	}
}

func TestExtractPDFTextEscapes(t *testing.T) {
	// \( \) \\ escapes and an octal code.
	pdf := makePDF([]byte(`BT (a\(b\)c \\ d) Tj ET`), false)
	got, _ := ExtractPDFText(pdf)
	if !strings.Contains(got, "a(b)c") {
		t.Errorf("escapes not decoded: %q", got)
	}
}

func TestExtractPDFTextArray(t *testing.T) {
	// TJ array with kerning numbers between strings.
	pdf := makePDF([]byte("BT [(Hel) -20 (lo)] TJ ET"), false)
	got, _ := ExtractPDFText(pdf)
	if !strings.Contains(strings.ReplaceAll(got, " ", ""), "Hello") {
		t.Errorf("TJ array text not extracted: %q", got)
	}
}

func TestExtractPDFTextRejectsNonPDF(t *testing.T) {
	_, err := ExtractPDFText([]byte("just plain text"))
	if err == nil {
		t.Errorf("expected error for non-PDF input")
	}
}

func TestReadLiteralStringNested(t *testing.T) {
	// Balanced nested parentheses are part of the string.
	data := []byte("(outer (inner) text)")
	s, next := readLiteralString(data, 0)
	if s != "outer (inner) text" {
		t.Errorf("nested parens mishandled: %q", s)
	}
	if next != len(data) {
		t.Errorf("next index wrong: %d", next)
	}
}
