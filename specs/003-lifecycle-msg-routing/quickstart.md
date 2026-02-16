# Quickstart: Message-Based Lifecycle Notifications

**Feature**: 003-lifecycle-msg-routing

## Before (LifecycleScreen interface)

```go
type myScreen struct {
    data    string
    stale   bool
}

func (s *myScreen) Init() tea.Cmd { return nil }

func (s *myScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
    // Only handles user messages — lifecycle is separate
    return s, nil
}

func (s *myScreen) View() string { return s.data }

// Extra interface methods required:
func (s *myScreen) Appeared() tea.Cmd {
    // Problem: can't update s.stale here and have it persist,
    // because Appeared() doesn't return a Screen value.
    return loadDataCmd
}

func (s *myScreen) Disappeared() {
    // Problem: can't update state at all — returns nothing.
}
```

## After (message-based)

```go
type myScreen struct {
    data    string
    stale   bool
}

func (s myScreen) Init() tea.Cmd { return nil }

func (s myScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
    switch msg.(type) {
    case nav.ScreenAppearedMsg:
        // State update works naturally — returned via (Screen, Cmd)
        s.stale = false
        return s, loadDataCmd
    case nav.ScreenDisappearedMsg:
        // Can update state and return commands
        s.stale = true
        return s, nil
    }
    // Handle other messages...
    return s, nil
}

func (s myScreen) View() string { return s.data }

// No extra interface methods needed!
```

## Key differences

1. **No extra interface** — just handle two message types in your existing `Update`
2. **State updates work** — `Update` returns `(Screen, Cmd)`, so state changes persist
3. **Commands from disappear** — `ScreenDisappearedMsg` handling can now return commands
4. **Uniform pattern** — lifecycle events use the same message-passing pattern as everything else

## Migration checklist

- [ ] Remove `Appeared() tea.Cmd` method
- [ ] Remove `Disappeared()` method
- [ ] Add `case nav.ScreenAppearedMsg:` to `Update`
- [ ] Add `case nav.ScreenDisappearedMsg:` to `Update`
- [ ] Move logic from old methods into the new cases
- [ ] Verify: `go test ./...` passes
