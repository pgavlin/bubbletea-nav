# Research: Unify on tea.Model

**Feature**: 004-unify-tea-model
**Date**: 2026-02-15

## R1: Replacing Screen with tea.Model in the Stack

**Decision**: Change `Stack.screens` from `[]Screen` to `[]tea.Model`. Update all navigation message types and helper functions to accept `tea.Model`.

**Rationale**: The `Screen` interface is identical to `tea.Model` except `Update` returns `(Screen, tea.Cmd)` instead of `(tea.Model, tea.Cmd)`. This forces consumers to use a non-standard return type, preventing reuse of existing Bubble Tea components. Since `Stack.Update` already returns `(tea.Model, tea.Cmd)`, the stack internally must handle `tea.Model` anyway — the `Screen` interface adds no type safety benefit.

**Alternatives considered**:
- Keep `Screen` and add adapters: Rejected — adds complexity, every consumer pays the cost.
- Make `Screen` an alias for `tea.Model`: Not possible in Go — `tea.Model` is an interface, and type aliases for interfaces don't change the method set.

## R2: Replacing Focusable with tea.Model + Focus/Blur Messages

**Decision**: Remove the `Focusable` interface. Change `FocusManager.items` from `[]Focusable` to `[]tea.Model`. Deliver `FocusMsg` and `BlurMsg` through each item's `Update` method. Store the returned model and collect commands.

**Rationale**: This applies the same message-based pattern that was successfully used for lifecycle events (feature 003). It eliminates a custom interface that prevents standard Bubble Tea components from participating in focus management.

**Alternatives considered**:
- Keep `Focusable` alongside `tea.Model` support: Rejected — maintaining two code paths adds complexity without benefit.
- Use a single `FocusChangeMsg{Focused bool}`: Rejected — separate `FocusMsg` and `BlurMsg` types are clearer for type-switching in `Update` and match the separate `ScreenAppearedMsg`/`ScreenDisappearedMsg` pattern.

## R3: Signature Changes for NewFocusManager and SetItems

**Decision**: `NewFocusManager` and `SetItems` must return `tea.Cmd` in addition to `FocusManager`, because delivering `FocusMsg`/`BlurMsg` via `Update` produces commands that must be surfaced.

**Rationale**: Currently `Focus()` returns `tea.Cmd` and `Blur()` returns nothing. With message-based delivery, both `FocusMsg` and `BlurMsg` go through `Update` which returns `(tea.Model, tea.Cmd)`. The commands from these calls must be returned to the caller. This is a breaking change to the `NewFocusManager` and `SetItems` signatures:
- `NewFocusManager(items ...tea.Model) (FocusManager, tea.Cmd)` — returns cmd from initial FocusMsg to first item
- `SetItems(items ...tea.Model) (FocusManager, tea.Cmd)` — returns cmds from BlurMsg to old focused item and FocusMsg to new first item

**Alternatives considered**:
- Queue the commands internally: Rejected — Bubble Tea's model is explicit command return, not internal queuing.
- Ignore commands from Focus/Blur: Rejected — violates FR-007 and loses functionality (e.g., a text input that starts a cursor blink timer on focus).

## R4: FocusIndex Signature Change

**Decision**: `FocusIndex(int)` currently returns `(FocusManager, tea.Cmd)` — this signature stays the same. The internal implementation changes from calling `Blur()`/`Focus()` methods to delivering `BlurMsg`/`FocusMsg` via `Update`, but the external signature already returns commands.

**Rationale**: No breaking change needed here since the return type already includes `tea.Cmd`.

## R5: Bounded Interface Unchanged

**Decision**: The `Bounded` interface remains as-is. It is checked via type assertion on `tea.Model` items (same as before on `Focusable` items).

**Rationale**: `Bounded` is orthogonal to focus management. It provides optional mouse-click targeting and doesn't depend on the `Focusable` interface. Type-asserting `tea.Model` to `Bounded` works identically to type-asserting `Focusable` to `Bounded`.

## R6: routeMessage Type Assertion

**Decision**: `routeMessage` currently type-asserts the `tea.Model` returned by `Update` back to `Focusable`. After removing `Focusable`, no type assertion is needed — the returned `tea.Model` is stored directly.

**Rationale**: Since items are now `[]tea.Model`, the `Update` method returns `tea.Model`, and we store `tea.Model` — the types align naturally without assertion.
