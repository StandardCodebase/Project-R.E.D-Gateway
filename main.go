package main

import (
	"crypto/sha256"
	"encoding/hex"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
)

// Structural mapping for the incoming YAML block
type PostMetadata struct {
	Title          string `yaml:"title"`
	AuthorIdentity string `yaml:"author_identity"`
	CreatedAt      string `yaml:"created_at"`
	DiscussionHub  string `yaml:"discussion_hub"`
}

// Data payload injected directly into layout.html
type PageData struct {
	PostMetadata
	NodeName    string
	ContentHash string
	HTMLContent template.HTML
}

const NodeName = "Alpha-Centauri-01" // Configure your localized node name here
const DataDir = "./data"

func main() {
	// 1. Serve static files (The CSS layout)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 2. Main content routing loop
	http.HandleFunc("/guides/", handleRenderGuide)

	log.Println("Project R.E.D. Engine initiating on port :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Engine panic: %v", err)
	}
}

func handleRenderGuide(w http.ResponseWriter, r *http.Request) {
	// A. Path Bounds Verification (Sanitize input against path-traversal attacks)
	requestedPath := strings.TrimPrefix(r.URL.Path, "/guides/")
	cleanedPath := filepath.Clean(requestedPath)
	if strings.HasPrefix(cleanedPath, "..") || strings.HasPrefix(cleanedPath, "/") {
		http.Error(w, "Access Denied: Bounds violation.", http.StatusForbidden)
		return
	}

	targetFile := filepath.Join(DataDir, cleanedPath+".md")

	// B. Stream Raw File bytes from Disk
	fileBytes, err := os.ReadFile(targetFile)
	if err != nil {
		http.Error(w, "Document not found inside local state volume.", http.StatusNotFound)
		return
	}

	// C. Extract Payload Fingerprint (Generate SHA-256 Checksum over raw bytes)
	hasher := sha256.New()
	hasher.Write(fileBytes)
	hashString := hex.EncodeToString(hasher.Sum(nil))

	// D. Extract YAML Front-Matter Metadata
	var meta PostMetadata
	markdownRaw, err := frontmatter.Parse(strings.NewReader(string(fileBytes)), &meta)
	if err != nil {
		http.Error(w, "Malformed front-matter blueprint.", http.StatusInternalServerError)
		return
	}

	// E. Execute Goldmark Engine (Convert Markdown block to native HTML)
	var buf strings.Builder
	if err := goldmark.Convert(markdownRaw, &buf); err != nil {
		http.Error(w, "Goldmark compilation loop exception.", http.StatusInternalServerError)
		return
	}

	// F. Inject Cryptographic Validation Headers
	w.Header().Set("X-RED-Content-Hash", hashString)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// G. Populate HTML Layout Template
	tmpl, err := template.ParseFiles("templates/layout.html")
	if err != nil {
		http.Error(w, "Template resolution compilation fault.", http.StatusInternalServerError)
		return
	}

	data := PageData{
		PostMetadata: meta,
		NodeName:     NodeName,
		ContentHash:  hashString,
		HTMLContent:  template.HTML(buf.String()),
	}

	tmpl.Execute(w, data)
}
