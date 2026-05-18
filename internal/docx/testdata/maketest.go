//go:build ignore

// Standalone helper: generates a minimal sample.docx in the current dir.
//   go run internal/docx/testdata/maketest.go
package main

import (
	"archive/zip"
	"fmt"
	"os"
	"strings"
)

const contentTypes = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`

const packageRels = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`

const docXMLTemplate = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>
    %s
  </w:body>
</w:document>`

func paragraph(text string) string {
	return fmt.Sprintf(`<w:p><w:r><w:t xml:space="preserve">%s</w:t></w:r></w:p>`, text)
}

func main() {
	paragraphs := []string{
		"In today's fast-paced world, it is important to note that artificial intelligence plays a crucial role.",
		"Furthermore, it cannot be overstated that this technology is paramount in the digital age.",
		"Let's dive into how AI is transforming various industries.",
		"In conclusion, the ever-evolving landscape of technology continues to reshape our lives.",
	}
	var body strings.Builder
	for _, p := range paragraphs {
		body.WriteString(paragraph(p))
	}
	docXML := fmt.Sprintf(docXMLTemplate, body.String())

	f, err := os.Create("sample.docx")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	zw := zip.NewWriter(f)
	for _, entry := range []struct{ name, body string }{
		{"[Content_Types].xml", contentTypes},
		{"_rels/.rels", packageRels},
		{"word/document.xml", docXML},
	} {
		w, err := zw.Create(entry.name)
		if err != nil {
			panic(err)
		}
		w.Write([]byte(entry.body))
	}
	if err := zw.Close(); err != nil {
		panic(err)
	}
	fmt.Println("wrote sample.docx")
}
