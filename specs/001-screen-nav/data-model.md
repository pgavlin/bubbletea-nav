# Data Model: Screen Navigation Library

**Branch**: `001-screen-nav` | **Date**: 2026-02-14

## Interfaces

### Screen

The minimal contract for a navigable screen. Identical to
`tea.Model` — any existing Bubble Tea model can be used as a screen.

```
Screen interface {
    Init() Cmd
    Update(Msg) (Screen, Cmd)
    View() string
}
```

Note: `Update` returns `Screen` (not `tea.Model`) to preserve type
safety within the navigation stack.

### LifecycleScreen (optional)

Screens that need lifecycle notifications implement this interface
in addition to `Screen`. Detected via interface assertion.

```
LifecycleScreen interface {
    Screen
    Appeared() Cmd
    Disappeared()
}
```

- `Appeared()` is called when the screen becomes the active screen
  (after push, after a pop reveals it, after replace). Returns a
  `Cmd` for async initialization (e.g., start a timer).
- `Disappeared()` is called when the screen loses active status
  (another screen pushed, this screen popped, this screen replaced).
  Returns nothing — cleanup is synchronous.

### Focusable

The contract for a component that can receive and lose focus.
Compatible with the `charmbracelet/bubbles` Focus/Blur convention.

```
Focusable interface {
    Focus() Cmd
    Blur()
    Focused() bool
}
```

### Bounded (optional)

Components that support mouse-click focus targeting implement this
in addition to `Focusable`. Detected via interface assertion.

```
Bounded interface {
    Bounds() (x, y, width, height int)
}
```

Returns the component's bounding rectangle in terminal coordinates
(zero-based, top-left origin). Used by `FocusManager` to determine
which component a mouse click targets.

## Structs

### Stack

The navigation stack. Implements `tea.Model`. Manages an ordered
list of screens and routes messages to the active (top) screen.

```
Stack struct {
    screens        []Screen   // ordered, last = active
    pendingOps     []Msg      // queued nav ops during lifecycle
    inLifecycle    bool       // re-entrancy guard
}
```

**Fields**:
- `screens`: Slice of Screen values. The last element is the active
  screen. Minimum length 1 (root screen).
- `pendingOps`: Navigation messages received during lifecycle event
  processing. Replayed after lifecycle completes.
- `inLifecycle`: True while processing Appeared/Disappeared calls.
  Prevents re-entrant stack modification.

**State transitions**:
- `Push(screen)`: Append screen to `screens`. Old top gets
  Disappeared. New top gets Appeared.
- `Pop()`: Remove last screen from `screens`. Popped screen gets
  Disappeared. New top gets Appeared. No-op if len == 1.
- `Replace(screen)`: Replace last element. Old top gets
  Disappeared. New top gets Appeared.

### FocusManager

Manages sequential focus order within a screen. Value type (not a
pointer receiver). Screens hold it as a field and delegate message
handling.

```
FocusManager struct {
    items       []Focusable  // ordered focus list
    focusIndex  int          // -1 = no focus, 0..len-1 = focused
}
```

**Fields**:
- `items`: Ordered slice of focusable components. Set by the screen
  during update phase. Static per render cycle.
- `focusIndex`: Index of the currently focused item. -1 when no
  items are focusable or the list is empty.

**State transitions**:
- `Tab`: `focusIndex = (focusIndex + 1) % len(items)`, skipping
  non-focusable items (if mixed with wrapper types).
- `Shift+Tab`: `focusIndex = (focusIndex - 1 + len(items)) %
  len(items)`, skipping non-focusable items.
- `Mouse click at (x, y)`: For each item implementing `Bounded`,
  check if (x, y) is within bounds. If match, set `focusIndex` to
  that item.
- `Empty list`: `focusIndex` = -1. Tab/Shift+Tab are no-ops.

## Message Types

### Navigation Messages (input to Stack)

```
PushMsg struct {
    Screen Screen   // the screen to push
}

PopMsg struct {}    // no fields needed

ReplaceMsg struct {
    Screen Screen   // the replacement screen
}
```

These are sent by screens as `tea.Cmd` return values to request
navigation. The Stack intercepts them in its `Update`.

### Lifecycle Messages (output from Stack to Screen)

```
ScreenAppearedMsg struct {}

ScreenDisappearedMsg struct {}
```

Delivered to screens via their `Appeared()` / `Disappeared()`
methods (LifecycleScreen interface), not via `Update`. This avoids
polluting the screen's message handling with lifecycle plumbing.

### Focus Messages (output from FocusManager)

```
FocusChangedMsg struct {
    Previous int    // index of previously focused item (-1 if none)
    Current  int    // index of newly focused item (-1 if none)
}
```

Emitted by FocusManager when focus changes, allowing the screen to
react (e.g., scroll to focused component).

## Relationships

```
Stack 1──* Screen          (stack contains 1+ screens)
Screen 1──0..1 FocusManager (screen optionally holds a FocusManager)
FocusManager 1──* Focusable (manager holds 0+ focusable items)
Focusable *──0..1 Bounded   (focusable optionally has bounds)
Screen *──0..1 LifecycleScreen (screen optionally supports lifecycle)
```

## Validation Rules

- Stack MUST contain at least 1 screen at all times.
- FocusManager's `focusIndex` MUST be -1 or within
  `[0, len(items))`.
- `Bounded.Bounds()` MUST return non-negative width and height.
- `PushMsg.Screen` MUST NOT be nil.
- `ReplaceMsg.Screen` MUST NOT be nil.
