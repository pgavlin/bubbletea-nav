package nav

// ScreenAppearCause identifies why a screen received ScreenAppearedMsg.
// Screens that need to distinguish a fresh push from a pop-reveal inspect
// this field; screens that don't care continue to ignore it. The zero
// value is ScreenAppearCausePushed so existing handlers that constructed
// ScreenAppearedMsg{} manually keep landing on the most-common intent.
type ScreenAppearCause int

const (
	// ScreenAppearCausePushed is the cause when a screen was pushed
	// onto the stack via PushMsg.
	ScreenAppearCausePushed ScreenAppearCause = iota

	// ScreenAppearCauseRevealed is the cause when a screen that was
	// already on the stack is now top because a screen above it was
	// popped via PopMsg.
	ScreenAppearCauseRevealed

	// ScreenAppearCauseReplaced is the cause when a screen was swapped
	// into the top slot via ReplaceMsg.
	ScreenAppearCauseReplaced
)

// ScreenAppearedMsg is sent to a screen's Update method when it
// becomes the active (topmost) screen after a push, pop-reveal,
// or replace operation. Cause names which of those happened.
type ScreenAppearedMsg struct {
	Cause ScreenAppearCause
}

// ScreenDisappearedMsg is sent to a screen's Update method when it
// loses active status due to being pushed over, popped, or replaced.
type ScreenDisappearedMsg struct{}
