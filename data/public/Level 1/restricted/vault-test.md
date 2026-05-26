---
title: "Verifying Node Integrity"
authors: ["R.E.D. Collective"]
created_at: "2026-04-05"
---

This guide explains how to independently verify that the content served by a R.E.D. node matches the original source files.

## Using the content hash

Every guide page displays a SHA-256 hash in its metadata bar and in the `X-RED-Content-Hash` response header. To verify:

```bash
# Download the raw file
curl -O http://localhost:8080/download/public/Level\ 1/restricted/vault-test

# Compute its hash
sha256sum vault-test.md

# Compare against the hash shown on the page
```

If the hashes match, the file was not modified between disk and your browser.

## Checking headers directly

```bash
curl -I http://localhost:8080/guides/public/Introduction
```

Look for `X-RED-Content-Hash` in the output.

## What this does not protect against

The hash proves the content is what the server sent. It does not prove the server was not compromised. For stronger guarantees, sign the markdown files with GPG and publish public keys out-of-band.
