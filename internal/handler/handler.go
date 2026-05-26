package handler

import (
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
	Cfg      config.Config
	Template *template.Template
}

func (n *Node) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","node":"` + n.Cfg.NodeName + `"}`))
}

func (n *Node) RenderGuide(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(path))
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
