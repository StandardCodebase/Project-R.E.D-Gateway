package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"red-engine/internal/config"
	"red-engine/internal/handler"
)

func main() {
	cfg := config.Load()

	tmpl, err := template.ParseFiles("templates/layout.html")
	if err != nil {
		log.Fatalf("failed to parse template: %v", err)
	}

	node := &handler.Node{
		Cfg:      cfg,
		Template: tmpl,
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("/health", node.Health)
	mux.HandleFunc("/guides/", node.RenderGuide)
	mux.HandleFunc("/download/", node.DownloadGuide)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Project R.E.D. Engine [%s] initiating on port :%s...\n", cfg.NodeName, cfg.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("engine panic: %v", err)
	}
}
