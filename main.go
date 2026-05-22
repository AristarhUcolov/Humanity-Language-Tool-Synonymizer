package main

import (
	"archive/zip"
	"bytes"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"humanity/internal/docx"
	"humanity/internal/humanize"
)

//go:embed web/*
var webFS embed.FS

type humanizeRequest struct {
	Text     string `json:"text"`
	Language string `json:"language,omitempty"` // "auto", "en", "ru", "mo"
}

type analysisDTO struct {
	Chars     int     `json:"chars"`
	Words     int     `json:"words"`
	Sentences int     `json:"sentences"`
	AvgWord   float64 `json:"avgWord"`
	AvgSent   float64 `json:"avgSent"`
	AIScore   int     `json:"aiScore"`
}

type synonymizeResponse struct {
	Output           string      `json:"output"`
	Before           analysisDTO `json:"before"`
	After            analysisDTO `json:"after"`
	DetectedLanguage string      `json:"detectedLanguage"`
	SynonymsApplied  int         `json:"synonymsApplied"`
	PunctFixes       int         `json:"punctFixes"`
	SentencesSplit   int         `json:"sentencesSplit"`
	Notes            []string    `json:"notes"`
}

type languagesResponse struct {
	Languages []humanize.LanguageInfo `json:"languages"`
}

func main() {
	addr := flag.String("addr", "0.0.0.0:8080", "address to listen on (0.0.0.0:8080 for network access, 127.0.0.1:8080 for localhost only)")
	noOpen := flag.Bool("no-open", false, "do not auto-open the browser")
	flag.Parse()

	sub, err := fs.Sub(webFS, "web")
	if err != nil {
		log.Fatalf("embed: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(sub)))
	mux.HandleFunc("/api/humanize", handleHumanize)
	mux.HandleFunc("/api/humanize-docx", handleHumanizeDocx)
	mux.HandleFunc("/api/humanize-docx-batch", handleHumanizeDocxBatch)
	mux.HandleFunc("/api/humanize-pdf", handleHumanizePDF)
	mux.HandleFunc("/api/languages", handleLanguages)
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("listen %s: %v", *addr, err)
	}

	// browserURL is the address the browser actually opens. A wildcard bind
	// (0.0.0.0 or [::]) is not a routable address — browsers cannot connect to
	// it — so we always open via loopback 127.0.0.1 on the actual port.
	port := "8080"
	if _, p, splitErr := net.SplitHostPort(ln.Addr().String()); splitErr == nil && p != "" {
		port = p
	}
	browserURL := "http://127.0.0.1:" + port

	// If listening on 0.0.0.0, also show localhost and IP addresses
	if strings.Contains(*addr, "0.0.0.0") {
		fmt.Println("\n━━━ Synonymaizer Server ━━━")
		fmt.Println("🌐 Network Access (from other devices):")
		getLocalIPs(port)
		fmt.Println("Local only (localhost):   " + browserURL)
		fmt.Println("\nStop the server: Ctrl+C")
	} else {
		fmt.Println("Synonymaizer Server")
		fmt.Println("Open in browser:", browserURL)
		fmt.Println("Stop the server: Ctrl+C")
	}

	if !*noOpen {
		go func() {
			time.Sleep(400 * time.Millisecond)
			openBrowser(browserURL)
		}()
	}

	srv := &http.Server{
		Handler:      mux,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func toDTO(a humanize.Analysis) analysisDTO {
	return analysisDTO{
		Chars:     a.Chars,
		Words:     a.Words,
		Sentences: a.Sentences,
		AvgWord:   roundTo(a.AvgWordLen, 2),
		AvgSent:   roundTo(a.AvgSentLen, 2),
		AIScore:   a.AIScore,
	}
}

func roundTo(v float64, digits int) float64 {
	factor := 1.0
	for i := 0; i < digits; i++ {
		factor *= 10
	}
	return float64(int(v*factor+0.5)) / factor
}

func handleLanguages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(languagesResponse{
		Languages: humanize.AllLanguages(),
	})
}

func handleHumanize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	var req humanizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json: "+err.Error(), http.StatusBadRequest)
		return
	}
	if len(req.Text) > 500_000 {
		http.Error(w, "text too long (max 500k chars)", http.StatusRequestEntityTooLarge)
		return
	}

	res := humanize.Process(humanize.Input{
		Text:     req.Text,
		Language: humanize.Language(req.Language),
	})

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(synonymizeResponse{
		Output:           res.Output,
		Before:           toDTO(res.Before),
		After:            toDTO(res.After),
		DetectedLanguage: string(res.DetectedLanguage),
		SynonymsApplied:  res.SynonymsApplied,
		PunctFixes:       res.PunctFixes,
		SentencesSplit:   res.SentencesSplit,
		Notes:            res.Notes,
	})
}

func handleHumanizeDocx(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		http.Error(w, "bad multipart: "+err.Error(), http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "no file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	if !strings.HasSuffix(strings.ToLower(header.Filename), ".docx") {
		http.Error(w, "only .docx is supported", http.StatusUnsupportedMediaType)
		return
	}

	// Get language from form (optional, defaults to auto-detection)
	lang := humanize.Language(r.FormValue("language"))

	srcBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "read file: "+err.Error(), http.StatusBadRequest)
		return
	}

	// If auto-detect, use the document's overall text to detect language once
	// (instead of per-run, which would be noisy and inconsistent).
	if lang == "" || lang == humanize.LangAuto {
		fullText := docx.ExtractPlainText(srcBytes, 50_000)
		lang = humanize.DetectLanguage(fullText)
	}

	totalRuns := 0
	changedRuns := 0
	totalSynonyms := 0

	// Use ProcessRun (not Process) for DOCX: each w:t node is a fragment, so
	// only synonyms are swapped — whitespace, capitalisation, punctuation, and
	// paragraph structure are left exactly as in the original document.
	out, err := docx.Process(srcBytes, func(s string) string {
		totalRuns++
		if strings.TrimSpace(s) == "" {
			return s
		}
		res := humanize.ProcessRun(humanize.Input{Text: s, Language: lang})
		totalSynonyms += res.SynonymsApplied
		if res.Output != s {
			changedRuns++
		}
		return res.Output
	})
	if err != nil {
		http.Error(w, "docx process: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Compute before/after word counts across the whole doc.
	before := humanize.Analyse(docx.ExtractPlainText(srcBytes, 200_000))
	after := humanize.Analyse(docx.ExtractPlainText(out, 200_000))

	outName := strings.TrimSuffix(header.Filename, ".docx") + ".synonymized.docx"
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, outName))
	w.Header().Set("X-Humanity-Language", string(lang))
	w.Header().Set("X-Humanity-Runs-Total", fmt.Sprintf("%d", totalRuns))
	w.Header().Set("X-Humanity-Runs-Changed", fmt.Sprintf("%d", changedRuns))
	w.Header().Set("X-Humanity-Synonyms", fmt.Sprintf("%d", totalSynonyms))
	w.Header().Set("X-Humanity-Words-Before", fmt.Sprintf("%d", before.Words))
	w.Header().Set("X-Humanity-Words-After", fmt.Sprintf("%d", after.Words))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(out)))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(out)
}

// synonymizeDocxBytes runs the DOCX synonymizer on one file's bytes and
// returns the processed document. Language is detected per file when lang
// is empty or "auto".
func synonymizeDocxBytes(src []byte, lang humanize.Language) ([]byte, error) {
	useLang := lang
	if useLang == "" || useLang == humanize.LangAuto {
		useLang = humanize.DetectLanguage(docx.ExtractPlainText(src, 50_000))
	}
	return docx.Process(src, func(s string) string {
		if strings.TrimSpace(s) == "" {
			return s
		}
		return humanize.ProcessRun(humanize.Input{Text: s, Language: useLang}).Output
	})
}

// handleHumanizeDocxBatch processes several .docx files in one request and
// returns them packed into a single ZIP archive.
func handleHumanizeDocxBatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseMultipartForm(200 << 20); err != nil {
		http.Error(w, "bad multipart: "+err.Error(), http.StatusBadRequest)
		return
	}
	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		http.Error(w, "no files", http.StatusBadRequest)
		return
	}
	lang := humanize.Language(r.FormValue("language"))

	var zipBuf bytes.Buffer
	zw := zip.NewWriter(&zipBuf)
	processed, failed := 0, 0

	for _, fh := range files {
		name := fh.Filename
		if !strings.HasSuffix(strings.ToLower(name), ".docx") {
			failed++
			continue
		}
		f, err := fh.Open()
		if err != nil {
			failed++
			continue
		}
		src, err := io.ReadAll(f)
		f.Close()
		if err != nil {
			failed++
			continue
		}
		out, err := synonymizeDocxBytes(src, lang)
		if err != nil {
			failed++
			continue
		}
		outName := strings.TrimSuffix(name, ".docx") + ".synonymized.docx"
		entry, err := zw.Create(outName)
		if err != nil {
			failed++
			continue
		}
		if _, err := entry.Write(out); err != nil {
			failed++
			continue
		}
		processed++
	}
	if err := zw.Close(); err != nil {
		http.Error(w, "zip: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if processed == 0 {
		http.Error(w, "no .docx files could be processed", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="synonymized_docx.zip"`)
	w.Header().Set("X-Humanity-Processed", fmt.Sprintf("%d", processed))
	w.Header().Set("X-Humanity-Failed", fmt.Sprintf("%d", failed))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", zipBuf.Len()))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(zipBuf.Bytes())
}

// handleHumanizePDF extracts the text from a PDF, synonymizes it, and returns
// the result as a plain-text file. PDF layout/formatting is not reconstructed
// — only the text is processed (unlike DOCX, where formatting is preserved).
func handleHumanizePDF(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		http.Error(w, "bad multipart: "+err.Error(), http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "no file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	if !strings.HasSuffix(strings.ToLower(header.Filename), ".pdf") {
		http.Error(w, "only .pdf is supported", http.StatusUnsupportedMediaType)
		return
	}
	src, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "read file: "+err.Error(), http.StatusBadRequest)
		return
	}

	text, err := docx.ExtractPDFText(src)
	if err != nil {
		http.Error(w, "pdf extract: "+err.Error(), http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(text) == "" {
		http.Error(w, "no extractable text found in PDF (it may be scanned images or encrypted)", http.StatusBadRequest)
		return
	}

	lang := humanize.Language(r.FormValue("language"))
	res := humanize.Process(humanize.Input{Text: text, Language: lang})

	outName := strings.TrimSuffix(header.Filename, ".pdf") + ".synonymized.txt"
	out := []byte(res.Output)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, outName))
	w.Header().Set("X-Humanity-Language", string(res.DetectedLanguage))
	w.Header().Set("X-Humanity-Synonyms", fmt.Sprintf("%d", res.SynonymsApplied))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(out)))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(out)
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	_ = cmd.Start()
}

func getLocalIPs(port string) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	shown := make(map[string]bool)
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			ip := ipnet.IP.String()
			if !shown[ip] && ipnet.IP.To4() != nil {
				fmt.Printf("   http://%s:%s\n", ip, port)
				shown[ip] = true
			}
		}
	}
}
