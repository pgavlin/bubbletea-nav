package nav

import (
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// mockScreen is a minimal tea.Model for testing.
type mockScreen struct {
	id       string
	initCmd  tea.Cmd
	lastMsg  tea.Msg
	viewText string
}

func newMockScreen(id string) mockScreen {
	return mockScreen{id: id, viewText: id}
}

func (s mockScreen) Init() tea.Cmd { return s.initCmd }
func (s mockScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	s.lastMsg = msg
	return s, nil
}
func (s mockScreen) View() string { return s.viewText }

// recordingScreen implements tea.Model, recording all messages received.
type recordingScreen struct {
	messages []tea.Msg
	viewText string
}

func (s recordingScreen) Init() tea.Cmd { return nil }
func (s recordingScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	s.messages = append(s.messages, msg)
	return s, nil
}
func (s recordingScreen) View() string { return s.viewText }

// T014: Tests for NewStack, Init, View, Depth, String.
func TestNewStack(t *testing.T) {
	t.Run("nil root panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic for nil root")
			}
		}()
		NewStack(nil)
	})

	t.Run("valid root", func(t *testing.T) {
		root := newMockScreen("root")
		s := NewStack(root)
		if s.Depth() != 1 {
			t.Fatalf("expected depth 1, got %d", s.Depth())
		}
	})
}

func TestStackInit(t *testing.T) {
	called := false
	root := mockScreen{
		id:       "root",
		viewText: "root",
		initCmd: func() tea.Msg {
			called = true
			return nil
		},
	}
	s := NewStack(root)
	cmd := s.Init()
	if cmd == nil {
		t.Fatal("expected non-nil cmd from Init")
	}
	cmd()
	if !called {
		t.Fatal("expected root Init to be called")
	}
}

func TestStackView(t *testing.T) {
	root := newMockScreen("root-view")
	s := NewStack(root)
	if got := s.View(); got != "root-view" {
		t.Fatalf("expected %q, got %q", "root-view", got)
	}
}

func TestStackDepthAndString(t *testing.T) {
	root := newMockScreen("root")
	s := NewStack(root)
	if s.Depth() != 1 {
		t.Fatalf("expected depth 1, got %d", s.Depth())
	}
	if got := s.String(); got != "Stack[1 screens, active: 0]" {
		t.Fatalf("expected %q, got %q", "Stack[1 screens, active: 0]", got)
	}
}

// T015: Tests for push, pop, replace, message forwarding.
func TestStackPush(t *testing.T) {
	tests := []struct {
		name           string
		pushCount      int
		wantDepth      int
		wantActiveView string
	}{
		{"push one", 1, 2, "screen-1"},
		{"push two", 2, 3, "screen-2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := newMockScreen("root")
			var model tea.Model = NewStack(root)

			for i := 1; i <= tt.pushCount; i++ {
				screen := newMockScreen(fmt.Sprintf("screen-%d", i))
				model, _ = model.Update(PushMsg{Screen: screen})
			}

			stack := model.(Stack)
			if stack.Depth() != tt.wantDepth {
				t.Fatalf("expected depth %d, got %d", tt.wantDepth, stack.Depth())
			}
			if got := stack.View(); got != tt.wantActiveView {
				t.Fatalf("expected view %q, got %q", tt.wantActiveView, got)
			}
		})
	}
}

func TestStackPop(t *testing.T) {
	tests := []struct {
		name           string
		pushCount      int
		popCount       int
		wantDepth      int
		wantActiveView string
	}{
		{"pop to root", 1, 1, 1, "root"},
		{"pop one of two", 2, 1, 2, "screen-1"},
		{"pop on root ignored", 0, 1, 1, "root"},
		{"pop on root ignored twice", 0, 2, 1, "root"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := newMockScreen("root")
			var model tea.Model = NewStack(root)

			for i := 1; i <= tt.pushCount; i++ {
				screen := newMockScreen(fmt.Sprintf("screen-%d", i))
				model, _ = model.Update(PushMsg{Screen: screen})
			}

			for i := 0; i < tt.popCount; i++ {
				model, _ = model.Update(PopMsg{})
			}

			stack := model.(Stack)
			if stack.Depth() != tt.wantDepth {
				t.Fatalf("expected depth %d, got %d", tt.wantDepth, stack.Depth())
			}
			if got := stack.View(); got != tt.wantActiveView {
				t.Fatalf("expected view %q, got %q", tt.wantActiveView, got)
			}
		})
	}
}

func TestStackPopPreservesState(t *testing.T) {
	root := mockScreen{id: "root", viewText: "root-with-state"}
	var model tea.Model = NewStack(root)

	model, _ = model.Update(PushMsg{Screen: newMockScreen("detail")})
	model, _ = model.Update(PopMsg{})

	stack := model.(Stack)
	if got := stack.View(); got != "root-with-state" {
		t.Fatalf("expected preserved view %q, got %q", "root-with-state", got)
	}
}

func TestStackReplace(t *testing.T) {
	root := newMockScreen("root")
	var model tea.Model = NewStack(root)

	model, _ = model.Update(PushMsg{Screen: newMockScreen("old")})
	model, _ = model.Update(ReplaceMsg{Screen: newMockScreen("new")})

	stack := model.(Stack)
	if stack.Depth() != 2 {
		t.Fatalf("expected depth 2 after replace, got %d", stack.Depth())
	}
	if got := stack.View(); got != "new" {
		t.Fatalf("expected view %q, got %q", "new", got)
	}

	// Pop to verify root is still there.
	model, _ = model.Update(PopMsg{})
	stack = model.(Stack)
	if got := stack.View(); got != "root" {
		t.Fatalf("expected root view %q, got %q", "root", got)
	}
}

func TestStackMessageForwarding(t *testing.T) {
	type customMsg struct{ data string }

	root := recordingScreen{viewText: "recorder"}
	var model tea.Model = NewStack(root)
	model, _ = model.Update(customMsg{data: "hello"})

	stack := model.(Stack)
	active := stack.screens[0].(recordingScreen)
	if len(active.messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(active.messages))
	}
	if active.messages[0].(customMsg).data != "hello" {
		t.Fatalf("expected message data %q, got %q", "hello", active.messages[0].(customMsg).data)
	}
}

func TestStackPushCallsInit(t *testing.T) {
	initCalled := false
	pushed := mockScreen{
		id:       "pushed",
		viewText: "pushed",
		initCmd: func() tea.Msg {
			initCalled = true
			return nil
		},
	}

	var model tea.Model = NewStack(newMockScreen("root"))
	_, cmd := model.Update(PushMsg{Screen: pushed})

	if cmd == nil {
		t.Fatal("expected Init cmd from pushed screen")
	}
	cmd()
	if !initCalled {
		t.Fatal("expected pushed screen Init to be called")
	}
}

func TestStackReplaceCallsInit(t *testing.T) {
	initCalled := false
	replacement := mockScreen{
		id:       "new",
		viewText: "new",
		initCmd: func() tea.Msg {
			initCalled = true
			return nil
		},
	}

	var model tea.Model = NewStack(newMockScreen("root"))
	_, cmd := model.Update(ReplaceMsg{Screen: replacement})

	if cmd == nil {
		t.Fatal("expected Init cmd from replacement screen")
	}
	cmd()
	if !initCalled {
		t.Fatal("expected replacement screen Init to be called")
	}
}

// --- Phase 6: US4 Integration Tests ---

// focusScreen is a tea.Model with a FocusManager for integration testing.
type focusScreen struct {
	id    string
	focus FocusManager
	items []*mockFocusable
}

func newFocusScreen(id string, itemCount int) *focusScreen {
	items := make([]*mockFocusable, itemCount)
	models := make([]tea.Model, itemCount)
	for i := range items {
		items[i] = &mockFocusable{}
		models[i] = items[i]
	}
	fm, _ := NewFocusManager(models...)
	return &focusScreen{
		id:    id,
		focus: fm,
		items: items,
	}
}

func (s *focusScreen) Init() tea.Cmd { return nil }
func (s *focusScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	s.focus, cmd = s.focus.Update(msg)
	return s, cmd
}
func (s *focusScreen) View() string { return s.id }

// T031: Screen A focus preserved across push/pop cycle.
func TestIntegrationFocusPreservedAcrossPushPop(t *testing.T) {
	screenA := newFocusScreen("A", 3)
	var model tea.Model = NewStack(screenA)

	// Move focus to item 1 on Screen A.
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	if screenA.focus.FocusedIndex() != 1 {
		t.Fatalf("expected A focus at 1, got %d", screenA.focus.FocusedIndex())
	}

	// Push Screen B.
	screenB := newFocusScreen("B", 2)
	model, _ = model.Update(PushMsg{Screen: screenB})

	// Pop Screen B.
	model, _ = model.Update(PopMsg{})

	// Screen A's focus should still be on item 1.
	stack := model.(Stack)
	active := stack.screens[len(stack.screens)-1].(*focusScreen)
	if active.id != "A" {
		t.Fatalf("expected active screen A, got %s", active.id)
	}
	if active.focus.FocusedIndex() != 1 {
		t.Fatalf("expected A focus at 1 after pop, got %d", active.focus.FocusedIndex())
	}
}

// T032: New Screen B instance gets default focus, not previous focus.
func TestIntegrationNewScreenDefaultFocus(t *testing.T) {
	screenA := newFocusScreen("A", 3)
	var model tea.Model = NewStack(screenA)

	// Push first Screen B and move focus.
	screenB1 := newFocusScreen("B1", 3)
	model, _ = model.Update(PushMsg{Screen: screenB1})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})
	if screenB1.focus.FocusedIndex() != 2 {
		t.Fatalf("expected B1 focus at 2, got %d", screenB1.focus.FocusedIndex())
	}

	// Pop B1.
	model, _ = model.Update(PopMsg{})

	// Push a new Screen B instance.
	screenB2 := newFocusScreen("B2", 3)
	model, _ = model.Update(PushMsg{Screen: screenB2})

	// New B should start with default focus (index 0).
	stack := model.(Stack)
	active := stack.screens[len(stack.screens)-1].(*focusScreen)
	if active.id != "B2" {
		t.Fatalf("expected active screen B2, got %s", active.id)
	}
	if active.focus.FocusedIndex() != 0 {
		t.Fatalf("expected B2 default focus at 0, got %d", active.focus.FocusedIndex())
	}
}
