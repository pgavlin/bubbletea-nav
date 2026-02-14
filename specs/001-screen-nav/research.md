# Research: Screen Navigation Library

**Branch**: `001-screen-nav` | **Date**: 2026-02-14

## R1: Bubble Tea Version Target

**Decision**: Target Bubble Tea v1 (stable v1.3.10).

**Rationale**: v2.0.0 is still in release candidate (v2.0.0-rc.2,
Nov 2024). A library targeting an unstable API risks breaking
changes. v1 has a proven, well-documented API. Migration to v2
can be a separate effort once v2 reaches stable.

**Alternatives considered**:
- Target v2-rc: Higher risk of API breakage; smaller user base.
- Dual v1/v2 support: Premature complexity; violates Principle VII.

## R2: Package Architecture

**Decision**: Single package (`package nav`) at the module root.

**Rationale**: The navigation stack and focus manager are closely
coupled (US4 requires combining them). Splitting into sub-packages
would force cross-package imports for a common use case. One package
keeps the import path simple and follows Go community conventions
for small-to-medium libraries.

**Alternatives considered**:
- `nav/stack` + `nav/focus` sub-packages: Adds import complexity
  without clear benefit at this scope.
- `nav` + `nav/debug` sub-package: Constitution VI suggests debug
  helpers SHOULD be separate. Deferred until debug helpers are
  implemented (not in scope for initial release).

## R3: Screen Interface Design

**Decision**: Define a `Screen` interface extending `tea.Model` with
no additional required methods. Lifecycle support via optional
interface assertion (`LifecycleScreen`).

**Rationale**: Any `tea.Model` can be a screen — this maximizes
compatibility with existing Bubble Tea components. Lifecycle events
are opt-in, keeping the barrier to entry low. This follows
Principle III (narrow interfaces) and Principle VII (simplicity).

**Alternatives considered**:
- Require lifecycle methods on all screens: Forces boilerplate on
  screens that don't need lifecycle events.
- Wrap screens in a struct with lifecycle callbacks: Adds a layer
  of indirection; less idiomatic.

## R4: Focus Manager Design

**Decision**: `FocusManager` is a value type (struct) that screens
embed or hold as a field. It processes `tea.Msg` and returns updated
state + `tea.Cmd`, following the Bubble Tea Update pattern.

**Rationale**: Making FocusManager a `tea.Model`-like value type
means it composes naturally with screens — the screen delegates
message handling to FocusManager in its own `Update`. This matches
the `textinput.Model` pattern from charmbracelet/bubbles.

**Alternatives considered**:
- FocusManager as a standalone `tea.Model`: Would require message
  forwarding and complicate screen composition.
- Focus logic built into Stack: Creates a God-object (violates
  Principle III).

## R5: Focusable Interface and Mouse Hit Testing

**Decision**: Define a `Focusable` interface with `Focus() tea.Cmd`,
`Blur()`, and `Focused() bool` methods (matching the bubbles
convention). For mouse hit testing, add a `Bounds() (x, y, w, h)`
method to an optional `Bounded` interface.

**Rationale**: Matching the existing bubbles Focus/Blur pattern
means existing components (textinput, textarea) are trivially
compatible. Mouse hit testing is a separate concern; components
that don't support mouse clicks don't need to implement `Bounded`.

**Alternatives considered**:
- Single interface with all methods: Forces mouse bounds on
  keyboard-only components.
- Zone-based hit testing (like lipgloss zones): More complex;
  couples to lipgloss; deferred.

## R6: Navigation Message Types

**Decision**: Use typed message structs for navigation commands:
`PushMsg`, `PopMsg`, `ReplaceMsg`. Lifecycle notifications use
`ScreenAppearedMsg` and `ScreenDisappearedMsg`.

**Rationale**: Typed messages are idiomatic Bubble Tea. They are
type-switchable in `Update`, inspectable for logging (Principle VI),
and composable with `tea.Batch`.

**Alternatives considered**:
- Method calls on Stack: Violates Principle II (state via Msg/Cmd
  only, no direct mutation).
- Generic `NavMsg` with action field: Less type-safe; harder to
  switch on.

## R7: Re-entrant Stack Modification Prevention

**Decision**: Stack operations during lifecycle event processing
are queued and executed after the current event completes.

**Rationale**: FR-013 requires this. The Stack's `Update` method
will track whether it is processing a lifecycle event. If a
navigation Msg arrives during lifecycle processing, it is appended
to a pending queue and replayed after the current operation.

**Alternatives considered**:
- Error/panic on re-entrant modification: Too harsh; hard to
  debug.
- Ignore re-entrant operations: Silent data loss; violates
  principle of least surprise.

## R8: Bubble Tea v1 API Details for Implementation

### tea.Model Interface (v1)
```
Init() tea.Cmd
Update(tea.Msg) (tea.Model, tea.Cmd)
View() string
```

### Key Detection (v1)
- Tab: `msg.Type == tea.KeyTab`
- Shift+Tab: `msg.Type == tea.KeyShiftTab`

### Mouse Events (v1)
- `tea.MouseMsg` with `X`, `Y`, `Action`, `Button` fields
- Action: `MouseActionPress`, `MouseActionRelease`, `MouseActionMotion`
- Enable via `tea.WithMouseCellMotion()` program option
- Coordinates: zero-based, (0,0) at upper-left

### Window Resize (v1)
- `tea.WindowSizeMsg{Width, Height}`
- Delivered automatically at startup and on resize

### Command Combinators
- `tea.Batch(cmds...)` — concurrent execution
- `tea.Sequence(cmds...)` — sequential execution
