package nav

// ScreenAppearedMsg is sent to a screen's Update method when it
// becomes the active (topmost) screen after a push, pop-reveal,
// or replace operation.
type ScreenAppearedMsg struct{}

// ScreenDisappearedMsg is sent to a screen's Update method when it
// loses active status due to being pushed over, popped, or replaced.
type ScreenDisappearedMsg struct{}
