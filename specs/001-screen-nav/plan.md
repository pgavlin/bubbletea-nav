# Implementation Plan: Screen Navigation Library

**Branch**: `001-screen-nav` | **Date**: 2026-02-14 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-screen-nav/spec.md`

## Summary

A Go library providing two composable primitives for Bubble Tea
applications: a **navigation stack** for push/pop/replace screen
transitions with lifecycle events, and a **focus manager** for
Tab/Shift+Tab and mouse-click-based focus cycling between
interactive components within a screen. Targets Bubble Tea v1
(stable v1.3.10).

## Technical Context

**Language/Version**: Go 1.25+ (module minimum TBD by `go.mod`)
**Primary Dependencies**: `github.com/charmbracelet/bubbletea` v1.3.x
**Storage**: N/A (in-memory state only)
**Testing**: `go test ./...` with table-driven tests
**Target Platform**: Any OS with terminal (Linux, macOS, Windows)
**Project Type**: Single Go module / library
**Performance Goals**: Navigation operations (push/pop/replace/focus)
  complete in <1ms; no perceptible latency for end users.
**Constraints**: Zero external dependencies beyond Bubble Tea itself.
  No goroutines or I/O outside of `tea.Cmd` functions (Constitution
  Principle II).
**Scale/Scope**: Library supporting arbitrary stack depth and up to
  ~100 focusable components per screen (practical TUI limit).

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1
design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Library-First | PASS | Single Go package, no `main`, importable via `go get`. |
| II. Bubble Tea Idioms | PASS | Stack implements `tea.Model`. FocusManager follows the bubbles value-type pattern (Update returns own type). State via Msg/Cmd. No goroutines. |
| III. Composable API | PASS | Two small types (Stack, FocusManager) compose independently. Narrow interfaces (Screen: 3 methods, Focusable: 3 methods). |
| IV. Test-Required | PASS | Table-driven tests for all exported types/functions. `go test ./...` only. |
| V. Backward Compatibility | PASS | Initial release; semver from v0.1.0. |
| VI. Observability | PASS | Stack and FocusManager implement `fmt.Stringer`. Navigation Msgs are inspectable. |
| VII. Simplicity | PASS | Single package. Two types. No config objects or option patterns unless justified. |

No violations. Complexity Tracking section not needed.

## Project Structure

### Documentation (this feature)

```text
specs/001-screen-nav/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
│   └── api.go           # Exported types and interfaces (pseudo-code)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
nav.go                   # Stack type, Screen interface, nav messages
focus.go                 # FocusManager, Focusable interface, focus msgs
lifecycle.go             # Lifecycle event types and processing
nav_test.go              # Stack tests (table-driven)
focus_test.go            # FocusManager tests (table-driven)
lifecycle_test.go        # Lifecycle event tests
go.mod                   # Module definition
go.sum                   # Dependency checksums
_examples/
└── basic/
    └── main.go          # 2-screen push/pop demo with focus
```

**Structure Decision**: Single package (`package nav`) at the module
root. All source files share one package. The `_examples/` directory
contains runnable demos with `package main`. This follows Go library
conventions and Constitution Principle I (no `main` in core module)
and Principle VII (simplest structure that works).
