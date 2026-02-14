# Quickstart: Screen Navigation Library

**Branch**: `001-screen-nav` | **Date**: 2026-02-14

## Installation

```bash
go get github.com/<module-path>
```

## Minimal Example: Two-Screen Push/Pop

This example creates a home screen and a detail screen. Pressing
Enter on the home screen pushes the detail screen. Pressing Escape
on the detail screen pops back to the home screen.

```go
package main

import (
    "fmt"
    tea "github.com/charmbracelet/bubbletea"
    nav "<module-path>"
)

// -- Home Screen --

type homeScreen struct {
    selected int
}

func (s homeScreen) Init() tea.Cmd { return nil }

func (s homeScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyEnter:
            // Push a detail screen when Enter is pressed
            detail := detailScreen{item: s.selected}
            return s, nav.Push(detail)
        case tea.KeyUp:
            if s.selected > 0 {
                s.selected--
            }
        case tea.KeyDown:
            if s.selected < 2 {
                s.selected++
            }
        }
    }
    return s, nil
}

func (s homeScreen) View() string {
    items := []string{"Item A", "Item B", "Item C"}
    out := "Home Screen\n\n"
    for i, item := range items {
        cursor := "  "
        if i == s.selected {
            cursor = "> "
        }
        out += cursor + item + "\n"
    }
    out += "\nPress Enter to view, q to quit"
    return out
}

// -- Detail Screen --

type detailScreen struct {
    item int
}

func (s detailScreen) Init() tea.Cmd { return nil }

func (s detailScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyEscape:
            // Pop back to home screen
            return s, nav.Pop()
        }
    }
    return s, nil
}

func (s detailScreen) View() string {
    return fmt.Sprintf(
        "Detail Screen\n\nViewing item %d\n\nPress Esc to go back",
        s.item,
    )
}

// -- Main --

func main() {
    stack := nav.NewStack(homeScreen{})
    p := tea.NewProgram(stack)
    if _, err := p.Run(); err != nil {
        fmt.Println("Error:", err)
    }
}
```

## Adding Focus Management

This extends a screen with multiple focusable text inputs.

```go
type formScreen struct {
    focus  nav.FocusManager
    name   textinput.Model
    email  textinput.Model
}

func newFormScreen() formScreen {
    name := textinput.New()
    name.Placeholder = "Name"
    email := textinput.New()
    email.Placeholder = "Email"

    return formScreen{
        focus: nav.NewFocusManager(&name, &email),
        name:  name,
        email: email,
    }
}

func (s formScreen) Init() tea.Cmd { return nil }

func (s formScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
    // Delegate to focus manager first (handles Tab/Shift+Tab)
    var cmd tea.Cmd
    s.focus, cmd = s.focus.Update(msg)
    if cmd != nil {
        return s, cmd
    }

    // Update the focused input
    switch idx := s.focus.FocusedIndex(); idx {
    case 0:
        s.name, cmd = s.name.Update(msg)
    case 1:
        s.email, cmd = s.email.Update(msg)
    }
    return s, cmd
}

func (s formScreen) View() string {
    return fmt.Sprintf(
        "Form Screen\n\n%s\n%s\n\nTab to switch fields, Esc to go back",
        s.name.View(),
        s.email.View(),
    )
}
```

## Lifecycle Events

Implement `LifecycleScreen` for setup/teardown notifications.

```go
type dataScreen struct {
    data    string
    loading bool
}

func (s dataScreen) Appeared() tea.Cmd {
    // Start loading data when screen becomes visible
    return func() tea.Msg {
        return dataLoadedMsg{data: "fetched data"}
    }
}

func (s dataScreen) Disappeared() {
    // Cleanup when screen is no longer visible
}
```

## Key Patterns

| Action | How |
|--------|-----|
| Push a screen | Return `nav.Push(screen)` from Update |
| Pop current screen | Return `nav.Pop()` from Update |
| Replace current screen | Return `nav.Replace(screen)` from Update |
| Cycle focus forward | FocusManager handles Tab automatically |
| Cycle focus backward | FocusManager handles Shift+Tab automatically |
| Click to focus | FocusManager handles mouse clicks if items implement Bounded |
| Lifecycle setup | Implement `Appeared() tea.Cmd` on your screen |
| Lifecycle cleanup | Implement `Disappeared()` on your screen |
