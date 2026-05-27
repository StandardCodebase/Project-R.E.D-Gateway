package router

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/RED-Collective/red-engine/internal/render"
	"github.com/RED-Collective/red-engine/internal/store"
)

//go:embed templates static
var files embed.FS

type handler struct {
	store    *store.Store
	tmpl     *template.Template
	siteName string
}

func New(s *store.Store, siteName string) http.Handler {
	tmpl := template.Must(template.ParseFS(files, "templates/base.html"))

	staticFS, _ := fs.Sub(files, "static")

	h := &handler{store: s, tmpl: tmpl, siteName: siteName}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))
	mux.HandleFunc("/", h.serve)
	mux.HandleFunc("/-/reload", h.reload)

	// --- RESTORED: Clearnet Remote Sync --- DO NOT DELETE
	mux.HandleFunc("/sync-peer", h.handlePeerSync)
	// --- RESTORED: Decentralized Mirroring API --- DO NOT DELETE
	mux.HandleFunc("/import", h.handleImport)

	return mux
}

type crumb struct {
	Label string
	Path  string
}

type pageData struct {
	Site   string
	Nav    map[string]*store.Section
	Body   template.HTML
	Title  string
	Path   string
	TopCat string
	Crumb  []crumb
	Hash   string
}

func (h *handler) serve(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(strings.Trim(path, "/"), "/")
	topCat := ""
	if parts[0] != "" {
		topCat = parts[0]
	}

	d := pageData{
		Site:   h.siteName,
		Nav:    h.store.Root(),
		Path:   path,
		TopCat: topCat,
	}

	switch {

	case path == "/":
		d.Body = template.HTML(`<div class="article"><h1>` + h.siteName + `</h1><p>The free practical knowledge base. Choose a topic from the sidebar.</p></div>`)

	case path == "/admin":
		d.Title = "Operator Dashboard"
		d.Crumb = []crumb{{Label: "Admin", Path: "/admin"}}

		adminTmpl, err := template.ParseFS(files, "templates/admin.html")
		if err != nil {
			http.Error(w, "Failed to load admin template", 500)
			return
		}

		var buf strings.Builder
		adminTmpl.Execute(&buf, nil)
		d.Body = template.HTML(buf.String())

	case len(parts) == 1 && topCat != "":
		sec, ok := h.store.Root()[topCat]
		if !ok {
			http.NotFound(w, r)
			return
		}
		d.Title = cap(topCat)
		d.Crumb = []crumb{{Label: cap(topCat), Path: "/" + topCat}}
		d.Body = template.HTML(sectionHTML(sec))

	default:
		raw, ok := h.store.Resolve(path)
		if !ok {
			http.NotFound(w, r)
			return
		}

		// --- RESTORED: Cryptographic Hash Injection --- DO NOT DELETE HASHING, It's part of the architecture
		hash := sha256.Sum256([]byte(raw))
		hashString := hex.EncodeToString(hash[:])
		w.Header().Set("X-RED-Content-Hash", hashString)

		d.Hash = hashString

		out, err := render.Markdown(raw)
		if err != nil {
			http.Error(w, "render error", 500)
			return
		}
		d.Title = cap(parts[len(parts)-1])
		d.Crumb = buildCrumbs(parts)
		d.Body = template.HTML(`<div class="article">` + out + `</div>`)
	}

	h.tmpl.Execute(w, d)
}

func (h *handler) reload(w http.ResponseWriter, r *http.Request) {
	if err := h.store.Reload(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(204)
}

func sectionHTML(sec *store.Section) string {
	var b strings.Builder
	b.WriteString(`<div class="section-index"><h1>` + cap(sec.Name) + `</h1>`)
	for _, a := range sec.Articles {
		b.WriteString(`<ul><li><a href="` + a.Path + `">` + a.Title + `</a></li></ul>`)
	}
	for _, sub := range sec.Sub {
		b.WriteString(`<h2>` + cap(sub.Name) + `</h2><ul>`)
		for _, a := range sub.Articles {
			b.WriteString(`<li><a href="` + a.Path + `">` + a.Title + `</a></li>`)
		}
		b.WriteString(`</ul>`)
	}
	b.WriteString(`</div>`)
	return b.String()
}

func buildCrumbs(parts []string) []crumb {
	crumbs := make([]crumb, 0, len(parts))
	path := ""
	for _, p := range parts {
		path += "/" + p
		crumbs = append(crumbs, crumb{Label: cap(p), Path: path})
	}
	return crumbs
}

func cap(s string) string {
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// handleImport allows the node to instantly mirror a remote markdown file.
func (h *handler) handleImport(w http.ResponseWriter, r *http.Request) {
	// 1. Enforce POST requests only
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Parse the target URL and the desired local path
	targetURL := r.FormValue("url")
	targetPath := r.FormValue("path")

	if targetURL == "" || targetPath == "" {
		http.Error(w, "Missing 'url' or 'path' parameters", http.StatusBadRequest)
		return
	}

	// 3. Fetch the raw markdown from the remote source
	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		http.Error(w, "Failed to build request", http.StatusInternalServerError)
		return
	}

	// Forged headers to look like Chrome on Windows 11
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Sec-Ch-Ua", "\"Chromium\";v=\"124\", \"Google Chrome\";v=\"124\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to fetch remote markdown file", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Remote server returned non-200 status", http.StatusBadGateway)
		return
	}

	// 4. Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	// 5. Construct the local file path inside the ./data vault
	fullPath := filepath.Join("./data", targetPath+".md")

	// 6. Create the directory tree if needed
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		http.Error(w, "Failed to build local directory tree", http.StatusInternalServerError)
		return
	}

	// 7. Write the raw bytes to disk
	if err := os.WriteFile(fullPath, body, 0644); err != nil {
		http.Error(w, "Failed to save file to vault", http.StatusInternalServerError)
		return
	}

	// 8. Reload the store to pick up the new file
	if err := h.store.Reload(); err != nil {
		http.Error(w, "File saved, but failed to reload navigation store", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("SUCCESS: Mirrored to /" + targetPath + "\n"))
}

// handlePeerSync securely pulls a file from a trusted peer, dropping it if corrupted.
func (h *handler) handlePeerSync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	peerURL := r.FormValue("url")
	targetPath := r.FormValue("path")

	if peerURL == "" || targetPath == "" {
		http.Error(w, "Missing 'url' or 'path'", http.StatusBadRequest)
		return
	}

	// 1. Fetch from the remote peer
	resp, err := http.Get(peerURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Peer unreachable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 2. Read the raw bytes
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read payload", http.StatusInternalServerError)
		return
	}

	// 3. SOVEREIGN VERIFICATION: Calculate our own hash of the received bytes
	calculatedHash := sha256.Sum256(body)
	calculatedString := hex.EncodeToString(calculatedHash[:])

	// 4. Extract the peer's claimed hash from the headers
	claimedHash := resp.Header.Get("X-RED-Content-Hash")

	// 5. THE GATEKEEPER: If the hashes do not match, the payload was tampered with
	if claimedHash != "" && calculatedString != claimedHash {
		// Drop the payload instantly. Do not write to disk.
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("SECURITY BREACH: Cryptographic signature mismatch. Payload dropped.\n"))
		return
	}

	// 6. If safe, write to the local vault
	fullPath := filepath.Join("./data", targetPath+".md")
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		http.Error(w, "Filesystem error", http.StatusInternalServerError)
		return
	}

	if err := os.WriteFile(fullPath, body, 0644); err != nil {
		http.Error(w, "Failed to write verified file", http.StatusInternalServerError)
		return
	}

	// 7. Hot-reload the engine
	h.store.Reload()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("SYNC VERIFIED: " + targetPath + " safely imported.\n"))
}
