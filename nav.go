package nav

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// PushMsg requests pushing a screen onto the navigation stack.
type PushMsg struct{ Screen tea.Model }

// PopMsg requests popping the top screen from the navigation stack.
type PopMsg struct{}

// ReplaceMsg requests replacing the top screen on the stack.
type ReplaceMsg struct{ Screen tea.Model }

// Push returns a command that pushes a screen onto the stack.
func Push(screen tea.Model) tea.Cmd {
	return func() tea.Msg {
		return PushMsg{Screen: screen}
	}
}

// Pop returns a command that pops the top screen from the stack.
func Pop() tea.Cmd {
	return func() tea.Msg {
		return PopMsg{}
	}
}

// Replace returns a command that replaces the top screen.
func Replace(screen tea.Model) tea.Cmd {
	return func() tea.Msg {
		return ReplaceMsg{Screen: screen}
	}
}

// Stack manages an ordered stack of screens. It implements tea.Model.
// The topmost screen is the active screen and receives all messages
// except navigation messages, which the Stack intercepts.
type Stack struct {
	screens     []tea.Model
	pendingOps  []tea.Msg
	inLifecycle bool
}

// NewStack creates a navigation stack with the given root screen.
// The root screen cannot be popped. Panics if root is nil.
func NewStack(root tea.Model) Stack {
	if root == nil {
		panic("nav: root screen must not be nil")
	}
	return Stack{screens: []tea.Model{root}}
}

// Init initializes the stack by calling Init on the root screen.
func (s Stack) Init() tea.Cmd {
	return s.screens[len(s.screens)-1].Init()
}

// Update processes a message. Navigation messages (PushMsg, PopMsg,
// ReplaceMsg) modify the stack. All other messages are forwarded to
// the active screen.
func (s Stack) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case PushMsg, PopMsg, ReplaceMsg:
		if s.inLifecycle {
			s.pendingOps = append(s.pendingOps, msg)
			return s, nil
		}
		return s.handleNav(msg)
	}

	// Forward to active screen.
	active := s.screens[len(s.screens)-1]
	updated, cmd := active.Update(msg)
	s.screens[len(s.screens)-1] = updated
	return s, cmd
}

// handleNav processes a navigation message and dispatches lifecycle events.
func (s Stack) handleNav(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case PushMsg:
		s.screens = append(s.screens, msg.Screen)
		newIdx := len(s.screens) - 1

		// Lifecycle: old top disappears, new top appears.
		s.inLifecycle = true
		updated, cmd := s.screens[newIdx-1].Update(ScreenDisappearedMsg{})
		s.screens[newIdx-1] = updated
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		updated, cmd = s.screens[newIdx].Update(ScreenAppearedMsg{})
		s.screens[newIdx] = updated
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		s.inLifecycle = false

		if cmd := s.screens[newIdx].Init(); cmd != nil {
			cmds = append(cmds, cmd)
		}

	case PopMsg:
		if len(s.screens) <= 1 {
			return s, nil
		}
		popped := s.screens[len(s.screens)-1]
		s.screens = s.screens[:len(s.screens)-1]
		topIdx := len(s.screens) - 1

		// Lifecycle: popped screen disappears, revealed screen appears.
		s.inLifecycle = true
		_, cmd := popped.Update(ScreenDisappearedMsg{})
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		updated, cmd := s.screens[topIdx].Update(ScreenAppearedMsg{})
		s.screens[topIdx] = updated
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		s.inLifecycle = false

	case ReplaceMsg:
		oldTop := s.screens[len(s.screens)-1]
		s.screens[len(s.screens)-1] = msg.Screen
		topIdx := len(s.screens) - 1

		// Lifecycle: old top disappears, new top appears.
		s.inLifecycle = true
		_, cmd := oldTop.Update(ScreenDisappearedMsg{})
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		updated, cmd := s.screens[topIdx].Update(ScreenAppearedMsg{})
		s.screens[topIdx] = updated
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		s.inLifecycle = false

		if cmd := s.screens[topIdx].Init(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	// Process any navigation messages queued during lifecycle.
	for len(s.pendingOps) > 0 {
		pending := s.pendingOps
		s.pendingOps = nil
		for _, op := range pending {
			var result tea.Model
			var cmd tea.Cmd
			result, cmd = s.handleNav(op)
			s = result.(Stack)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}

	return s, tea.Batch(cmds...)
}

// View renders the active screen.
func (s Stack) View() string {
	return s.screens[len(s.screens)-1].View()
}

// Depth returns the number of screens on the stack.
func (s Stack) Depth() int {
	return len(s.screens)
}

// String returns a human-readable representation of the stack.
func (s Stack) String() string {
	return fmt.Sprintf("Stack[%d screens, active: %d]", len(s.screens), len(s.screens)-1)
}
