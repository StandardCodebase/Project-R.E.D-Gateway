
# Project R.E.D. – Node Operator Installation Guide

This document explains how to set up, configure, and run a Project R.E.D. sovereign knowledge node – including how to write guides with proper metadata.

## Prerequisites

- **Docker** and **Docker Compose** (for production deployment)
- **Go 1.26+** (only for local development)
- Basic familiarity with the command line and markdown

---

## Quick Start (Development)

1. Clone the repository:
   ```bash
   git clone https://github.com/your-org/project-red.git
   cd project-red
   ```

2. Create a test guide:
   ```bash
   mkdir -p data/public
   ```

   Create `data/public/hello.md` with the following content (see [Metadata section](#writing-guide-metadata-yaml-frontmatter) for details):
   ```yaml
   ---
   title: "Hello, Sovereign World"
   authors: ["Your Name"]
   created_at: "2025-01-01"
   discussion_hub: "https://discord.gg/example"
   ---
   # Welcome

   This is your first R.E.D. guide.
   ```

3. Run the engine:
   ```bash
   RED_DATA_DIR=./data/public go run main.go
   ```

4. Open `http://localhost:8080/guides/hello`

---

## Production Deployment with Docker Compose

The repository includes a `docker-compose.yml` that launches two independent nodes:

- **Clearnet node** (public internet) – served via Caddy proxy on ports 80/443
- **Darknet node** (Tor hidden service) – isolated on an internal network, accessible only via Tor

### Steps

1. Build the Docker image:
   ```bash
   docker build -t project-red-node .
   ```

2. Prepare your data directories:
   ```bash
   mkdir -p data/public data/restricted
   # Add your .md files into these folders
   ```

3. (Optional) Edit `Caddyfile` to set your domain.

4. Launch the stack:
   ```bash
   docker-compose up -d
   ```

5. Check logs:
   ```bash
   docker-compose logs -f
   ```

6. The clearnet node will be available at `http://localhost` (or your domain).  
   The darknet node’s `.onion` address appears in the Tor container logs:
   ```bash
   docker logs tor_sidecar
   ```

---

## Writing Guide Metadata (YAML Frontmatter)

Every markdown guide **must** start with a YAML frontmatter block between `---` lines. This metadata is displayed in the node’s web interface and used for content integrity.

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `title` | string | Human‑readable title of the guide |
| `authors` | list of strings | Primary author(s) of the content (displayed as “AUTHORS”) |
| `created_at` | string | ISO date `YYYY-MM-DD` when the guide was first written |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `contributors` | list of strings | Additional people who helped but aren’t primary authors |
| `updated_at` | string | Date of last modification (ISO `YYYY-MM-DD`) |
| `last_editor` | string | Name or handle of the person who made the latest change |
| `discussion_hub` | string | URL to a community thread (Reddit, Discord, Nostr, etc.) for peer review |

### Example with Multiple Authors & Contributors

```yaml
---
title: "Emergency Lithium-Ion Cell Bypass"
authors: ["Alex Rivera", "Sam Chen"]
contributors: ["Jordan Lee", "Taylor Smith"]
created_at: "2026-05-21"
updated_at: "2026-05-22"
last_editor: "Sam Chen"
discussion_hub: "https://reddit.com/r/ProjectRED/comments/xyz123"
---
# Emergency Lithium-Ion Cell Bypass

Content goes here...
```

### How to Update a Guide

1. Edit the markdown file.
2. Change `updated_at` to today’s date.
3. Set `last_editor` to your name or handle.
4. If you are a new contributor, add your name to the `contributors` list (or `authors` if you become a primary author).
5. Commit the change (if using git) or simply restart the node.

> **Note:** The node does not enforce any rules on names – they are free‑text. Use consistent handles to build reputation within your community.

### Display Behaviour

- If `updated_at` and `last_editor` are present, the metadata ribbon shows an **UPDATED** line and a **LAST EDITOR** line.
- The `discussion_hub` appears as a button at the bottom of the page: “Review or Discuss this Guide on Public Hub”.
- The **SHA-256** hash is automatically computed from the **raw file bytes** (including frontmatter). This allows external verifiers to compare hashes.

---

## Verifying Content Integrity

Every rendered page includes an HTTP header `X-RED-Content-Hash` and a visible hash string. To verify a guide:

1. Download the raw markdown (click **📄 Markdown Source**).
2. Compute its SHA-256 hash:
   ```bash
   sha256sum guide.md
   ```
3. Compare with the hash shown on the node’s page.

If the hash matches, the content has not been altered since it was published by the node operator.

For community‑driven verification, publish the hash of a trusted version on an external platform (e.g., Discord, Nostr, signed git tag). Readers can then compare any node’s hash against the community‑approved hash.

---

## Next Steps

- Join the Project R.E.D. Discord (link in main README).
- Set up your own node and start writing guides.
- Help others by verifying hashes and discussing content.
