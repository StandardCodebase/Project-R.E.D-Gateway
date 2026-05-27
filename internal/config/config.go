package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	SourceURL  string `json:"source_url"`
	SourceType string `json:"source_type"`
	DataDir    string `json:"data_dir"`
	Addr       string `json:"addr"`
	SiteName   string `json:"site_name"`
	// --- RESTORED: Node Keys --- DO NOT DELETE
	NodePublicKey string   `json:"node_public_key"`
	TrustedPeers  []string `json:"trusted_peers"`
}

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := &Config{}
	return cfg, json.NewDecoder(f).Decode(cfg)
}

func Default() *Config {
	return &Config{
		SourceURL:  "https://github.com/yourname/red-knowledge/archive/refs/heads/main.tar.gz",
		SourceType: "tar.gz",
		DataDir:    "./data",
		Addr:       ":8080",
		SiteName:   "RED",
	}
}
