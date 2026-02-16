# Data Model: Message-Based Lifecycle Notifications

**Feature**: 003-lifecycle-msg-routing
**Date**: 2026-02-15

## Entities

### ScreenAppearedMsg (unchanged)

A zero-field message type delivered to a screen's `Update` method when it becomes the active (topmost) screen.

```text
ScreenAppearedMsg {}
```

- **Delivered on**: push (to new screen), pop (to revealed screen), replace (to new screen)
- **Ordering**: Always delivered after `ScreenDisappearedMsg` to the outgoing screen
- **For push/replace**: Delivered after the new screen's `Init()` command is collected

### ScreenDisappearedMsg (unchanged)

A zero-field message type delivered to a screen's `Update` method when it loses active status.

```text
ScreenDisappearedMsg {}
```

- **Delivered on**: push (to old top screen), pop (to popped screen), replace (to old screen)
- **Ordering**: Always delivered before `ScreenAppearedMsg` to the incoming screen

### LifecycleScreen (removed)

Previously an optional interface for screens to receive lifecycle callbacks. Removed in this feature — all lifecycle notifications now flow through `Screen.Update`.

## State Transitions

### Push Operation

```text
Stack: [... A]  →  [... A, B]

1. A.Update(ScreenDisappearedMsg{})  → store updated A, collect cmd
2. B.Update(ScreenAppearedMsg{})     → store updated B, collect cmd
3. B.Init()                          → collect cmd
4. Return batched cmds
```

### Pop Operation

```text
Stack: [... A, B]  →  [... A]

1. B.Update(ScreenDisappearedMsg{})  → discard updated B, collect cmd
2. A.Update(ScreenAppearedMsg{})     → store updated A, collect cmd
3. Return batched cmds
```

### Replace Operation

```text
Stack: [... A]  →  [... B]

1. A.Update(ScreenDisappearedMsg{})  → discard updated A, collect cmd
2. B.Update(ScreenAppearedMsg{})     → store updated B, collect cmd
3. B.Init()                          → collect cmd
4. Return batched cmds
```

## Removed Entities

| Entity | Reason |
|--------|--------|
| `LifecycleScreen` interface | Replaced by message-based delivery through `Screen.Update` |
| `Stack.dispatchAppeared()` helper | Replaced by direct `screen.Update(ScreenAppearedMsg{})` call |
| `Stack.dispatchDisappeared()` helper | Replaced by direct `screen.Update(ScreenDisappearedMsg{})` call |
