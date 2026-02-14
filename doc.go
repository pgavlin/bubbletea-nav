// Package nav provides screen navigation and focus management for
// Bubble Tea terminal UI applications.
//
// The two core primitives are Stack and FocusManager:
//
//   - Stack manages an ordered stack of screens with push, pop, and
//     replace operations. It implements tea.Model and can be passed
//     directly to tea.NewProgram.
//
//   - FocusManager manages sequential focus order within a screen,
//     handling Tab/Shift+Tab cycling and mouse-click targeting. It is
//     a value type that screens hold as a field.
//
// # Basic Usage
//
//	// Create a navigation stack with a root screen.
//	stack := nav.NewStack(myRootScreen{})
//	p := tea.NewProgram(stack)
//	p.Run()
//
//	// In a screen's Update method, push a new screen:
//	func (s myScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
//	    return s, nav.Push(detailScreen{})
//	}
//
//	// Or pop back to the previous screen:
//	func (s myScreen) Update(msg tea.Msg) (nav.Screen, tea.Cmd) {
//	    return s, nav.Pop()
//	}
package nav
