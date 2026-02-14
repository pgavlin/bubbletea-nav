package nav

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// lifecycleScreen records Appeared/Disappeared events.
type lifecycleScreen struct {
	id       string
	events   []string
	viewText string
	initCmd  tea.Cmd
	onAppear tea.Cmd
}

func newLifecycleScreen(id string) *lifecycleScreen {
	return &lifecycleScreen{id: id, viewText: id}
}

func (s *lifecycleScreen) Init() tea.Cmd { return s.initCmd }
func (s *lifecycleScreen) Update(msg tea.Msg) (Screen, tea.Cmd) {
	return s, nil
}
func (s *lifecycleScreen) View() string { return s.viewText }
func (s *lifecycleScreen) Appeared() tea.Cmd {
	s.events = append(s.events, "appeared")
	return s.onAppear
}
func (s *lifecycleScreen) Disappeared() {
	s.events = append(s.events, "disappeared")
}

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

func TestLifecycleNonLifecycleScreenSkipped(t *testing.T) {
	// mockScreen does not implement LifecycleScreen.
	root := newMockScreen("root")
	detail := newMockScreen("detail")

	var model tea.Model = NewStack(root)
	// This should not panic even though screens don't implement lifecycle.
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
	if got := stack.View(); got != "B" {
		t.Fatalf("expected active view %q, got %q", "B", got)
	}
}
