# Data Model: Unify on tea.Model

**Feature**: 004-unify-tea-model
**Date**: 2026-02-15

## New Types

### FocusMsg

A message delivered to a model's `Update` when it gains keyboard focus.

| Field | Type | Description |
|-------|------|-------------|
| (none) | — | Empty struct, identity is the type itself |

### BlurMsg

A message delivered to a model's `Update` when it loses keyboard focus.

| Field | Type | Description |
|-------|------|-------------|
| (none) | — | Empty struct, identity is the type itself |

## Modified Types

### Stack

| Field | Before | After |
|-------|--------|-------|
| screens | `[]Screen` | `[]tea.Model` |

### PushMsg

| Field | Before | After |
|-------|--------|-------|
| Screen | `Screen` | `tea.Model` |

### ReplaceMsg

| Field | Before | After |
|-------|--------|-------|
| Screen | `Screen` | `tea.Model` |

### FocusManager

| Field | Before | After |
|-------|--------|-------|
| items | `[]Focusable` | `[]tea.Model` |

## Modified Functions

| Function | Before Signature | After Signature |
|----------|-----------------|-----------------|
| NewStack | `NewStack(root Screen) Stack` | `NewStack(root tea.Model) Stack` |
| Push | `Push(screen Screen) tea.Cmd` | `Push(screen tea.Model) tea.Cmd` |
| Replace | `Replace(screen Screen) tea.Cmd` | `Replace(screen tea.Model) tea.Cmd` |
| NewFocusManager | `NewFocusManager(items ...Focusable) FocusManager` | `NewFocusManager(items ...tea.Model) (FocusManager, tea.Cmd)` |
| SetItems | `SetItems(items ...Focusable) FocusManager` | `SetItems(items ...tea.Model) (FocusManager, tea.Cmd)` |

## Removed Types

| Type | Reason |
|------|--------|
| Screen | Replaced by `tea.Model` |
| Focusable | Replaced by `tea.Model` + `FocusMsg`/`BlurMsg` messages |

## Unchanged Types

| Type | Notes |
|------|-------|
| PopMsg | No Screen field to change |
| Bounded | Orthogonal optional interface |
| FocusChangedMsg | Unchanged behavior |
| ScreenAppearedMsg | Unchanged |
| ScreenDisappearedMsg | Unchanged |

## State Transitions

### Focus Transition (blur-before-focus)

```
Old focused item: Update(BlurMsg{}) → store returned model, collect cmd
New focused item: Update(FocusMsg{}) → store returned model, collect cmd
Emit FocusChangedMsg command
Return batched commands
```
