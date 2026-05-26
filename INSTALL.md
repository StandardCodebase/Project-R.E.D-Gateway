# Project R.E.D. – Node Operator Installation Guide

This document explains how to set up, configure, and run a Project R.E.D. sovereign knowledge node – including how to write guides with proper metadata across different operating systems and terminals.

---

## 🛠️ Prerequisites

Before you begin, ensure you have the following installed based on your workflow:

* **Docker & Docker Compose** (Highly recommended for production deployment)
* **Go 1.26+** (Only required for local source development)
* **Terminal Environment:**
    * **Linux/macOS:** Bash or Zsh
    * **Windows:** PowerShell 7+ (Recommended) or Command Prompt (CMD)

---

## 🚀 Quick Start (Development)

Follow the setup steps corresponding to your specific operating system and terminal interface below.

### 1. Clone the Repository (All Platforms)

```bash
git clone [https://github.com/your-org/project-red.git](https://github.com/your-org/project-red.git)
cd project-red

```

### 2. Create the Public Data Directory

Choose the block matching your terminal environment:

#### 🐧 Linux / 🍏 macOS (Bash/Zsh)

```bash
mkdir -p data/public

```

#### 🪟 Windows (PowerShell)

```powershell
New-Item -ItemType Directory -Force -Path "data\public"

```

#### 🪟 Windows (Command Prompt - CMD)

```cmd
mkdir data\public

```

### 3. Create a Test Guide

Create a file named `hello.md` inside your newly created `data/public/` folder. Add the following content exactly as shown:

```yaml
---
title: "Hello, Sovereign World"
authors: ["Your Name"]
created_at: "2025-01-01"
discussion_hub: "[https://discord.gg/example](https://discord.gg/example)"
---
# Welcome

This is your first R.E.D. guide.

```

### 4. Boot up the Engine

Run the local Go development server by initializing the data path inline:

#### 🐧 Linux / 🍏 macOS (Bash/Zsh)

```bash
RED_DATA_DIR=./data/public go run main.go

```

#### 🪟 Windows (PowerShell)

```powershell
$env:RED_DATA_DIR=".\data\public"; go run main.go

```

#### 🪟 Windows (Command Prompt - CMD)

```cmd
set RED_DATA_DIR=.\data\public
go run main.go

```

### 5. Access the Node

Open your web browser and navigate to: **`http://localhost:8080/guides/hello`**

---

## 🐳 Production Deployment with Docker Compose

The repository includes a pre-configured `docker-compose.yml` that launches two completely isolated, independent environments:

* **Clearnet Node:** Public web-facing node proxying through Caddy on standard ports (`80`/`443`).
* **Darknet Node:** Isolated Tor hidden service running completely over the `.onion` network.

### Execution Steps

#### 1. Build the Production Engine Image

```bash
docker build -t project-red-node .

```

#### 2. Initialize Directory Layout

* **Linux / macOS:**
```bash
mkdir -p data/public data/restricted

```


* **Windows (PowerShell):**
```powershell
New-Item -ItemType Directory -Force -Path "data\public", "data\restricted"

```


* **Windows (CMD):**
```cmd
mkdir data\public
mkdir data\restricted

```



> *Drop your curated `.md` files directly into these directories before proceeding.*

#### 3. Domain Mapping (Optional)

Modify the project `Caddyfile` configuration block to reflect your personal domain or staging endpoints.

#### 4. Spin Up the Daemon Pipeline

```bash
docker-compose up -d

```

#### 5. Track Run Logs

```bash
docker-compose logs -f

```

#### 6. Verification & Onion Link Retrieval

* **Clearnet access endpoint:** `http://localhost/guides/index` (or your configured custom domain domain)
* **Darknet Onion Link extraction:** To view your freshly minted Tor hidden service string, parse the active sidecar logs:
```bash
docker logs tor_sidecar

```



---

## 📝 Writing Guide Metadata (YAML Frontmatter)

Every markdown document processing through the engine **must** be declared with an explicit YAML frontmatter block bound tightly between dual `---` separators.

This tracking system ensures seamless UI generation and strict cryptographic verification.

### Required Core Metadata

| Field | Type | Description |
| --- | --- | --- |
| `title` | `string` | The reader-facing clean title of your guide. |
| `authors` | `list of strings` | Primary documentation authors (rendered as **AUTHORS** ribbon). |
| `created_at` | `string` | Absolute origin tracking date formatted as strict ISO `YYYY-MM-DD`. |

### Extended Optional Metadata

| Field | Type | Description |
| --- | --- | --- |
| `contributors` | `list of strings` | Peer community editors who helped expand the documentation footprint. |
| `updated_at` | `string` | Date of the latest revision formatted in ISO `YYYY-MM-DD`. |
| `last_editor` | `string` | The distinct handle or nickname of the person who committed the edit. |
| `discussion_hub` | `string` | External verification link (Nostr, Reddit, Discord, etc.) for open peer review. |

### Complete Multi-Author Blueprint Example

```yaml
---
title: "Emergency Lithium-Ion Cell Bypass"
authors: ["Alex Rivera", "Sam Chen"]
contributors: ["Jordan Lee", "Taylor Smith"]
created_at: "2026-05-21"
updated_at: "2026-05-22"
last_editor: "Sam Chen"
discussion_hub: "[https://reddit.com/r/ProjectRED/comments/xyz123](https://reddit.com/r/ProjectRED/comments/xyz123)"
---
# Emergency Lithium-Ion Cell Bypass

Content goes here...

```

### Protocol for Updating Files

1. Modify the target markdown file data.
2. Bump the tracking date string under `updated_at` to the current timestamp.
3. Update the `last_editor` variable string.
4. Append additional identities to the `contributors` array if applicable.
5. Push changes down your standard version control channel or hit a fast-restart flag on the daemon structure.

> ⚠️ **Note on Identity:** Project R.E.D. does not run complex validation checks on names or identity strings—they are un-sanctioned free text. Choose a persistent handle to manually establish your cryptographic footprint and reputation within peer verification circles.

### Interface Behavior Mechanics

* **Dynamic Ribbons:** Setting an active `updated_at` and `last_editor` dynamically generates explicit **UPDATED** and **LAST EDITOR** data matrices in the web viewer.
* **Interactions:** Linking a valid `discussion_hub` string automatically injects an explicit interactive node action element at the bottom footer: `Review or Discuss this Guide on Public Hub`.
* **Cryptographic Integrity:** The application layer instantly registers a byte-level **SHA-256** checksum calculated directly from the absolute file payload (including headers). This permits distributed network nodes to easily confirm origin parity.

---

## 🔒 Verifying Content Integrity

Every page served through a R.E.D. engine endpoint ships with an external `X-RED-Content-Hash` header payload alongside a clearly displayed visual checksum string in the UI footer.

To cross-verify a file, parse the base file through a local crypto loop:

### 1. Fetch the Source Asset

Click on the **📄 Markdown Source** action link inside the document interface and download the raw `.md` content layout.

### 2. Compute the Local Target Checksum

#### 🐧 Linux

```bash
sha256sum guide.md

```

#### 🍏 macOS

```bash
shasum -a 256 guide.md

```

#### 🪟 Windows (PowerShell)

```powershell
Get-FileHash -Algorithm SHA256 .\guide.md

```

#### 🪟 Windows (Command Prompt - CMD)

```cmd
certutil -hashfile guide.md SHA256

```

### 3. Compare Results

Verify that the output hash match exactly with the hash shown on the web UI layout.

For robust community-driven defense, publish verified file hashes across external decentralized protocols (e.g., Nostr, signed Git tags, or PGP-clearsigned messages). Readers can cross-check their local node's hash against these signatures to spot altered instances.

---

## ⚡ Next Steps

* Deploy your instance, structure a custom layout, and spin up an active node pipeline.
* Audit peer endpoints, calculate independent hashes, and track community documentation logs.

No heroes are coming. **Claim your agency.**
