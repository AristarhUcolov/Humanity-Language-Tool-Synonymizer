package docx

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"strings"
)

// ExtractPDFText extracts readable text from a PDF file.
//
// It is a pragmatic, dependency-free extractor: it locates content streams,
// inflates the FlateDecode-compressed ones, and pulls text out of the
// text-showing operators (Tj, TJ, ', "). It handles literal strings "(…)"
// with escapes and hex strings "<…>".
//
// Limitations: it does not reconstruct exact layout, does not decode custom
// font CMaps, and cannot read encrypted PDFs or scanned-image PDFs (which
// contain no text). For well-behaved text PDFs it returns clean, usable text.
func ExtractPDFText(src []byte) (string, error) {
	if len(src) < 5 || !bytes.HasPrefix(src, []byte("%PDF-")) {
		return "", fmt.Errorf("not a PDF file")
	}

	var out strings.Builder
	for _, stream := range findStreams(src) {
		data := stream
		if inflated, ok := tryInflate(stream); ok {
			data = inflated
		}
		// Only treat it as a content stream if it contains text operators.
		if !bytes.Contains(data, []byte("BT")) &&
			!bytes.Contains(data, []byte("Tj")) &&
			!bytes.Contains(data, []byte("TJ")) {
			continue
		}
		extractTextOps(data, &out)
	}
	return cleanupPDFText(out.String()), nil
}

// findStreams returns the raw bytes between every `stream`/`endstream` pair.
func findStreams(src []byte) [][]byte {
	var streams [][]byte
	const kw = "stream"
	const end = "endstream"
	i := 0
	for {
		idx := bytes.Index(src[i:], []byte(kw))
		if idx < 0 {
			break
		}
		start := i + idx + len(kw)
		// Skip the EOL after the `stream` keyword (CRLF or LF).
		if start < len(src) && src[start] == '\r' {
			start++
		}
		if start < len(src) && src[start] == '\n' {
			start++
		}
		endIdx := bytes.Index(src[start:], []byte(end))
		if endIdx < 0 {
			break
		}
		body := src[start : start+endIdx]
		// Trim a trailing EOL before `endstream`.
		body = bytes.TrimRight(body, "\r\n")
		streams = append(streams, body)
		i = start + endIdx + len(end)
	}
	return streams
}

// tryInflate attempts zlib (FlateDecode) decompression.
func tryInflate(data []byte) ([]byte, bool) {
	zr, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, false
	}
	defer zr.Close()
	out, err := io.ReadAll(zr)
	if err != nil || len(out) == 0 {
		return nil, false
	}
	return out, true
}

// extractTextOps scans a content stream and appends the text it shows.
// It recognises:
//   (literal) Tj      — show a literal string
//   <hex> Tj          — show a hex string
//   [ (a) -10 (b) ] TJ — show an array of strings with kerning
//   (literal) '       — move to next line and show
//   T*                — move to next line
func extractTextOps(data []byte, out *strings.Builder) {
	i := 0
	n := len(data)
	for i < n {
		c := data[i]
		switch c {
		case '(':
			s, next := readLiteralString(data, i)
			out.WriteString(s)
			i = next
		case '<':
			// Hex string — but skip dictionary openers "<<".
			if i+1 < n && data[i+1] == '<' {
				i += 2
				continue
			}
			s, next := readHexString(data, i)
			out.WriteString(s)
			i = next
		case 'T':
			// T* — text line move → newline.
			if i+1 < n && data[i+1] == '*' {
				out.WriteByte('\n')
				i += 2
				continue
			}
			i++
		case '\'', '"':
			// ' and " show a string on a new line.
			out.WriteByte('\n')
			i++
		default:
			i++
		}
	}
	out.WriteByte('\n')
}

// readLiteralString reads a PDF literal string starting at data[i]=='(' and
// returns the decoded text plus the index just past the closing ')'.
func readLiteralString(data []byte, i int) (string, int) {
	var b strings.Builder
	depth := 0
	n := len(data)
	for ; i < n; i++ {
		c := data[i]
		if c == '\\' && i+1 < n {
			i++
			switch e := data[i]; e {
			case 'n':
				b.WriteByte('\n')
			case 'r':
				b.WriteByte('\r')
			case 't':
				b.WriteByte('\t')
			case 'b':
				b.WriteByte('\b')
			case 'f':
				b.WriteByte('\f')
			case '(', ')', '\\':
				b.WriteByte(e)
			case '\r':
				if i+1 < n && data[i+1] == '\n' {
					i++
				}
			case '\n':
				// line continuation — emit nothing
			default:
				if e >= '0' && e <= '7' {
					// up to 3 octal digits
					val := int(e - '0')
					for k := 0; k < 2 && i+1 < n && data[i+1] >= '0' && data[i+1] <= '7'; k++ {
						i++
						val = val*8 + int(data[i]-'0')
					}
					b.WriteByte(byte(val))
				} else {
					b.WriteByte(e)
				}
			}
			continue
		}
		switch c {
		case '(':
			if depth > 0 {
				b.WriteByte(c)
			}
			depth++
		case ')':
			depth--
			if depth == 0 {
				return decodePDFBytes(b.String()), i + 1
			}
			b.WriteByte(c)
		default:
			b.WriteByte(c)
		}
	}
	return decodePDFBytes(b.String()), i
}

// readHexString reads a PDF hex string starting at data[i]=='<'.
func readHexString(data []byte, i int) (string, int) {
	n := len(data)
	i++ // skip '<'
	var hex strings.Builder
	for ; i < n && data[i] != '>'; i++ {
		c := data[i]
		if isHexDigit(c) {
			hex.WriteByte(c)
		}
	}
	if i < n {
		i++ // skip '>'
	}
	h := hex.String()
	if len(h)%2 == 1 {
		h += "0" // an odd final digit is padded with zero per the PDF spec
	}
	var raw strings.Builder
	for k := 0; k+1 < len(h); k += 2 {
		raw.WriteByte(hexVal(h[k])<<4 | hexVal(h[k+1]))
	}
	return decodePDFBytes(raw.String()), i
}

// decodePDFBytes turns a raw PDF string into UTF-8 text. It detects a
// UTF-16BE byte-order mark; otherwise it treats the bytes as Latin-1, which
// is correct for the standard PDF text encodings used by most documents.
func decodePDFBytes(s string) string {
	b := []byte(s)
	if len(b) >= 2 && b[0] == 0xFE && b[1] == 0xFF {
		// UTF-16BE
		var out strings.Builder
		for i := 2; i+1 < len(b); i += 2 {
			out.WriteRune(rune(int(b[i])<<8 | int(b[i+1])))
		}
		return out.String()
	}
	// Latin-1 → UTF-8
	var out strings.Builder
	for _, c := range b {
		out.WriteRune(rune(c))
	}
	return out.String()
}

// cleanupPDFText collapses excess whitespace produced by the extractor while
// keeping paragraph-like line breaks.
func cleanupPDFText(s string) string {
	lines := strings.Split(s, "\n")
	var kept []string
	for _, ln := range lines {
		ln = strings.TrimRight(ln, " \t\r")
		// Collapse internal runs of spaces.
		for strings.Contains(ln, "  ") {
			ln = strings.ReplaceAll(ln, "  ", " ")
		}
		ln = strings.TrimSpace(ln)
		if ln != "" {
			kept = append(kept, ln)
		}
	}
	return strings.Join(kept, "\n")
}

func isHexDigit(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func hexVal(c byte) byte {
	switch {
	case c >= '0' && c <= '9':
		return c - '0'
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10
	}
	return 0
}
