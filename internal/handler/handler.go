package handler

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"red-engine/internal/config"
	redfs "red-engine/internal/fs"
	"red-engine/internal/render"
)

type Node struct {
	Cfg           config.Config
	Template      *template.Template
	IndexTemplate *template.Template
}

func (n *Node) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"node":   n.Cfg.NodeName,
	})
}

func (n *Node) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	guides := walkGuides(n.Cfg.DataDir)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := n.IndexTemplate.Execute(w, map[string]interface{}{
		"NodeName": n.Cfg.NodeName,
		"Guides":   guides,
	}); err != nil {
		log.Printf("index template error: %v", err)
	}
}

type GuideEntry struct {
	Path  string
	Title string
}

func walkGuides(dataDir string) []GuideEntry {
	var entries []GuideEntry
	filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		rel, _ := filepath.Rel(dataDir, path)
		rel = strings.TrimSuffix(rel, ".md")
		title := rel
		raw, rerr := os.ReadFile(path)
		if rerr == nil {
			result, merr := render.Markdown(raw)
			if merr == nil && result.Meta.Title != "" {
				title = result.Meta.Title
			}
		}
		entries = append(entries, GuideEntry{Path: rel, Title: title})
		return nil
	})
	return entries
}

func (n *Node) Manifest(w http.ResponseWriter, r *http.Request) {
	guides := walkGuides(n.Cfg.DataDir)
	m := map[string]map[string]string{}
	for _, g := range guides {
		m[g.Path] = map[string]string{"title": g.Title}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

func (n *Node) RenderGuide(w http.ResponseWriter, r *http.Request) {
	// 1. Intercept static resources (images, PDFs, etc.)
	ext := strings.ToLower(filepath.Ext(r.URL.Path))
	isResource := ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" || ext == ".webp" || ext == ".svg" || ext == ".pdf"
	if isResource {
		n.serveResource(w, r)
		return
	}

	// 2. Markdown routing logic
	path, ok := n.resolveGuidePath(w, r, "/guides/")
	if !ok {
		return
	}
	raw, ok := n.readGuide(w, r, path)
	if !ok {
		return
	}
	result, err := render.Markdown(raw)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("X-RED-Content-Hash", result.Hash)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	data := config.PageData{
		PostMetadata: result.Meta,
		NodeName:     n.Cfg.NodeName,
		ContentHash:  result.Hash,
		ContentPath:  filepath.Clean(strings.TrimPrefix(r.URL.Path, "/guides/")),
		HTMLContent:  result.HTMLContent,
	}
	if err := n.Template.Execute(w, data); err != nil {
		log.Printf("template execution error: %v", err)
	}
}

// ViewSource serves the raw markdown as plain text — like raw.githubusercontent.com
func (n *Node) ViewSource(w http.ResponseWriter, r *http.Request) {
	path, ok := n.resolveGuidePath(w, r, "/source/")
	if !ok {
		return
	}
	raw, ok := n.readGuide(w, r, path)
	if !ok {
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Write(raw)
}

func (n *Node) DownloadGuide(w http.ResponseWriter, r *http.Request) {
	path, ok := n.resolveGuidePath(w, r, "/download/")
	if !ok {
		return
	}
	raw, ok := n.readGuide(w, r, path)
	if !ok {
		return
	}
	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="`+filepath.Base(path)+`"`)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Write(raw)
}

func (n *Node) resolveGuidePath(w http.ResponseWriter, r *http.Request, prefix string) (string, bool) {
	requested := strings.TrimPrefix(r.URL.Path, prefix)
	if requested == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return "", false
	}
	cleaned := filepath.Clean(requested)
	if cleaned == "." || strings.HasPrefix(cleaned, "..") || strings.Contains(cleaned, string(os.PathSeparator)+"..") {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return "", false
	}
	resolved, err := redfs.SecureJoin(n.Cfg.DataDir, cleaned+".md")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return "", false
	}
	return resolved, true
}

func (n *Node) readGuide(w http.ResponseWriter, r *http.Request, path string) ([]byte, bool) {
	raw, err := redfs.ReadFileWithContext(r.Context(), path)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return nil, false
	}
	return raw, true
}

func (n *Node) serveResource(w http.ResponseWriter, r *http.Request) {
	requested := strings.TrimPrefix(r.URL.Path, "/guides/")
	cleaned := filepath.Clean(requested)

	// Path traversal protection
	if cleaned == "." || strings.HasPrefix(cleaned, "..") || strings.Contains(cleaned, string(os.PathSeparator)+"..") {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	//NO ".md" is appended here
	resolved, err := redfs.SecureJoin(n.Cfg.DataDir, cleaned)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	// Optional but highly recommended: Cache images on the reader's browser
	// for 24 hours to drastically reduce bandwidth costs for Node Operators.
	w.Header().Set("Cache-Control", "public, max-age=86400")
	http.ServeFile(w, r, resolved)
}
