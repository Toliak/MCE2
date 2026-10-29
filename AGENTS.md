# AGENTS.md

## Project

MCE2 (Make Configuration Easier) is a Go CLI/TUI that installs and configures developer tooling. Work is done in "Tegns" (individual tools/configs), grouped into "Tegnsetts" (ordered collections). User-facing configuration data lives in the separate `toliak/mce2config` repository, cloned as a sibling `../config` directory.

## Layout

- `configapp/` — all Go source; run commands from here.
- `configapp/cmd/mce/` — CLI and tview-based TUI.
- `configapp/cmd/mce-e2e/`, `configapp/e2e/` — Docker end-to-end tests and overlay.
- `configapp/tegnbuilder/` — Tegn/Tegnsett interfaces, parameters, and ordering.
- `configapp/tegns/` — concrete Tegns, downloads, and registry in `tegns.go`.
- `configapp/inspector/`, `configapp/osinfo/`, `configapp/platform/` — OS detection, platform helpers, exec/fs/git/package-manager wrappers.
- `configapp/sedparody/` — sed-like config editing helpers.

## Build and test

Run from `configapp/`:

- `make build` — local binaries (`mce`, debug `mce-dbg`, `mce-e2e`).
- `make build_release` — cross-platform release binaries.
- `make test_e2e` — Docker e2e; requires `../config` cloned from `https://github.com/toliak/mce2config`.
- `make test_pkg` — package-manager tests in Docker.

There are currently no unit tests. CI builds all binaries and runs Docker e2e.

## Conventions

- Go 1.25.6; keep `gofmt` clean.
- Default to `CGO_ENABLED=0` and preserve all cross-platform targets.
- Prefer `platform` helpers over raw `os/exec`, filesystem, git, or package-manager calls.
- Wrap errors with context using `fmt.Errorf("...: %w", err)`.
- Add new Tegns under `configapp/tegns/tegn/`; scaffold with `create_new_tegn.py`, implement `tegnbuilder.Tegn`, and register in `configapp/tegns/tegns.go`.
- Uninstall paths are destructive: preserve user data and avoid removing unrelated files.
- Do not commit build artifacts or binaries.

## Nullability

See `configapp/NULLABILITY.md`. Pointers are non-nil by default; optional values use `optional.Optional[T]`.
