package nav

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

// lifecycleScreen records lifecycle events via Update.
type lifecycleScreen struct {
	id          string
	events      []string
	viewText    string
	initCmd     tea.Cmd
	onAppear    tea.Cmd
	onDisappear tea.Cmd
}

func newLifecycleScreen(id string) *lifecycleScreen {
	return &lifecycleScreen{id: id, viewText: id}
}

func (s *lifecycleScreen) Init() tea.Cmd { return s.initCmd }
func (s *lifecycleScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case ScreenAppearedMsg:
		s.events = append(s.events, "appeared")
		return s, s.onAppear
	case ScreenDisappearedMsg:
		s.events = append(s.events, "disappeared")
		return s, s.onDisappear
	}
	return s, nil
}
func (s *lifecycleScreen) View() tea.View { return tea.NewView(s.viewText) }

// T029: Lifecycle events for push, pop, replace, and non-lifecycle screens.
func TestLifecyclePush(t *testing.T) {
	root := newLifecycleScreen("root")
	detail := newLifecycleScreen("detail")

	var model tea.Model = NewStack(root)
	model, _ = model.Update(PushMsg{Screen: detail})

	if len(root.events) != 1 || root.events[0] != "disappeared" {
		t.Fatalf("expected root [disappeared], got %v", root.events)
	}
	if len(detail.events) != 1 || detail.events[0] != "appeared" {
		t.Fatalf("expected detail [appeared], got %v", detail.events)
	}
}

func TestLifecyclePop(t *testing.T) {
	root := newLifecycleScreen("root")
	detail := newLifecycleScreen("detail")

	var model tea.Model = NewStack(root)
	model, _ = model.Update(PushMsg{Screen: detail})

	// Reset events.
	root.events = nil
	detail.events = nil

	model, _ = model.Update(PopMsg{})

	if len(detail.events) != 1 || detail.events[0] != "disappeared" {
		t.Fatalf("expected detail [disappeared], got %v", detail.events)
	}
	if len(root.events) != 1 || root.events[0] != "appeared" {
		t.Fatalf("expected root [appeared], got %v", root.events)
	}
}

func TestLifecycleReplace(t *testing.T) {
	root := newLifecycleScreen("root")
	old := newLifecycleScreen("old")
	replacement := newLifecycleScreen("new")

	var model tea.Model = NewStack(root)
	model, _ = model.Update(PushMsg{Screen: old})

	// Reset events.
	root.events = nil
	old.events = nil

	model, _ = model.Update(ReplaceMsg{Screen: replacement})

	if len(old.events) != 1 || old.events[0] != "disappeared" {
		t.Fatalf("expected old [disappeared], got %v", old.events)
	}
	if len(replacement.events) != 1 || replacement.events[0] != "appeared" {
		t.Fatalf("expected replacement [appeared], got %v", replacement.events)
	}
	// Root should not have received any events.
	if len(root.events) != 0 {
		t.Fatalf("expected root [], got %v", root.events)
	}
}

func TestLifecycleScreenIgnoresUnhandledMessages(t *testing.T) {
	// mockScreen does not handle ScreenAppearedMsg/ScreenDisappearedMsg.
	root := newMockScreen("root")
	detail := newMockScreen("detail")

	var model tea.Model = NewStack(root)
	// This should not panic even though screens don't handle lifecycle messages.
	model, _ = model.Update(PushMsg{Screen: detail})
	model, _ = model.Update(PopMsg{})

	// If we get here without panic, the test passes.
	stack := model.(Stack)
	if stack.Depth() != 1 {
		t.Fatalf("expected depth 1, got %d", stack.Depth())
	}
}

func TestLifecycleAppearedReturnsCmd(t *testing.T) {
	type loadMsg struct{}
	root := newLifecycleScreen("root")
	detail := newLifecycleScreen("detail")
	detail.onAppear = func() tea.Msg { return loadMsg{} }

	var model tea.Model = NewStack(root)
	_, cmd := model.Update(PushMsg{Screen: detail})

	if cmd == nil {
		t.Fatal("expected command from Appeared")
	}
}

func TestLifecycleFullSequence(t *testing.T) {
	root := newLifecycleScreen("root")
	a := newLifecycleScreen("A")
	b := newLifecycleScreen("B")

	var model tea.Model = NewStack(root)

	// Push A.
	model, _ = model.Update(PushMsg{Screen: a})
	// Push B.
	model, _ = model.Update(PushMsg{Screen: b})
	// Pop B.
	model, _ = model.Update(PopMsg{})
	// Pop A.
	model, _ = model.Update(PopMsg{})

	expectedRoot := []string{"disappeared", "appeared"}
	expectedA := []string{"appeared", "disappeared", "appeared", "disappeared"}
	expectedB := []string{"appeared", "disappeared"}

	assertEvents(t, "root", root.events, expectedRoot)
	assertEvents(t, "A", a.events, expectedA)
	assertEvents(t, "B", b.events, expectedB)
}

func assertEvents(t *testing.T, name string, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s: expected %d events %v, got %d events %v", name, len(want), want, len(got), got)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("%s: event %d expected %q, got %q", name, i, want[i], got[i])
		}
	}
}

func TestLifecycleAppearedStatePersists(t *testing.T) {
	// Screen updates its viewText when it receives ScreenAppearedMsg.
	root := newLifecycleScreen("root")
	detail := newLifecycleScreen("detail")
	detail.onAppear = nil // no command, but Update still mutates state
	// Override detail's Update to change viewText on appear.
	var model tea.Model = NewStack(root)

	// Push detail — its Update receives ScreenAppearedMsg and records "appeared".
	model, _ = model.Update(PushMsg{Screen: detail})

	// The lifecycleScreen mock appends "appeared" to events.
	// Verify the state change persists via the stack's View.
	stack := model.(Stack)
	if got := stack.View().Content; got != "detail" {
		t.Fatalf("expected view %q, got %q", "detail", got)
	}
	if len(detail.events) != 1 || detail.events[0] != "appeared" {
		t.Fatalf("expected detail events [appeared], got %v", detail.events)
	}

	// Use a custom screen that changes View on appear.
	adapter := &appearStateScreen{viewText: "not-visited"}
	model = NewStack(newLifecycleScreen("root2"))
	model, _ = model.Update(PushMsg{Screen: adapter})

	stack = model.(Stack)
	if got := stack.View().Content; got != "visited" {
		t.Fatalf("expected view %q after appear, got %q", "visited", got)
	}
}

// appearStateScreen changes its view on ScreenAppearedMsg.
type appearStateScreen struct {
	viewText string
}

func (s *appearStateScreen) Init() tea.Cmd { return nil }
func (s *appearStateScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(ScreenAppearedMsg); ok {
		s.viewText = "visited"
	}
	return s, nil
}
func (s *appearStateScreen) View() tea.View { return tea.NewView(s.viewText) }

func TestLifecycleDisappearedReturnsCmd(t *testing.T) {
	type cleanupMsg struct{}
	root := newLifecycleScreen("root")
	root.onDisappear = func() tea.Msg { return cleanupMsg{} }
	detail := newLifecycleScreen("detail")

	var model tea.Model = NewStack(root)
	_, cmd := model.Update(PushMsg{Screen: detail})

	// The batched command should include root's onDisappear command.
	if cmd == nil {
		t.Fatal("expected command from disappeared handler")
	}
}

func TestLifecycleDisappearedStatePersists(t *testing.T) {
	// Screen that changes viewText when it receives ScreenDisappearedMsg.
	root := &disappearStateScreen{viewText: "active"}
	var model tea.Model = NewStack(root)

	// Push detail — root receives ScreenDisappearedMsg, sets viewText to "hidden".
	model, _ = model.Update(PushMsg{Screen: newLifecycleScreen("detail")})

	// Pop detail — root is revealed.
	model, _ = model.Update(PopMsg{})

	// Root's viewText should reflect the state change from disappeared.
	stack := model.(Stack)
	if got := stack.View().Content; got != "hidden" {
		t.Fatalf("expected view %q after disappear+reappear, got %q", "hidden", got)
	}
}

// disappearStateScreen changes its view on ScreenDisappearedMsg.
type disappearStateScreen struct {
	viewText string
}

func (s *disappearStateScreen) Init() tea.Cmd { return nil }
func (s *disappearStateScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(ScreenDisappearedMsg); ok {
		s.viewText = "hidden"
	}
	return s, nil
}
func (s *disappearStateScreen) View() tea.View { return tea.NewView(s.viewText) }

// T030: Re-entrant stack modification during lifecycle.
func TestLifecycleReentrantPush(t *testing.T) {
	// When screen A appears, it tries to push screen B.
	root := newLifecycleScreen("root")
	b := newLifecycleScreen("B")
	a := newLifecycleScreen("A")
	a.onAppear = func() tea.Msg {
		return PushMsg{Screen: b}
	}

	var model tea.Model = NewStack(root)
	model, cmd := model.Update(PushMsg{Screen: a})

	// The push of A should complete. B should be queued.
	// Execute the batched command to process the queued PushMsg from Appeared.
	if cmd != nil {
		msg := cmd()
		if msg != nil {
			model, _ = model.Update(msg)
		}
	}

	stack := model.(Stack)
	// After processing: root -> A -> B
	if stack.Depth() != 3 {
		t.Fatalf("expected depth 3 (root->A->B), got %d", stack.Depth())
	}
	if got := stack.View().Content; got != "B" {
		t.Fatalf("expected active view %q, got %q", "B", got)
	}
}

// causeRecordingScreen records the Cause of every ScreenAppearedMsg.
type causeRecordingScreen struct {
	id     string
	causes []ScreenAppearCause
}

func newCauseRecordingScreen(id string) *causeRecordingScreen {
	return &causeRecordingScreen{id: id}
}

func (s *causeRecordingScreen) Init() tea.Cmd { return nil }
func (s *causeRecordingScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m, ok := msg.(ScreenAppearedMsg); ok {
		s.causes = append(s.causes, m.Cause)
	}
	return s, nil
}
func (s *causeRecordingScreen) View() tea.View { return tea.NewView(s.id) }

func TestScreenAppearedCausePushed(t *testing.T) {
	root := newCauseRecordingScreen("root")
	detail := newCauseRecordingScreen("detail")

	var model tea.Model = NewStack(root)
	model, _ = model.Update(PushMsg{Screen: detail})

	if len(detail.causes) != 1 || detail.causes[0] != ScreenAppearCausePushed {
		t.Fatalf("expected detail causes [Pushed], got %v", detail.causes)
	}
}

func TestScreenAppearedCauseRevealed(t *testing.T) {
	root := newCauseRecordingScreen("root")
	detail := newCauseRecordingScreen("detail")

	var model tea.Model = NewStack(root)
	model, _ = model.Update(PushMsg{Screen: detail})
	root.causes = nil // reset to isolate the pop-back

	model, _ = model.Update(PopMsg{})

	if len(root.causes) != 1 || root.causes[0] != ScreenAppearCauseRevealed {
		t.Fatalf("expected root causes [Revealed] after pop, got %v", root.causes)
	}
}

func TestScreenAppearedCauseReplaced(t *testing.T) {
	root := newCauseRecordingScreen("root")
	old := newCauseRecordingScreen("old")
	replacement := newCauseRecordingScreen("new")

	var model tea.Model = NewStack(root)
	model, _ = model.Update(PushMsg{Screen: old})
	replacement.causes = nil // sanity: replacement hasn't been on stack yet

	model, _ = model.Update(ReplaceMsg{Screen: replacement})

	if len(replacement.causes) != 1 || replacement.causes[0] != ScreenAppearCauseReplaced {
		t.Fatalf("expected replacement causes [Replaced], got %v", replacement.causes)
	}
}
