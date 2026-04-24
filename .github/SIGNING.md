# Windows Binary Signing in CI

The project is configured to sign the Windows binary and publish a GitHub Release in `.github/workflows/tag-build.yml`.

Signing and release publishing run on pushing a `v*` tag after `wails build`.

## What to Configure in GitHub

### Secrets

- `WIN_SIGNING_CERT` — base64-encoded certificate content (`.pfx`/`.p12`)
- `WIN_SIGNING_CERT_PASSWORD` — certificate password

### Variables (optional)

- `WIN_SIGNING_TIMESTAMP_URL` — timestamp server URL  
  If not set, the default value `http://ts.ssl.com` is used.

## How to Prepare a Base64 Certificate (PowerShell)

```powershell
certutil -encode .\my-cert.p12 my-cert-base64.txt
```

Put the contents of `my-cert-base64.txt` into the `WIN_SIGNING_CERT` secret.

## How the Pipeline Works

1. Build the application with `wails build` using version from the tag via `-ldflags`.
2. Decode the certificate from the secret and create a `.pfx` on the runner.
3. Sign `build/bin/game-radar.exe` via `signtool`.
4. Verify signature (`signtool verify`).
5. Prepare release assets:
   - `game-radar-windows-amd64.exe` (signed binary),
   - `game-radar-windows-amd64.zip` (archive with binary),
   - `checksums.txt` (SHA256).
6. Publish GitHub Release with auto-generated release notes.
7. GitHub automatically adds source code archives (`zip` and `tar.gz`) to the release.
8. Clean up temporary certificate files.

## Release Notes Template

Automatic release-note categories are configured in `.github/release.yml`.
