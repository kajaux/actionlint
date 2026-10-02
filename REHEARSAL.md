# Isolated release rehearsal

Source stack: `2cf1ae998a21b78e8471945ed2fd18be4903f1f0` (#204).
Destination: `kajaux/actionlint-release-rehearsal` only.

Exercise signed immutable GitHub releases, the published Action on Linux/macOS/Windows/ubuntu-slim,
npm packages through GitHub Packages, and CLI/Action containers through GHCR.
No npmjs, Docker Hub, Homebrew, Scoop, AUR, or WinGet publication. Nix checks excluded.

The rehearsal uses version `1.17.0` in this fork's independent release namespace.
It does not move any production tags or publish any production packages.
GitHub Packages publication does not exercise npmjs trusted publishing.
