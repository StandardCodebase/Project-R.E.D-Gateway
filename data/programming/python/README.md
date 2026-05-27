# Obsidian Red Signer (Project R.E.D. Network)

**The official authoring bridge for the Project R.E.D. decentralized knowledge base.**

This plugin allows maintainers to cryptographically sign their Markdown notes with Ed25519 directly inside Obsidian. It automatically generates and updates the network's `manifest.json`, securely storing the file hash, your public key, and your cryptographic signature.

---

## 🦅 The Philosophy: Why We Built This

Project R.E.D. is built on a strict engineering philosophy: **stateless, lightweight, and execution-focused.** When we designed the security grid for the network, we needed a visual interface for maintainers to hash, sign, and manage their guides. The standard industry reaction is to reinvent the wheel—to waste weeks building a bloated, custom desktop application in C++ and Qt6 just to render a file tree and click a "Sign" button. 

We rejected that. 

Instead of building an authoring environment from scratch, we simply made use of one of the best local-first Markdown editor in the world: Obsidian. By turning our signing tool into an Obsidian plugin, we achieved zero context-switching. Maintainers can write their guides, hit a hotkey, and have the Ed25519 cryptography handled completely invisibly in the background. 

We build tools that work. We don't reinvent the wheel.

---

## ⚡ Features

* **Zero-Friction Signing:** Sign your files with one click via the ribbon icon, editor menu, or command palette.
* **The Sovereign Identity Vault:** Automatically generates and stores your permanent `maintainer.key` safely outside your working directory.
* **Network Manifest Injection:** Automatically discovers and updates the `manifest.json` at your vault root, formatting the keys exactly as the Project R.E.D. Go engine requires.
* **Real-Time Security Status:** A status bar indicator that displays `✓ Signed` (green) or `Unsigned` (gray) by dynamically comparing your live file hash against the manifest.
* **Public Key Clipboard:** Instantly copy your public key to add to the server’s `TrustedMaintainers` ring.

---


We need to rewrite the Installation & Setup section of the README to match the updated plugin logic (no renaming required, platform‑specific binary names).
Also, we must provide two clear methods: a simple drag‑and‑drop for inexperienced users, and a lightweight method that avoids downloading the whole zip.

Below is the updated section. Replace the existing ## Installation & Setup in your README.md with this:

markdown
## 🛠️ Installation & Setup

Choose the method that fits your workflow.

---

### 1. Script (recommended)

- Before you begin installation make sure you have `npm` installed.

#### macOS / Linux

Open a terminal and run:

```bash
bash <(curl -s https://raw.githubusercontent.com/StandardCodebase/obsidian-red-signer/main/install-red-signer.sh)
```



#### Windows

Open PowerShell as Administrator and run:


```powershell
iex (iwr -UseBasicParsing https://raw.githubusercontent.com/StandardCodebase/obsidian-red-signer/main/install-red-signer.ps1).Content
```


### 2.Manual Installation

1. Go to the Releases page.

2. Download the latest red-signer.zip and the binary for your operating system:

| Your OS              | Binary Name            |
| -------------------- | ---------------------- |
| Linux                | signer-linux-x64       |
| macOS(Intel)         | signer-windows-x64.exe |
| macOS(Apple Silicon) | signer-macos-arm64     |
| Windows              | signer-windows-x64     |

3. Extract red-signer.zip → you get a folder named red-signer.

4. Move the downloaded binary (e.g., signer-linux-x64) inside that red-signer folder.

5 Move the whole red-signer folder into your Obsidian vault’s plugins directory:
`YourVault/.obsidian/plugins/`

#### macOS / Linux only – make the binary executable:


```bash
chmod +x /path/to/YourVault/.obsidian/plugins/red-signer/signer-*
```


6. Restart Obsidian (or reload community plugins) and enable Red Signer in Settings → Community plugins.


✅ No renaming needed – the plugin now uses the exact binary names listed above.


### Next Steps

- Open any Markdown note in your vault.

- Click the signature icon (✍️) in the left ribbon, or use the command palette.

- Copy your public key and add it to your Project R.E.D. node’s TrustedMaintainers.

- Sign your first note – the status bar will show ✓ Signed.

- That’s it. You’re ready to contribute to the network.

---

## 🚀 The Genesis Workflow (First Use)

1. Open any Markdown guide in your vault.
2. Click the **Signature Icon** (✍️) in the left ribbon.
3. A modal will display your newly generated **Public Key**. Copy this and add it to your Project R.E.D. node's `TrustedMaintainers` map.
4. Click **Sign this note**.  

**What happens under the hood:**
* If you are a new contributor/maintainer, the engine generates your private key at `~/.red-network/maintainer.key` (with strict `0600` permissions).  
* It creates or locates the `manifest.json` in your vault root.  
* The guide is cryptographically hashed, signed, and instantly injected into the network ledger.

Once signed, the bottom status bar will glow **`Signed`**. If you modify a single character in that file, the hash changes, and the status will immediately revert to **`Unsigned`** until you authorize the changes.

---

## 🛡️ Security Architecture

* **The Private Key:** Your `maintainer.key` is stored at `~/.red-network/maintainer.key`  (owner read/write only). **Never upload this or share it.**
* **The Verification:** The Obsidian plugin only signs the files. Verification happens strictly on the Project R.E.D. Go server side, ensuring no compromised files are ever rendered to the end-user.

---

## 💻 Building from Source (For Contributors)

If you want to audit or modify the plugin's TypeScript architecture:

```bash
git clone [https://github.com/StandardCodebase/obsidian-red-signer](https://github.com/StandardCodebase/obsidian-red-signer)
cd obsidian-red-signer
npm install
npm install typescript
npm run build   # or `npx tsc`
```
