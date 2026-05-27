package config

import "html/template"

type PostMetadata struct {
	Title         string   `yaml:"title"`
	Authors       []string `yaml:"authors"`
	Contributors  []string `yaml:"contributors"`
	CreatedAt     string   `yaml:"created_at"`
	UpdatedAt     string   `yaml:"updated_at"`
	LastEditor    string   `yaml:"last_editor"`
	DiscussionHub string   `yaml:"discussion_hub"`
}

type PageData struct {
	PostMetadata
	NodeName    string
	ContentHash string
	ContentPath string
	HTMLContent template.HTML
}

type GuideEntry struct {
	Path  string
	Title string
}

type ImportRequest struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
}
