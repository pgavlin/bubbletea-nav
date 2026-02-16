# Implementation Plan: Message-Based Lifecycle Notifications

**Branch**: `003-lifecycle-msg-routing` | **Date**: 2026-02-15 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/003-lifecycle-msg-routing/spec.md`

## Summary

Replace the `LifecycleScreen` interface (direct method calls) with message-based lifecycle notifications. The Stack will send `ScreenAppearedMsg` and `ScreenDisappearedMsg` through each screen's standard `Update` method instead of calling `Appeared()`/`Disappeared()` directly. This lets screens update their own state in response to lifecycle events using the same pattern they use for all other messages.

## Technical Context

**Language/Version**: Go 1.25+ (per existing `go.mod`)
**Primary Dependencies**: `github.com/charmbracelet/bubbletea` v1.3.10
**Storage**: N/A (in-memory state only)
**Testing**: `go test ./...` (standard library testing package)
**Target Platform**: Cross-platform library (any OS supported by Go + Bubble Tea)
**Project Type**: Single Go library package
**Performance Goals**: N/A (in-process synchronous message delivery)
**Constraints**: Must maintain existing lifecycle ordering semantics (disappear before appear)
**Scale/Scope**: 3 source files modified, ~30 lines changed in nav.go, lifecycle.go simplified, tests rewritten

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Library-First | PASS | All changes remain in the `nav` package, importable via `go get` |
| II. Bubble Tea Idioms | PASS | Moves lifecycle handling INTO the `tea.Msg`/`tea.Cmd` flow — better alignment than current interface approach |
| III. Composable API | PASS | Removes one interface (`LifecycleScreen`), reduces API surface |
| IV. Test-Required | PASS | All existing lifecycle tests will be rewritten for message-based approach; new state-persistence test added |
| V. Backward Compatibility | NOTE | Removing `LifecycleScreen` is a breaking change. Pre-v1.0, acceptable without deprecation cycle per semver. |
| VI. Observability | PASS | `ScreenAppearedMsg`/`ScreenDisappearedMsg` are `tea.Msg` types — already inspectable and loggable |
| VII. Simplicity | PASS | Removes an interface and two helper methods; uses existing message-passing pattern instead |

**Gate result**: PASS. Breaking change is acceptable pre-v1.0.

## Project Structure

### Documentation (this feature)

```text
specs/003-lifecycle-msg-routing/
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
├── nav.go               # Stack — modify handleNav to send messages via Update
├── lifecycle.go         # Remove LifecycleScreen interface, keep message types
├── lifecycle_test.go    # Rewrite tests for message-based approach
├── nav_test.go          # No changes expected (mockScreen doesn't use lifecycle)
├── focus.go             # No changes
├── focus_test.go        # No changes
└── examples/
    ├── basic/           # No changes (doesn't use lifecycle)
    ├── focus-routing/   # No changes
    └── focus-nav/       # No changes
```

**Structure Decision**: Existing single-package library layout. No new files or directories needed beyond spec artifacts.
