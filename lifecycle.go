package nav

import tea "github.com/charmbracelet/bubbletea"

// LifecycleScreen is an optional interface for screens that need
// notifications when they become visible or hidden due to stack
// changes. Implement this in addition to Screen.
type LifecycleScreen interface {
	Screen

	// Appeared is called when this screen becomes the active screen
	// (after push, pop-reveal, or replace). Returns a command for
	// async initialization.
	Appeared() tea.Cmd

	// Disappeared is called when this screen loses active status
	// (pushed over, popped, or replaced). Synchronous cleanup.
	Disappeared()
}

// ScreenAppearedMsg is sent to a screen when it becomes active.
type ScreenAppearedMsg struct{}

// ScreenDisappearedMsg is sent to a screen when it loses active status.
type ScreenDisappearedMsg struct{}
