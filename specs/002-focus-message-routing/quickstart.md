# Quickstart: Focus Message Routing

**Branch**: `002-focus-message-routing` | **Date**: 2026-02-14

## Before (manual routing boilerplate)

Without message routing, developers must manually dispatch messages
to the focused component:

```go
func (s formScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
    // Step 1: Delegate to focus manager for Tab/Shift+Tab
    var cmd tea.Cmd
    s.focus, cmd = s.focus.Update(msg)
    if cmd != nil {
        return s, cmd
    }

    // Step 2: Manually route to the focused component
    switch s.focus.FocusedIndex() {
    case 0:
        s.name, cmd = s.name.Update(msg)
    case 1:
        s.email, cmd = s.email.Update(msg)
    case 2:
        s.phone, cmd = s.phone.Update(msg)
    }
    return s, cmd
}
```

This switch statement grows with every component and must be updated
whenever components are added, removed, or reordered.

## After (automatic routing)

With message routing, the FocusManager handles everything:

```go
func (s formScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
    // One call handles focus cycling AND message routing
    var cmd tea.Cmd
    s.focus, cmd = s.focus.Update(msg)
    return s, cmd
}
```

The FocusManager automatically forwards non-focus messages to the
focused component via its `tea.Model.Update` method.

## The Focusable Interface

Focusable now embeds `tea.Model`. All focusable components must
implement `Init`, `Update`, and `View`:

```go
// nav.Focusable interface (updated):
//     tea.Model (Init, Update, View)
//     Focus() tea.Cmd
//     Blur()
//     Focused() bool

type myInput struct {
    value   string
    focused bool
}

// tea.Model methods (required)
func (m myInput) Init() tea.Cmd { return nil }

func (m myInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        m.value += msg.String()
    }
    return m, nil
}

func (m myInput) View() string {
    if m.focused {
        return "> " + m.value
    }
    return "  " + m.value
}

// Focusable methods (required)
func (m *myInput) Focus() tea.Cmd { m.focused = true; return nil }
func (m *myInput) Blur()          { m.focused = false }
func (m *myInput) Focused() bool  { return m.focused }
```

## Wrapping Existing Bubbles Components

For existing charmbracelet/bubbles components (e.g., `textinput`),
write a thin adapter:

```go
type textInputAdapter struct {
    textinput.Model
}

// tea.Model methods delegate to embedded textinput.Model
// (Init and View are promoted automatically)

func (a textInputAdapter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    a.Model, cmd = a.Model.Update(msg)
    return a, cmd
}

// Focusable methods
func (a *textInputAdapter) Focus() tea.Cmd { return a.Model.Focus() }
func (a *textInputAdapter) Blur()          { a.Model.Blur() }
func (a *textInputAdapter) Focused() bool  { return a.Model.Focused() }
```

Then use it with the FocusManager:

```go
name := &textInputAdapter{Model: textinput.New()}
email := &textInputAdapter{Model: textinput.New()}

fm := nav.NewFocusManager(name, email)
```

## Key Patterns

| Scenario | Behavior |
|----------|----------|
| Tab/Shift+Tab pressed | Focus cycles. Message NOT forwarded to component. |
| Regular key pressed | Forwarded to focused component via Update. |
| Mouse click on bounded component | Focus moves, then click forwarded via Update. |
| Mouse click, focus doesn't change | Click forwarded to focused component via Update. |
| No component focused (index -1) | Message silently ignored. |
