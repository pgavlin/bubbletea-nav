# Quickstart: Unify on tea.Model

**Feature**: 004-unify-tea-model
**Date**: 2026-02-15

## Before / After Comparison

### Navigation Screen (Before)

```go
type myScreen struct{ data string }

func (s myScreen) Init() tea.Cmd { return nil }
func (s myScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
    return s, nil
}
func (s myScreen) View() string { return s.data }

stack := nav.NewStack(myScreen{data: "hello"})
```

### Navigation Screen (After)

```go
type myScreen struct{ data string }

func (s myScreen) Init() tea.Cmd { return nil }
func (s myScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    return s, nil
}
func (s myScreen) View() string { return s.data }

stack := nav.NewStack(myScreen{data: "hello"})
```

The only change: `Update` returns `(tea.Model, tea.Cmd)` instead of `(nav.Screen, tea.Cmd)`.

### Focusable Component (Before)

```go
type myField struct {
    value   string
    focused bool
}

func (f *myField) Init() tea.Cmd                           { return nil }
func (f *myField) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return f, nil }
func (f *myField) View() string                            { return f.value }
func (f *myField) Focus() tea.Cmd                          { f.focused = true; return nil }
func (f *myField) Blur()                                   { f.focused = false }

fm := nav.NewFocusManager(field1, field2, field3)
```

### Focusable Component (After)

```go
type myField struct {
    value   string
    focused bool
}

func (f *myField) Init() tea.Cmd { return nil }
func (f *myField) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg.(type) {
    case nav.FocusMsg:
        f.focused = true
    case nav.BlurMsg:
        f.focused = false
    }
    return f, nil
}
func (f *myField) View() string { return f.value }

fm, cmd := nav.NewFocusManager(field1, field2, field3)
```

Changes:
1. Remove `Focus()` and `Blur()` methods
2. Handle `FocusMsg` and `BlurMsg` in `Update`
3. `NewFocusManager` now returns `(FocusManager, tea.Cmd)`

### Screen with FocusManager (Before)

```go
func (s editScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
    var cmd tea.Cmd
    s.focus, cmd = s.focus.Update(msg)
    return s, cmd
}
```

### Screen with FocusManager (After)

```go
func (s editScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    s.focus, cmd = s.focus.Update(msg)
    return s, cmd
}
```

Only the `Update` return type changes. FocusManager delegation is identical.

## Test Scenarios

### Scenario 1: Plain tea.Model with Stack

```go
// Any tea.Model works — no Screen interface needed.
type counterModel struct{ count int }
func (m counterModel) Init() tea.Cmd                           { return nil }
func (m counterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m counterModel) View() string                            { return fmt.Sprint(m.count) }

stack := nav.NewStack(counterModel{})
// Push, pop, replace all work with standard tea.Model.
```

### Scenario 2: Plain tea.Model with FocusManager

```go
// Any tea.Model works — no Focusable interface needed.
type simpleItem struct{ name string }
func (s simpleItem) Init() tea.Cmd                           { return nil }
func (s simpleItem) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return s, nil }
func (s simpleItem) View() string                            { return s.name }

fm, _ := nav.NewFocusManager(simpleItem{name: "A"}, simpleItem{name: "B"})
// FocusMsg is delivered to A's Update. A ignores it (no handler). No panic.
```

### Scenario 3: Focus State Persists

```go
type field struct{ focused bool }
func (f field) Init() tea.Cmd { return nil }
func (f field) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg.(type) {
    case nav.FocusMsg:
        f.focused = true
    case nav.BlurMsg:
        f.focused = false
    }
    return f, nil
}
func (f field) View() string {
    if f.focused { return "[focused]" }
    return "[blurred]"
}

fm, _ := nav.NewFocusManager(field{}, field{})
// First field receives FocusMsg, sets focused=true — state persists.
fm, _ = fm.Update(tea.KeyMsg{Type: tea.KeyTab})
// First field receives BlurMsg (focused=false), second receives FocusMsg (focused=true).
```
