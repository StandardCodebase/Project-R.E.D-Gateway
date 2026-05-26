package render

import (
	"crypto/sha256"
	"encoding/hex"
	"html/template"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"

	"red-engine/internal/config"
)

type Result struct {
	Meta        config.PostMetadata
	HTMLContent template.HTML
	Hash        string
}

func Markdown(raw []byte) (*Result, error) {
	sum := sha256.Sum256(raw)
	hash := hex.EncodeToString(sum[:])

	var meta config.PostMetadata
	body, err := frontmatter.Parse(bytes.NewReader(raw), &meta)
	if err != nil {
		return nil, err
	}

	var buf strings.Builder
	if err := goldmark.Convert(body, &buf); err != nil {
		return nil, err
	}

	return &Result{
		Meta:        meta,
		HTMLContent: template.HTML(buf.String()),
		Hash:        hash,
	}, nil
}
