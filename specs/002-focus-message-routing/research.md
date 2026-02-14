# Research: Focus Message Routing

**Branch**: `002-focus-message-routing` | **Date**: 2026-02-14

## R1: Routing via tea.Model (No New Interface)

**Decision**: Embed `tea.Model` into the existing `Focusable`
interface. The FocusManager routes messages by calling the standard
`tea.Model.Update(tea.Msg) (tea.Model, tea.Cmd)` method on the
focused item.

**Rationale**: The spec requires that all Focusable items implement
`tea.Model` and that no new interface be introduced. Embedding
`tea.Model` into `Focusable` enforces this constraint at compile
time. The FocusManager can call `item.Update(msg)` directly without
runtime type assertions for routing.

The updated Focusable interface:

```go
type Focusable interface {
    tea.Model
    Focus() tea.Cmd
    Blur()
    Focused() bool
}
```

This is a breaking change to the existing `Focusable` interface.
The spec explicitly states backward compatibility is not required.

**Alternatives considered**:

- Runtime type assertion to `tea.Model`: Keeps Focusable unchanged
  but defers the constraint to runtime. Since the spec says ALL
  items must implement tea.Model, compile-time enforcement is
  preferable.
- New combined interface (e.g., `FocusableModel`): Spec says no new
  interface. Embedding tea.Model into Focusable achieves the same
  result without a new type.
- Pointer-mutation approach (`Update(tea.Msg) tea.Cmd`): Would
  require a new interface. Rejected per spec clarification.

## R2: Storing Updated State After Routing

**Decision**: After calling `item.Update(msg)`, the returned
`tea.Model` is cast back to `Focusable` and stored in the
`[]Focusable` slice, replacing the previous value.

**Rationale**: `tea.Model.Update` returns `(tea.Model, tea.Cmd)`.
The returned `tea.Model` reflects the component's updated state.
To satisfy FR-003 (retain updated state), the FocusManager must
store this value back. Since the slice stores `Focusable` values,
the returned `tea.Model` must be asserted to `Focusable`:

```go
updated, cmd := fm.items[fm.focusIndex].Update(msg)
fm.items[fm.focusIndex] = updated.(Focusable)
```

This assertion is safe because:
1. The component implemented `Focusable` before the call.
2. Well-behaved components return the same concrete type from
   `Update`.
3. The spec requires all items to implement both interfaces.

**Alternatives considered**:

- Not storing the updated value: Violates FR-003.
- Storing as `tea.Model` (changing the slice type): Would require
  type assertions elsewhere (Focus, Blur, Focused calls). More
  disruptive than a single assertion at the routing site.

## R3: Message Routing Flow

**Decision**: The updated `FocusManager.Update` follows this flow:

1. If items list is empty → return immediately.
2. Match on message type:
   a. `tea.KeyMsg` with `KeyTab` or `KeyShiftTab` → handle focus
      cycling. Do NOT route to any component. Return.
   b. `tea.MouseMsg` with press action → check bounded items for
      hit. If hit moves focus, move focus first, then route the
      click to the newly focused item. If already focused, just
      route. Return.
   c. All other messages (including non-press mouse events) →
      route to the focused item.
3. If focus index is -1 (no focus) → do not route.

**Rationale**: This satisfies:
- FR-001: Non-focus messages routed to focused component via
  tea.Model.Update.
- FR-004: Tab/Shift+Tab consumed, never forwarded.
- FR-008: Mouse clicks that move focus are forwarded after focus
  change.
- FR-005: No routing when no component has focus.

## R4: Command Combining

**Decision**: When a focus change and message routing both produce
commands (as in the mouse click case), combine them with `tea.Batch`.

**Rationale**: FR-006 requires all commands to be combined. The
FocusChangedMsg command from focus cycling is batched with the
command from the routed message. `tea.Batch` is the idiomatic
Bubble Tea combinator and handles nil commands gracefully.

## R5: Breaking Changes to Existing Code

**Decision**: This feature introduces the following breaking changes:

1. `Focusable` interface now embeds `tea.Model` — existing types
   that implement `Focusable` must also implement `Init`, `Update`,
   and `View`.
2. `FocusManager.Update` now routes messages — existing code that
   relied on the FocusManager NOT routing (and handled routing
   manually) may see double-processing if the manual routing is
   not removed.
3. Existing test helpers (`mockFocusable`, `mockBoundedFocusable`)
   must be updated to implement `tea.Model`.

**Rationale**: The spec explicitly states backward compatibility is
not required. These breaking changes simplify the API surface and
eliminate boilerplate for all consumers.
