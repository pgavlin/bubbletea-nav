# Data Model: Focus Message Routing

**Branch**: `002-focus-message-routing` | **Date**: 2026-02-14

## Modified Interfaces

### Focusable (modified — embeds tea.Model)

The Focusable interface is expanded to embed `tea.Model`. All
focusable components must now implement `Init`, `Update`, and `View`
in addition to the existing focus methods. This is a breaking change.

```
Focusable interface {
    tea.Model       // Init, Update, View
    Focus() Cmd
    Blur()
    Focused() bool
}
```

The `tea.Model.Update` method is used by FocusManager for message
routing. The returned `tea.Model` is cast back to `Focusable` and
stored in the items slice to preserve updated state.

### Bounded (unchanged)

```
Bounded interface {
    Bounds() (x, y, width, height int)
}
```

No changes.

## Modified Structs

### FocusManager (modified behavior, same fields)

```
FocusManager struct {
    items       []Focusable  // ordered focus list (unchanged type)
    focusIndex  int          // -1 = no focus (unchanged)
}
```

**Fields**: No new fields. The struct layout is identical.

**Modified behavior in Update**:

1. Tab/Shift+Tab → focus cycling only (no routing). Unchanged.
2. Mouse press on bounded item:
   - Different item: change focus, then route click via
     `item.Update(msg)`. Store returned Focusable. Batch commands.
   - Same item: route click via `item.Update(msg)`. Store returned
     Focusable. Modified (was no-op).
   - No hit: fall through to default routing. New.
3. All other messages → route to focused item via
   `item.Update(msg)`. Store returned Focusable. New.
4. No focus (index -1) → no routing. New guard.

**State transitions** (additions):

- Routed message: `item.Update(msg)` returns `(tea.Model, tea.Cmd)`.
  The `tea.Model` is asserted to `Focusable` and stored back in the
  slice at the focused index. The command is returned (or batched
  with focus-change commands).

## Message Types

No new message types. Existing `FocusChangedMsg` is unchanged.

## Relationships

```
FocusManager 1──* Focusable    (unchanged)
Focusable ───── tea.Model      (NEW: Focusable embeds tea.Model)
Focusable *──0..1 Bounded      (unchanged)
```

## Validation Rules

Existing rules unchanged. New rules:

- After routing, the `tea.Model` returned by `Update` MUST be
  assertable to `Focusable`. If it is not, this is a programming
  error in the component.
- `tea.Model.Update` on the focused item MUST NOT be called with
  Tab/Shift+Tab messages (these are consumed by FocusManager).
