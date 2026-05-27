package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/RED-Collective/red-engine/internal/config"
	"github.com/RED-Collective/red-engine/internal/fetch"
	"github.com/RED-Collective/red-engine/internal/router"
	"github.com/RED-Collective/red-engine/internal/store"
)

func main() {
	cfgPath := flag.String("config", "config.json", "path to config file")
	pull := flag.Bool("pull", false, "fetch knowledge base before starting")
	flag.Parse()

	cfg := config.Default()
	if _, err := os.Stat(*cfgPath); err == nil {
		loaded, err := config.Load(*cfgPath)
		if err != nil {
			log.Fatalf("config: %v", err)
		}
		cfg = loaded
	}

	if *pull {
		log.Printf("pulling %s", cfg.SourceURL)
		if err := fetch.Pull(cfg.SourceURL, cfg.SourceType, cfg.DataDir); err != nil {
			log.Fatalf("fetch: %v", err)
		}
		log.Println("fetch complete")
	}

	s := store.New(cfg.DataDir)
	if err := s.Reload(); err != nil {
		log.Fatalf("store: %v", err)
	}

	h := router.New(s, cfg.SiteName)
	log.Printf("RED listening on %s", cfg.Addr)
	log.Fatal(http.ListenAndServe(cfg.Addr, h))
}
