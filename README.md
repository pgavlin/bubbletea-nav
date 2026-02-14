# bubbletea-nav

[![Go Reference](https://pkg.go.dev/badge/github.com/pgavlin/bubbletea-nav.svg)](https://pkg.go.dev/github.com/pgavlin/bubbletea-nav)

Screen navigation and focus management for [Bubble Tea](https://github.com/charmbracelet/bubbletea) terminal UI applications.

## Features

**Stack** manages an ordered stack of screens with push, pop, and replace operations. It implements `tea.Model` and can be passed directly to `tea.NewProgram`.

**FocusManager** manages sequential focus order within a screen. Tab/Shift+Tab cycle focus between components, and mouse clicks on bounded components move focus automatically. All other messages are routed to the focused component.

**Lifecycle hooks** let screens react to visibility changes. Implement the optional `LifecycleScreen` interface to receive `Appeared` and `Disappeared` callbacks when screens are pushed over, popped, or replaced.

## Install

```
go get github.com/pgavlin/bubbletea-nav
```

## Quick start

### Navigation stack

```go
// Create a stack with a root screen.
stack := nav.NewStack(homeScreen{})
p := tea.NewProgram(stack)
p.Run()
```

Screens implement `nav.Screen` (like `tea.Model`, but `Update` returns `nav.Screen`):

```go
func (s myScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyEnter:
            return s, nav.Push(detailScreen{})
        case tea.KeyEscape:
            return s, nav.Pop()
        }
    }
    return s, nil
}
```

### Focus management

Create focusable components by implementing `nav.Focusable` (`tea.Model` + `Focus`/`Blur`), then delegate messages:

```go
type formScreen struct {
    focus  nav.FocusManager
    fields []*textField
}

func newFormScreen() formScreen {
    fields := []*textField{
        newTextField("Name"),
        newTextField("Email"),
    }
    return formScreen{
        focus:  nav.NewFocusManager(fields[0], fields[1]),
        fields: fields,
    }
}

func (s formScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
    // Handle screen-level keys first.
    if msg, ok := msg.(tea.KeyMsg); ok {
        if msg.Type == tea.KeyEscape {
            return s, nav.Pop()
        }
    }
    // Delegate everything else to FocusManager.
    var cmd tea.Cmd
    s.focus, cmd = s.focus.Update(msg)
    return s, cmd
}
```

## Examples

- **[basic](examples/basic/)** — Two-screen push/pop navigation
- **[focus-routing](examples/focus-routing/)** — FocusManager with a three-field form
- **[focus-nav](examples/focus-nav/)** — Stack navigation combined with FocusManager

Run an example:

```
go run ./examples/basic
```

## API overview

### Navigation

| Type/Function | Description |
|---|---|
| `Screen` | Interface: `Init`, `Update` (returns `Screen`), `View` |
| `LifecycleScreen` | Optional: `Appeared`, `Disappeared` callbacks |
| `Stack` | Navigation stack implementing `tea.Model` |
| `NewStack(root)` | Create a stack with a root screen |
| `Push(screen)` | Command to push a screen |
| `Pop()` | Command to pop the top screen |
| `Replace(screen)` | Command to replace the top screen |

### Focus

| Type/Function | Description |
|---|---|
| `Focusable` | Interface: `tea.Model` + `Focus`, `Blur` |
| `Bounded` | Optional: `Bounds()` for mouse-click targeting |
| `FocusManager` | Value type managing focus order |
| `NewFocusManager(items...)` | Create with focusable items (first gets focus) |
| `FocusManager.Update(msg)` | Route messages; Tab/Shift+Tab cycle focus |
| `FocusManager.FocusedIndex()` | Index of focused item (-1 if none) |
| `FocusManager.FocusIndex(i)` | Set focus to index programmatically |
| `FocusManager.SetItems(items...)` | Replace the items list |
| `FocusChangedMsg` | Emitted when focus moves between items |

---

This package was written using [Claude Code](https://claude.ai/claude-code) and [Spec Kit](https://speckit.org/).
