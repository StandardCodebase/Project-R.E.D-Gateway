package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
)

type PostMetadata struct {
	Title          string `yaml:"title"`
	AuthorIdentity string `yaml:"author_identity"`
	CreatedAt      string `yaml:"created_at"`
	DiscussionHub  string `yaml:"discussion_hub"`
}

type PageData struct {
	PostMetadata
	NodeName    string
	ContentHash string
	ContentPath string // e.g., "solar-array-build/00-index" (without .md)
	HTMLContent template.HTML
}

// Configurable state
var (
	NodeName string
	DataDir  string
	Port     string

	// Pre-parsed template
	indexTemplate *template.Template
)

func init() {
	NodeName = os.Getenv("RED_NODE_NAME")
	if NodeName == "" {
		NodeName = "Alpha-Centauri-01"
	}

	DataDir = os.Getenv("RED_DATA_DIR")
	if DataDir == "" {
		DataDir = "./data"
	}

	Port = os.Getenv("RED_PORT")
	if Port == "" {
		Port = "8080"
	}
}

func main() {
	// Parse template once
	var err error
	indexTemplate, err = template.ParseFiles("templates/layout.html")
	if err != nil {
		log.Fatalf("Failed to parse template: %v", err)
	}

	// Static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","node":"` + NodeName + `"}`))
	})

	http.HandleFunc("/guides/", handleRenderGuide)
	http.HandleFunc("/download/", handleDownloadGuide)

	// Configure server with timeouts
	srv := &http.Server{
		Addr:         ":" + Port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Project R.E.D. Engine [%s] initiating on port :%s...\n", NodeName, Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Engine panic: %v", err)
	}
}

// secureJoin ensures the requested path cannot escape DataDir.
func secureJoin(baseDir, requestedPath string) (string, error) {
	baseAbs, err := filepath.Abs(baseDir)
	if err != nil {
		return "", err
	}
	targetAbs, err := filepath.Abs(filepath.Join(baseDir, requestedPath))
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(targetAbs, baseAbs+string(os.PathSeparator)) && targetAbs != baseAbs {
		return "", os.ErrPermission
	}
	return targetAbs, nil
}

// readFileWithContext reads a file while respecting request cancellation.
func readFileWithContext(ctx context.Context, path string) ([]byte, error) {
	type result struct {
		data []byte
		err  error
	}
	ch := make(chan result, 1)

	go func() {
		data, err := os.ReadFile(path)
		ch <- result{data, err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case r := <-ch:
		return r.data, r.err
	}
}

func handleRenderGuide(w http.ResponseWriter, r *http.Request) {
	httpError := func(status int) {
		http.Error(w, http.StatusText(status), status)
	}

	requestedPath := strings.TrimPrefix(r.URL.Path, "/guides/")
	if requestedPath == "" {
		httpError(http.StatusNotFound)
		return
	}

	cleanedPath := filepath.Clean(requestedPath)
	if cleanedPath == "." || strings.HasPrefix(cleanedPath, "..") || strings.Contains(cleanedPath, string(os.PathSeparator)+"..") {
		httpError(http.StatusForbidden)
		return
	}

	targetFile, err := secureJoin(DataDir, cleanedPath+".md")
	if err != nil {
		httpError(http.StatusForbidden)
		return
	}

	fileBytes, err := readFileWithContext(r.Context(), targetFile)
	if err != nil {
		if os.IsNotExist(err) {
			httpError(http.StatusNotFound)
		} else {
			httpError(http.StatusInternalServerError)
		}
		return
	}

	hasher := sha256.New()
	hasher.Write(fileBytes)
	hashString := hex.EncodeToString(hasher.Sum(nil))

	var meta PostMetadata
	markdownRaw, err := frontmatter.Parse(strings.NewReader(string(fileBytes)), &meta)
	if err != nil {
		httpError(http.StatusInternalServerError)
		return
	}

	var buf strings.Builder
	if err := goldmark.Convert(markdownRaw, &buf); err != nil {
		httpError(http.StatusInternalServerError)
		return
	}

	w.Header().Set("X-RED-Content-Hash", hashString)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	data := PageData{
		PostMetadata: meta,
		NodeName:     NodeName,
		ContentHash:  hashString,
		ContentPath:  cleanedPath, // store for download link
		HTMLContent:  template.HTML(buf.String()),
	}

	if err := indexTemplate.Execute(w, data); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

func handleDownloadGuide(w http.ResponseWriter, r *http.Request) {
	httpError := func(status int) {
		http.Error(w, http.StatusText(status), status)
	}

	requestedPath := strings.TrimPrefix(r.URL.Path, "/download/")
	if requestedPath == "" {
		httpError(http.StatusNotFound)
		return
	}

	cleanedPath := filepath.Clean(requestedPath)
	if cleanedPath == "." || strings.HasPrefix(cleanedPath, "..") || strings.Contains(cleanedPath, string(os.PathSeparator)+"..") {
		httpError(http.StatusForbidden)
		return
	}

	targetFile, err := secureJoin(DataDir, cleanedPath+".md")
	if err != nil {
		httpError(http.StatusForbidden)
		return
	}

	fileBytes, err := readFileWithContext(r.Context(), targetFile)
	if err != nil {
		if os.IsNotExist(err) {
			httpError(http.StatusNotFound)
		} else {
			httpError(http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(targetFile))
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Write(fileBytes)
}
