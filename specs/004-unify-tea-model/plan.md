# Implementation Plan: Unify on tea.Model

**Branch**: `004-unify-tea-model` | **Date**: 2026-02-15 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/004-unify-tea-model/spec.md`

## Summary

Remove the `Screen` and `Focusable` custom interfaces in favor of the standard `tea.Model` interface. The navigation stack (`Stack`) will store `[]tea.Model` instead of `[]Screen`. The focus manager (`FocusManager`) will store `[]tea.Model` instead of `[]Focusable` and deliver `FocusMsg`/`BlurMsg` through each item's `Update` method instead of calling `Focus()`/`Blur()` directly. This unifies the entire library on the standard Bubble Tea interface, making all existing `tea.Model` implementations directly compatible.

## Technical Context

**Language/Version**: Go 1.25+ (per existing `go.mod`)
**Primary Dependencies**: `github.com/charmbracelet/bubbletea` v1.3.10
**Storage**: N/A (in-memory state only)
**Testing**: `go test ./...` (standard library testing package)
**Target Platform**: Cross-platform library (any OS supported by Go + Bubble Tea)
**Project Type**: Single Go library package
**Performance Goals**: N/A (in-process synchronous message delivery)
**Constraints**: Must maintain blur-before-focus ordering semantics
**Scale/Scope**: 8 source files modified, 3 examples updated, 2 interfaces removed, 2 message types added

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Library-First | PASS | All changes remain in the `nav` package, importable via `go get` |
| II. Bubble Tea Idioms | PASS | Moves entirely to `tea.Model` — maximum alignment with Bubble Tea conventions |
| III. Composable API | PASS | Removes two interfaces (`Screen`, `Focusable`), reduces API surface by 2 types and 2 methods |
| IV. Test-Required | PASS | All existing tests rewritten; new tests added for message-based focus delivery and state persistence |
| V. Backward Compatibility | NOTE | Removing `Screen` and `Focusable` are breaking changes. `NewFocusManager` and `SetItems` signatures change to return `tea.Cmd`. Pre-v1.0, acceptable without deprecation cycle per semver. |
| VI. Observability | PASS | `FocusMsg`/`BlurMsg` are `tea.Msg` types — inspectable and loggable. `Stack.String()` and `FocusManager.String()` unchanged. |
| VII. Simplicity | PASS | Removes two interfaces, replaces with message-passing pattern already established by lifecycle feature (003) |

**Gate result**: PASS. Breaking changes acceptable pre-v1.0.

**Post-design re-check**: PASS. No new complexity introduced. `FocusMsg`/`BlurMsg` follow the same empty-struct message pattern as `ScreenAppearedMsg`/`ScreenDisappearedMsg`.

## Project Structure

### Documentation (this feature)

```text
specs/004-unify-tea-model/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
│   └── api.go.txt       # Updated public API contract
└── tasks.md             # Phase 2 output (created by /speckit.tasks)
```

### Source Code (repository root)

```text
.
├── nav.go               # Remove Screen interface; change Stack to use []tea.Model
├── focus.go             # Remove Focusable interface; add FocusMsg/BlurMsg; update FocusManager to use []tea.Model
├── lifecycle.go         # No changes (already message-based)
├── nav_test.go          # Update mockScreen, recordingScreen, focusScreen to return tea.Model
├── lifecycle_test.go    # Update lifecycleScreen and helper types to return tea.Model
├── focus_test.go        # Update mockFocusable to handle FocusMsg/BlurMsg in Update; remove Focus()/Blur()
├── doc.go               # Update package doc examples
└── examples/
    ├── basic/main.go        # Update Screen return types to tea.Model
    ├── focus-routing/main.go # Replace Focus()/Blur() with FocusMsg/BlurMsg; update NewFocusManager call
    └── focus-nav/main.go     # Both: tea.Model returns + FocusMsg/BlurMsg + NewFocusManager signature
```

**Structure Decision**: Existing single-package library layout. No new files needed. `FocusMsg` and `BlurMsg` added to `focus.go` alongside existing focus types.
