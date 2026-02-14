# Implementation Plan: Focus Message Routing

**Branch**: `002-focus-message-routing` | **Date**: 2026-02-14 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/002-focus-message-routing/spec.md`

## Summary

Enhance the FocusManager to automatically route non-focus messages
to the currently focused component via the standard `tea.Model`
Update method. The existing `Focusable` interface is expanded to
embed `tea.Model`, requiring all focusable components to implement
`Init`, `Update`, and `View`. This is a breaking change; backward
compatibility is not required. No new interfaces are introduced.

## Technical Context

**Language/Version**: Go 1.25+ (per existing `go.mod`)
**Primary Dependencies**: `github.com/charmbracelet/bubbletea` v1.3.10
**Storage**: N/A (in-memory state only)
**Testing**: `go test ./...` with table-driven tests
**Target Platform**: Any OS with terminal (Linux, macOS, Windows)
**Project Type**: Single Go module / library (existing)
**Performance Goals**: Message routing adds <1ms latency; single
  method call + type assertion per message.
**Constraints**: Zero new dependencies. No goroutines or I/O outside
  `tea.Cmd`. No new interfaces. Breaking change to Focusable is
  accepted.
**Scale/Scope**: Incremental feature addition to existing library.
  Modifies `focus.go` and `focus_test.go`. Updates test helpers in
  `nav_test.go`, `lifecycle_test.go`, and example app.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1
design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Library-First | PASS | No new packages. Modified interface is importable via `go get`. |
| II. Bubble Tea Idioms | PASS | FocusManager remains a value type (bubbles convention). Routing uses standard `tea.Model.Update`. State via Msg/Cmd. No goroutines. |
| III. Composable API | PASS | No new types. Focusable interface expands from 3 to 6 methods (via tea.Model embedding). FocusManager composes with existing patterns. |
| IV. Test-Required | PASS | Table-driven tests for all new routing behavior. Existing tests updated. |
| V. Backward Compatibility | JUSTIFIED | Breaking change: Focusable embeds tea.Model. Spec explicitly states backward compatibility is not required. Semver major version bump needed. |
| VI. Observability | PASS | FocusManager.String() unchanged. Routing behavior transparent through commands. |
| VII. Simplicity | PASS | No new interfaces or types. One modified method. Standard tea.Model pattern reused. |

**Post-Phase 1 re-check**: All gates still pass. The Focusable
interface change is a justified breaking change per the spec.

## Project Structure

### Documentation (this feature)

```text
specs/002-focus-message-routing/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
│   └── api.go.txt       # Modified API surface (pseudo-code)
└── tasks.md             # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
focus.go                 # MODIFIED: Focusable embeds tea.Model,
                         #   FocusManager.Update() routing logic
focus_test.go            # MODIFIED: Update test helpers to implement
                         #   tea.Model, add routing tests
nav_test.go              # MODIFIED: Update focusScreen test helper
lifecycle_test.go        # MODIFIED: Update test helpers if needed
examples/basic/main.go   # MODIFIED: Update to new Focusable contract
```

**Structure Decision**: Same single package (`package nav`) at the
module root. This feature modifies existing files only. No new
packages or directories needed.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| Breaking Focusable interface | Spec requires tea.Model embedding; no new interface | Runtime type assertion defers constraint to runtime; compile-time is safer |
