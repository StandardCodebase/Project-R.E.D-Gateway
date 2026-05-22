
# Project R.E.D

## Sovereign Knowledge Node Engine

Project R.E.D. rejects both centralized database monopolies and overly complex distributed consensus protocols. Instead, it systematically decouples the **Independent Data Layer** from the **Social Curation Layer**.

The engine operates as a stateless, high-performance Go runtime that compiles raw Markdown files into visually polarized technical templates, dynamically injecting cryptographic integrity signatures on every request loop.

## 1. The Philosophy: Red Sovereignty vs. Blue Illusion

The internet is broken not from a lack of information, but from a fundamental crisis of architecture. Well-meaning idealists continually attempt to fix this by proposing centralized chokepoints like B.L.U.E. (Broad Learning Universal Education)—a theoretical centralized archive designed to hold all human "how-to" knowledge in a single, strictly moderated, unified database.

While noble in its intent to provide universal free education, the BLUE System suffers from structural flaws that guarantee its immediate compromise or destruction:

- **The Single Point of Failure:** By attempting to gather all technical knowledge onto a single domain, the Blue System creates an existential target for corporate lawfare and state suppression.
- **The Conflict of Interest:** The moment a platform's revenue is derived from a percentage of material sales via affiliate links, its neutrality is permanently compromised.
- **The Jury-Style Delusion:** Relying on an uncompensated "jury" of randomly selected expert verifiers to audit complex technical workflows is an economic fantasy.
- **The Dependency Hell of Linear Hierarchies:** Forcing all human knowledge into a rigid, linear hierarchy where Level 4 strictly depends on Level 1 creates a catastrophic maintenance bottleneck.

## 2. How R.E.D Fixes B.L.U.E


Project R.E.D. systematically dismantles the architectural vulnerabilities inherent in centralized knowledge repositories.

- **Eradicating the Single Point of Failure:** B.L.U.E. relies on a master domain, creating a massive target for corporate lawfare and global de-indexing. R.E.D. operates on a stateless, dual-tier Docker network topology (Clearnet/Tor), eliminating the centralized attack vector entirely.

- **Eliminating Financial Conflicts:** B.L.U.E.'s reliance on contextual advertising fundamentally destroys its neutrality by financially incentivizing the use of expensive corporate materials. R.E.D. requires zero centralized funding, ensuring information remains free from commercial manipulation and review-bombing bots.
  
- **Beyond Zero Funding – Ethical Node Monetization:** While R.E.D. the engine takes no cut, individual node operators may choose to include their own affiliate links or donation addresses within the guides they host. The revenue goes directly to the operator who maintains the node and curates the content – never to a central entity. Readers can freely choose a different node if they prefer an ad‑free or different affiliate experience. This decentralized, opt‑in model aligns incentives without corrupting the knowledge layer.

- **Outsourcing the Verifier System:** B.L.U.E. expects specialized professionals to donate countless hours to an internal, uncompensated moderation jury. R.E.D. outsources the voting and review of guides to highly capable, already-existing social platforms like Reddit and Discord. By tapping into established, free networks, it leverages active communities to organically filter content and create a more robust, unbiased verifier system.

- **Bypassing Dependency Hell:** B.L.U.E. enforces a rigid, linear learning hierarchy where a single obsolete foundational guide can collapse the entire structure. R.E.D. prevents this by leveraging native filesystem directories and dynamic versioning, allowing knowledge to adapt organically without cascading failures.

- **Resolving the "Spin-Off" Paradox:** B.L.U.E. mandates one definitive guide per topic, yet paradoxically suggests forking contested guides during disputes, guaranteeing a cluttered, redundant database. R.E.D. removes moderation logic from the runtime entirely; the end-user's local client seamlessly curates the best version based on established network consensus.

## 3. Architecture: The Dual-Tier Deployment

Project R.E.D operates via a highly secure, dual-tier Docker network topology. It separates public "Clearnet" knowledge from restricted "Darknet" knowledge using strict Docker network isolation.

```text
GLOBAL INTERNET & TOR NETWORK
       │                   │
  [Port 80/443]       [Tor Network Bridge]
       │                   │
+------▼-------+    +------▼------------------+
| Caddy Proxy  |    | Tor Sidecar Container   |
+------┬-------+    +------┬---------┬--------+
       │                   │         │
[Clearnet-Tier Network]    │    [Onion-Tier Network (Internal)]
       │                   │         │
+------▼-------+           │  +------▼-------+
|  Light Node  |           │  |  Dark Node   |
| (Go Engine)  |           │  | (Go Engine)  |
+------┬-------+           │  +------┬-------+
       │                   │         │
       +----------+        │         |
                  │        │         │
               +--▼--------▼---------▼--+
               |   Host Volume (/data)  |
               | - /public              |
               | - /restricted          |
               +------------------------+

```

### Key Security Mechanisms:

* **Stateless Web Routing:** No database layer (SQL/NoSQL) exists to be corrupted, breached, or subpoenaed.
* **Cryptographic Validation:** The engine computes a standard SHA-256 hash over the raw file bytes on demand, injecting it into response headers (`X-RED-Content-Hash`).
* **Air-Gapped Vault:** The `Dark Node` lives on a strictly internal Docker network. It has zero outbound internet access. The `Tor Sidecar` acts as the exclusive cryptographic bridge, publishing the hidden service descriptor to the global Tor network while routing incoming requests securely to the Vault.

## 4. Technical Stack Blueprint

* **Language & Runtime:** Go (Golang) — Compiles down to a single memory-safe static binary with zero virtual machine overhead.
* **Markdown Parser:** `goldmark` — Highly efficient, fully CommonMark compliant.
* **Container Environment:** Multi-stage minimal Docker pipeline (`alpine`) for complete process isolation and instant replication.

## 5. Installation & Deployment

For full setup instructions—including configuring your local workspace, building the Docker matrix, and launching both the Clearnet Gateway and the Air-Gapped Onion Vault—please read the installation guide:

**[Node Operator Installation Guide](INSTALL.md)**.

## 6. Simulating Knowledge Base Structures

Rather than enforcing a brittle, top-down linear hierarchy that breaks under dependency hell, R.E.D. leverages your computer's native filesystem directory tracking. Group, version, and fork your files within folders inside your local storage volume dynamically:

```plaintext
data/
└── restricted/
    └── solar-array-build/
        ├── 00-index.md
        ├── 01-wiring-schematics-v1.0.md
        └── 02-inverter-fusing.md

```

The node pathing mechanics resolve these automatically into accessible downstream routes:

* `http://<your-onion>.onion/guides/solar-array-build/00-index`

## 7. Curation Philosophy: The Sovereign Web-of-Trust

1. **The Software Only Knows State:** The Go runtime handles zero moderation logic. It treats strings objectively.
2. **Cryptographic Peer Review:** Instead of relying on a centralized platform's voting buttons, users and curators sign content hashes with their independent cryptographic keys (PGP/Nostr).
3. **End-User Agency:** You, the reader, choose which curators to trust. Your local index client aggregates links signed by your trusted network.

The creator of the Blue System lamented that he lacked the millions of dollars needed for a startup, resigning himself to the belief that his vision would never be a reality. He failed because he thought like a corporate founder seeking capital.

We do not need capital. We do not need permission. We do not need a master website.

It's about time we stop waiting for someone or something to fix our problems for us. No hero is coming.

**Claim Your Agency. Run a Node.**

---

## License

This project is licensed under the GNU Affero General Public License v3.0 - see the [LICENSE.md](LICENSE.md) file for details.
