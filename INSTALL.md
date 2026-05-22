# Project R.E.D: Node Operator Installation Guide

Welcome to the matrix. This guide will walk you through initializing your own sovereign knowledge node, structuring your Markdown data, and booting both the Clearnet Gateway and the Air-Gapped Onion Vault.

## Prerequisites

Before you begin, ensure your host machine has the following installed:

* **Docker Engine** (v20.10+)
* **Docker Compose** (v2.0+)
* A basic text editor for writing Markdown files.

---

## Step 1: Establish the Data Vaults

Project R.E.D uses strict Docker volume mounting to read your knowledge base directly from your host machine. **You do not need to rebuild the Docker container when you add or edit files.** However, because the system isolates public internet traffic from the darknet vault, you must strictly organize your files into two separate directories.

Run this command in the root of your project to build the necessary folders:

```bash
mkdir -p data/public data/restricted tor_keys

```

---

## Step 2: Seed Your Knowledge Base

The Go engine expects standard Markdown files (`.md`) equipped with YAML front-matter.

Create a test file for your Onion Vault to verify the connection later:

**1. Create a file at:** `./data/restricted/vault-test.md`
**2. Paste the following blueprint:**

```markdown
---
title: "Dark Node Initializer"
author_identity: "red://alpha-omega-99"
created_at: 2026-05-22
discussion_hub: "nostr://vault-discussions"
---
# Welcome to the Onion Vault

If you are reading this over the Tor network, the cryptographic volume mount was successful. The files in your host's `/data/restricted/` directory are actively syncing.

```

*(Note: If you want to host files on the standard internet, place them inside `./data/public/` instead.)*

---

## Step 3: Boot the Docker Matrix

You are now ready to launch the system. The `docker-compose` stack will automatically build the Go engine from source, establish the isolated internal networks, and spin up the Tor sidecar.

Run the following command from the root of the project:

```bash
docker-compose up -d --build

```

* The `-d` flag runs the node silently in the background.
* The `--build` flag ensures the static Go binaries are freshly compiled.

---

## Step 4: Retrieve Your Sovereign Addresses

Your nodes are now live. Because the Go routing engine serves content under the `/guides/` path, you must append that to your URLs to view your files.

### Accessing the Clearnet Gateway (Local)

If you placed files in `./data/public/`, you can view them locally via your standard web browser at:
`http://localhost:8080/guides/<filename>`

### Accessing the Air-Gapped Onion Vault

Upon its first boot, the Tor sidecar generates a brand new V3 hidden service address and saves it to your host machine.

To reveal your Vault's `.onion` link, run:

```bash
cat ./tor_keys/hostname

```

Copy the output string and open your Tor Browser. To view the test file you created in Step 2, navigate to:
`http://<YOUR_ONION_ADDRESS_HERE>.onion/guides/vault-test`

**⚠️ Crucial Tor Networking Note:** * Do not add `.md` to the end of the URL in your browser.

* If the Tor browser throws an **"Onion Site Not Found"** error immediately after booting, *be patient*. It takes 5 to 15 minutes for a brand-new cryptographic descriptor to fully propagate across the global Tor network.

---

## Step 5: Shutting Down & Maintenance

To safely spin down the entire node matrix without losing your data or cryptographic keys:

```bash
docker-compose down

```

To view the live access logs of your Tor sidecar or Go engine:

```bash
# View Go Engine logs
docker logs -f red_dark_node

# View Tor connection logs
docker logs -f tor_sidecar

```
