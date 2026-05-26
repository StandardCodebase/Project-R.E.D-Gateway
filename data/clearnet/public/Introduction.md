---
title: "Introduction to Project R.E.D."
authors: ["R.E.D. Collective"]
created_at: "2026-01-10"
updated_at: "2026-05-20"
last_editor: "R.E.D. Collective"
---

Project R.E.D. (Resilient Encrypted Distribution) is a self-hosted knowledge node system built on static markdown files served by a lightweight Go binary. No database, no CMS, no JavaScript framework.

## How it works

Each `.md` file in the `data/` directory becomes a guide accessible at `/guides/<path>`. The engine reads the file, parses YAML frontmatter, renders markdown to HTML, and computes a SHA-256 hash of the raw content for integrity verification.

The hash is injected into the response header `X-RED-Content-Hash` and shown in the page metadata, allowing any reader to verify the content has not been tampered with.

## File structure

Place markdown files anywhere under `data/`. Subdirectories become part of the URL path. A file at `data/public/Level 1/soldering.md` is served at `/guides/public/Level 1/soldering`.

## Frontmatter fields

| Field | Required | Description |
|---|---|---|
| `title` | yes | Display title |
| `authors` | yes | List of authors |
| `created_at` | yes | ISO date |
| `updated_at` | no | ISO date |
| `last_editor` | no | Name |
| `discussion_hub` | no | URL for external discussion |

## Running the node

```bash
go run .
```

The server starts on port `8080` by default. Override with `RED_PORT`, `RED_NODE_NAME`, and `RED_DATA_DIR` environment variables.
