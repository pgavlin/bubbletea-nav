package nav

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// mockFocusable is a minimal tea.Model for focus testing.
// It handles FocusMsg and BlurMsg in Update to track focus state.
type mockFocusable struct {
	focused  bool
	messages []tea.Msg
}

func (f *mockFocusable) Init() tea.Cmd { return nil }
func (f *mockFocusable) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case FocusMsg:
		f.focused = true
	case BlurMsg:
		f.focused = false
	}
	f.messages = append(f.messages, msg)
	return f, nil
}
func (f *mockFocusable) View() string { return "" }

// mockBoundedFocusable is a tea.Model with Bounded for mouse testing.
type mockBoundedFocusable struct {
	mockFocusable
	x, y, w, h int
}

func (f *mockBoundedFocusable) Bounds() (x, y, width, height int) {
	return f.x, f.y, f.w, f.h
}

// tabMsg creates a KeyMsg for Tab.
func tabMsg() tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyTab}
}

// shiftTabMsg creates a KeyMsg for Shift+Tab.
func shiftTabMsg() tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyShiftTab}
}

// mouseClick creates a MouseMsg press at (x, y).
func mouseClick(x, y int) tea.MouseMsg {
	return tea.MouseMsg{
		X:      x,
		Y:      y,
		Action: tea.MouseActionPress,
		Button: tea.MouseButtonLeft,
	}
}

func TestNewFocusManager(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		fm, cmd := NewFocusManager()
		if fm.FocusedIndex() != -1 {
			t.Fatalf("expected -1, got %d", fm.FocusedIndex())
		}
		if cmd != nil {
			t.Fatal("expected nil cmd for empty focus manager")
		}
	})

	t.Run("one item focused", func(t *testing.T) {
		item := &mockFocusable{}
		fm, _ := NewFocusManager(item)
		if fm.FocusedIndex() != 0 {
			t.Fatalf("expected 0, got %d", fm.FocusedIndex())
		}
		if !item.focused {
			t.Fatal("expected first item to be focused")
		}
	})

	t.Run("three items first focused", func(t *testing.T) {
		items := []*mockFocusable{{}, {}, {}}
		fm, _ := NewFocusManager(items[0], items[1], items[2])
		if fm.FocusedIndex() != 0 {
			t.Fatalf("expected 0, got %d", fm.FocusedIndex())
		}
		if !items[0].focused {
			t.Fatal("expected first item focused")
		}
		if items[1].focused || items[2].focused {
			t.Fatal("expected only first item focused")
		}
	})
}

// T022: Tab forward cycling, Shift+Tab backward cycling, wrap, empty.
func TestFocusManagerTabCycling(t *testing.T) {
	tests := []struct {
		name      string
		msgs      []tea.Msg
		wantIndex int
	}{
		{"tab once from 0", []tea.Msg{tabMsg()}, 1},
		{"tab twice from 0", []tea.Msg{tabMsg(), tabMsg()}, 2},
		{"tab wraps to 0", []tea.Msg{tabMsg(), tabMsg(), tabMsg()}, 0},
		{"shift+tab from 0 wraps to last", []tea.Msg{shiftTabMsg()}, 2},
		{"shift+tab twice wraps", []tea.Msg{shiftTabMsg(), shiftTabMsg()}, 1},
		{"tab then shift+tab returns", []tea.Msg{tabMsg(), shiftTabMsg()}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items := []*mockFocusable{{}, {}, {}}
			fm, _ := NewFocusManager(items[0], items[1], items[2])

			for _, msg := range tt.msgs {
				fm, _ = fm.Update(msg)
			}

			if fm.FocusedIndex() != tt.wantIndex {
				t.Fatalf("expected index %d, got %d", tt.wantIndex, fm.FocusedIndex())
			}
			// Verify the correct item is actually focused.
			for i, item := range items {
				if i == tt.wantIndex && !item.focused {
					t.Fatalf("expected item %d to be focused", i)
				}
				if i != tt.wantIndex && item.focused {
					t.Fatalf("expected item %d to not be focused", i)
				}
			}
		})
	}
}

func TestFocusManagerEmptyNoOp(t *testing.T) {
	fm, _ := NewFocusManager()
	fm, cmd := fm.Update(tabMsg())
	if cmd != nil {
		t.Fatal("expected nil cmd for empty focus manager")
	}
	if fm.FocusedIndex() != -1 {
		t.Fatalf("expected -1, got %d", fm.FocusedIndex())
	}

	fm, cmd = fm.Update(shiftTabMsg())
	if cmd != nil {
		t.Fatal("expected nil cmd for empty focus manager on shift+tab")
	}
}

// T023: Mouse click tests.
func TestFocusManagerMouseClick(t *testing.T) {
	t.Run("click on bounded item moves focus", func(t *testing.T) {
		items := []*mockBoundedFocusable{
			{mockFocusable: mockFocusable{}, x: 0, y: 0, w: 10, h: 1},
			{mockFocusable: mockFocusable{}, x: 0, y: 1, w: 10, h: 1},
			{mockFocusable: mockFocusable{}, x: 0, y: 2, w: 10, h: 1},
		}
		fm, _ := NewFocusManager(items[0], items[1], items[2])

		// Click on item 2 (y=2).
		fm, cmd := fm.Update(mouseClick(5, 2))
		if fm.FocusedIndex() != 2 {
			t.Fatalf("expected index 2, got %d", fm.FocusedIndex())
		}
		if cmd == nil {
			t.Fatal("expected FocusChangedMsg command")
		}
		if !items[2].focused {
			t.Fatal("expected item 2 to be focused")
		}
		if items[0].focused {
			t.Fatal("expected item 0 to be blurred")
		}
	})

	t.Run("click outside all bounds no change", func(t *testing.T) {
		items := []*mockBoundedFocusable{
			{mockFocusable: mockFocusable{}, x: 0, y: 0, w: 10, h: 1},
		}
		fm, _ := NewFocusManager(items[0])

		fm, _ = fm.Update(mouseClick(20, 20))
		if fm.FocusedIndex() != 0 {
			t.Fatalf("expected index 0, got %d", fm.FocusedIndex())
		}
	})

	t.Run("click on non-bounded item no change", func(t *testing.T) {
		plain := &mockFocusable{}
		bounded := &mockBoundedFocusable{
			mockFocusable: mockFocusable{}, x: 0, y: 1, w: 10, h: 1,
		}
		fm, _ := NewFocusManager(plain, bounded)

		// Click in the area of the bounded item.
		fm, _ = fm.Update(mouseClick(5, 1))
		if fm.FocusedIndex() != 1 {
			t.Fatalf("expected index 1, got %d", fm.FocusedIndex())
		}
	})

	t.Run("click on already focused item no-op", func(t *testing.T) {
		items := []*mockBoundedFocusable{
			{mockFocusable: mockFocusable{}, x: 0, y: 0, w: 10, h: 1},
		}
		fm, _ := NewFocusManager(items[0])

		fm, cmd := fm.Update(mouseClick(5, 0))
		// Click is routed to the focused item, so cmd comes from routeMessage.
		_ = cmd
		if fm.FocusedIndex() != 0 {
			t.Fatalf("expected index 0, got %d", fm.FocusedIndex())
		}
	})

	t.Run("mouse release ignored", func(t *testing.T) {
		items := []*mockBoundedFocusable{
			{mockFocusable: mockFocusable{}, x: 0, y: 0, w: 10, h: 1},
			{mockFocusable: mockFocusable{}, x: 0, y: 1, w: 10, h: 1},
		}
		fm, _ := NewFocusManager(items[0], items[1])

		releaseMsg := tea.MouseMsg{
			X: 5, Y: 1,
			Action: tea.MouseActionRelease,
			Button: tea.MouseButtonLeft,
		}
		fm, _ = fm.Update(releaseMsg)
		if fm.FocusedIndex() != 0 {
			t.Fatalf("expected index 0, got %d", fm.FocusedIndex())
		}
	})
}

func TestFocusManagerFocusChangedMsg(t *testing.T) {
	items := []*mockFocusable{{}, {}, {}}
	fm, _ := NewFocusManager(items[0], items[1], items[2])

	fm, cmd := fm.Update(tabMsg())
	if cmd == nil {
		t.Fatal("expected FocusChangedMsg command")
	}
}

// T024: SetItems, FocusedIndex, FocusIndex, dynamic list changes.
func TestFocusManagerSetItems(t *testing.T) {
	t.Run("set items resets focus", func(t *testing.T) {
		old := []*mockFocusable{{}, {}}
		fm, _ := NewFocusManager(old[0], old[1])
		fm, _ = fm.Update(tabMsg()) // focus on item 1

		newItems := []*mockFocusable{{}, {}, {}}
		fm, _ = fm.SetItems(newItems[0], newItems[1], newItems[2])

		if fm.FocusedIndex() != 0 {
			t.Fatalf("expected index 0 after SetItems, got %d", fm.FocusedIndex())
		}
		if !newItems[0].focused {
			t.Fatal("expected first new item to be focused")
		}
		if old[1].focused {
			t.Fatal("expected old focused item to be blurred")
		}
	})

	t.Run("set empty items", func(t *testing.T) {
		items := []*mockFocusable{{}}
		fm, _ := NewFocusManager(items[0])
		fm, _ = fm.SetItems()
		if fm.FocusedIndex() != -1 {
			t.Fatalf("expected -1 for empty, got %d", fm.FocusedIndex())
		}
	})
}

func TestFocusManagerFocusIndex(t *testing.T) {
	tests := []struct {
		name      string
		index     int
		wantIndex int
	}{
		{"set to 1", 1, 1},
		{"set to 2", 2, 2},
		{"set to -1 blurs all", -1, -1},
		{"out of range high clamps", 10, 2},
		{"out of range low clamps to -1", -5, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items := []*mockFocusable{{}, {}, {}}
			fm, _ := NewFocusManager(items[0], items[1], items[2])

			fm, _ = fm.FocusIndex(tt.index)
			if fm.FocusedIndex() != tt.wantIndex {
				t.Fatalf("expected index %d, got %d", tt.wantIndex, fm.FocusedIndex())
			}

			// Verify focus state of items.
			for i, item := range items {
				if i == tt.wantIndex && !item.focused {
					t.Fatalf("expected item %d to be focused", i)
				}
				if i != tt.wantIndex && item.focused {
					t.Fatalf("expected item %d to not be focused", i)
				}
			}
		})
	}
}

func TestFocusManagerFocusIndexSameNoOp(t *testing.T) {
	items := []*mockFocusable{{}, {}}
	fm, _ := NewFocusManager(items[0], items[1])

	fm, cmd := fm.FocusIndex(0) // already focused on 0
	if cmd != nil {
		t.Fatal("expected nil cmd when setting same index")
	}
}

func TestFocusManagerString(t *testing.T) {
	items := []*mockFocusable{{}, {}, {}}
	fm, _ := NewFocusManager(items[0], items[1], items[2])
	expected := "FocusManager[3 items, focused: 0]"
	if got := fm.String(); got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

// --- T004: Message routing tests ---

// mockFocusableWithCmd returns a command from Update.
type mockFocusableWithCmd struct {
	mockFocusable
	cmd tea.Cmd
}

func (f *mockFocusableWithCmd) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case FocusMsg:
		f.focused = true
	case BlurMsg:
		f.focused = false
	}
	f.messages = append(f.messages, msg)
	return f, f.cmd
}

// T004: Table-driven tests for automatic message routing.
func TestFocusManagerMessageRouting(t *testing.T) {
	type customMsg struct{ data string }

	t.Run("non-focus key routed only to focused item", func(t *testing.T) {
		items := []*mockFocusable{{}, {}, {}}
		fm, _ := NewFocusManager(items[0], items[1], items[2])

		// Move focus to item 1.
		fm, _ = fm.Update(tabMsg())

		// Clear any messages from focus transitions.
		for _, item := range items {
			item.messages = nil
		}

		// Send a non-focus message.
		msg := customMsg{data: "hello"}
		fm, _ = fm.Update(msg)

		// Only item 1 should have received the message.
		if len(items[1].messages) != 1 {
			t.Fatalf("expected focused item to receive 1 message, got %d", len(items[1].messages))
		}
		if items[1].messages[0].(customMsg).data != "hello" {
			t.Fatal("focused item received wrong message")
		}
		if len(items[0].messages) != 0 {
			t.Fatalf("item 0 should not receive messages, got %d", len(items[0].messages))
		}
		if len(items[2].messages) != 0 {
			t.Fatalf("item 2 should not receive messages, got %d", len(items[2].messages))
		}
	})

	t.Run("command from focused item returned", func(t *testing.T) {
		type resultMsg struct{}
		item := &mockFocusableWithCmd{
			cmd: func() tea.Msg { return resultMsg{} },
		}
		fm, _ := NewFocusManager(item)

		_, cmd := fm.Update(customMsg{data: "trigger"})
		if cmd == nil {
			t.Fatal("expected command from focused item")
		}
		msg := cmd()
		if _, ok := msg.(resultMsg); !ok {
			t.Fatalf("expected resultMsg, got %T", msg)
		}
	})

	t.Run("updated state retained after routing", func(t *testing.T) {
		items := []*mockFocusable{{}, {}}
		fm, _ := NewFocusManager(items[0], items[1])

		// Clear messages from initial focus.
		items[0].messages = nil

		// Send messages to item 0 (focused).
		fm, _ = fm.Update(customMsg{data: "first"})
		fm, _ = fm.Update(customMsg{data: "second"})

		// The focused item should have recorded both messages.
		if len(items[0].messages) != 2 {
			t.Fatalf("expected 2 messages on focused item, got %d", len(items[0].messages))
		}
	})

	t.Run("unfocused items do not receive messages", func(t *testing.T) {
		items := []*mockFocusable{{}, {}, {}}
		fm, _ := NewFocusManager(items[0], items[1], items[2])

		// Clear messages from initial focus.
		for _, item := range items {
			item.messages = nil
		}

		fm, _ = fm.Update(customMsg{data: "test"})

		// Only item 0 (focused) receives the message.
		if len(items[0].messages) != 1 {
			t.Fatalf("expected 1 message on item 0, got %d", len(items[0].messages))
		}
		for i := 1; i < 3; i++ {
			if len(items[i].messages) != 0 {
				t.Fatalf("item %d should not receive messages, got %d", i, len(items[i].messages))
			}
		}
	})
}

// T005: Edge case tests for message routing.
func TestFocusManagerRoutingEdgeCases(t *testing.T) {
	type customMsg struct{ data string }

	t.Run("no routing when focusIndex is -1", func(t *testing.T) {
		items := []*mockFocusable{{}, {}}
		fm, _ := NewFocusManager(items[0], items[1])
		fm, _ = fm.FocusIndex(-1) // blur all

		// Clear messages from focus transitions.
		for _, item := range items {
			item.messages = nil
		}

		fm, cmd := fm.Update(customMsg{data: "hello"})
		if cmd != nil {
			t.Fatal("expected nil cmd when no item focused")
		}
		for i, item := range items {
			if len(item.messages) != 0 {
				t.Fatalf("item %d should not receive messages when unfocused, got %d", i, len(item.messages))
			}
		}
	})

	t.Run("no routing when items list is empty", func(t *testing.T) {
		fm, _ := NewFocusManager()
		fm, cmd := fm.Update(customMsg{data: "hello"})
		if cmd != nil {
			t.Fatal("expected nil cmd for empty focus manager")
		}
		if fm.FocusedIndex() != -1 {
			t.Fatalf("expected -1, got %d", fm.FocusedIndex())
		}
	})

	t.Run("navigation command returned unmodified", func(t *testing.T) {
		item := &mockFocusableWithCmd{
			cmd: func() tea.Msg { return PushMsg{} },
		}
		fm, _ := NewFocusManager(item)

		// Clear messages from initial focus.
		item.messages = nil

		_, cmd := fm.Update(customMsg{data: "trigger"})
		if cmd == nil {
			t.Fatal("expected command from focused item")
		}
		msg := cmd()
		if _, ok := msg.(PushMsg); !ok {
			t.Fatalf("expected PushMsg, got %T", msg)
		}
	})
}

// T007: Tab/Shift+Tab are consumed — not forwarded to any component.
func TestFocusManagerTabConsumed(t *testing.T) {
	t.Run("tab not forwarded to any item", func(t *testing.T) {
		items := []*mockFocusable{{}, {}, {}}
		fm, _ := NewFocusManager(items[0], items[1], items[2])

		// Clear messages from initial focus.
		for _, item := range items {
			item.messages = nil
		}

		fm, _ = fm.Update(tabMsg())

		// Focus should have moved to item 1.
		if fm.FocusedIndex() != 1 {
			t.Fatalf("expected index 1, got %d", fm.FocusedIndex())
		}
		// Items should only have received FocusMsg/BlurMsg, not the Tab key.
		for i, item := range items {
			for _, msg := range item.messages {
				if _, ok := msg.(tea.KeyMsg); ok {
					t.Fatalf("item %d should not receive tab key message", i)
				}
			}
		}
	})

	t.Run("shift+tab not forwarded to any item", func(t *testing.T) {
		items := []*mockFocusable{{}, {}, {}}
		fm, _ := NewFocusManager(items[0], items[1], items[2])

		// Clear messages from initial focus.
		for _, item := range items {
			item.messages = nil
		}

		fm, _ = fm.Update(shiftTabMsg())

		// Focus should have wrapped to item 2.
		if fm.FocusedIndex() != 2 {
			t.Fatalf("expected index 2, got %d", fm.FocusedIndex())
		}
		// Items should only have received FocusMsg/BlurMsg, not the Shift+Tab key.
		for i, item := range items {
			for _, msg := range item.messages {
				if _, ok := msg.(tea.KeyMsg); ok {
					t.Fatalf("item %d should not receive shift+tab key message", i)
				}
			}
		}
	})
}

// T008: Mouse click forwarding tests.
func TestFocusManagerMouseClickForwarding(t *testing.T) {
	t.Run("click on different bounded item moves focus and forwards click", func(t *testing.T) {
		items := []*mockBoundedFocusable{
			{mockFocusable: mockFocusable{}, x: 0, y: 0, w: 10, h: 1},
			{mockFocusable: mockFocusable{}, x: 0, y: 1, w: 10, h: 1},
		}
		fm, _ := NewFocusManager(items[0], items[1])

		// Clear messages from initial focus.
		for _, item := range items {
			item.messages = nil
		}

		click := mouseClick(5, 1)
		fm, cmd := fm.Update(click)

		// Focus moved to item 1.
		if fm.FocusedIndex() != 1 {
			t.Fatalf("expected index 1, got %d", fm.FocusedIndex())
		}
		// Click was forwarded to the newly focused item (among other messages).
		hasMouseMsg := false
		for _, msg := range items[1].messages {
			if _, ok := msg.(tea.MouseMsg); ok {
				hasMouseMsg = true
			}
		}
		if !hasMouseMsg {
			t.Fatal("expected MouseMsg on item 1")
		}
		// Old focused item should have received BlurMsg but NOT the click.
		hasClick := false
		for _, msg := range items[0].messages {
			if _, ok := msg.(tea.MouseMsg); ok {
				hasClick = true
			}
		}
		if hasClick {
			t.Fatal("item 0 should not receive click")
		}
		// Command should be a batch (focus change + route).
		if cmd == nil {
			t.Fatal("expected batched command")
		}
	})

	t.Run("click on already focused bounded item forwards click", func(t *testing.T) {
		items := []*mockBoundedFocusable{
			{mockFocusable: mockFocusable{}, x: 0, y: 0, w: 10, h: 1},
		}
		fm, _ := NewFocusManager(items[0])

		// Clear messages from initial focus.
		items[0].messages = nil

		click := mouseClick(5, 0)
		fm, _ = fm.Update(click)

		// Focus unchanged.
		if fm.FocusedIndex() != 0 {
			t.Fatalf("expected index 0, got %d", fm.FocusedIndex())
		}
		// Click forwarded to focused item.
		hasMouseMsg := false
		for _, msg := range items[0].messages {
			if _, ok := msg.(tea.MouseMsg); ok {
				hasMouseMsg = true
			}
		}
		if !hasMouseMsg {
			t.Fatal("expected MouseMsg on item 0")
		}
	})

	t.Run("click outside all bounds routes to focused item", func(t *testing.T) {
		items := []*mockBoundedFocusable{
			{mockFocusable: mockFocusable{}, x: 0, y: 0, w: 10, h: 1},
			{mockFocusable: mockFocusable{}, x: 0, y: 1, w: 10, h: 1},
		}
		fm, _ := NewFocusManager(items[0], items[1])

		// Clear messages from initial focus.
		for _, item := range items {
			item.messages = nil
		}

		click := mouseClick(20, 20)
		fm, _ = fm.Update(click)

		// Focus unchanged.
		if fm.FocusedIndex() != 0 {
			t.Fatalf("expected index 0, got %d", fm.FocusedIndex())
		}
		// Click routed to focused item via default path.
		hasMouseMsg := false
		for _, msg := range items[0].messages {
			if _, ok := msg.(tea.MouseMsg); ok {
				hasMouseMsg = true
			}
		}
		if !hasMouseMsg {
			t.Fatal("expected MouseMsg on focused item 0")
		}
		// Unfocused item should not receive the click.
		if len(items[1].messages) != 0 {
			t.Fatalf("item 1 should not receive click, got %d messages", len(items[1].messages))
		}
	})

	t.Run("mouse release forwarded to focused item via default routing", func(t *testing.T) {
		items := []*mockBoundedFocusable{
			{mockFocusable: mockFocusable{}, x: 0, y: 0, w: 10, h: 1},
			{mockFocusable: mockFocusable{}, x: 0, y: 1, w: 10, h: 1},
		}
		fm, _ := NewFocusManager(items[0], items[1])

		// Clear messages from initial focus.
		for _, item := range items {
			item.messages = nil
		}

		releaseMsg := tea.MouseMsg{
			X: 5, Y: 1,
			Action: tea.MouseActionRelease,
			Button: tea.MouseButtonLeft,
		}
		fm, _ = fm.Update(releaseMsg)

		// Focus unchanged (release doesn't trigger focus change).
		if fm.FocusedIndex() != 0 {
			t.Fatalf("expected index 0, got %d", fm.FocusedIndex())
		}
		// Release forwarded to focused item via default routing.
		if len(items[0].messages) != 1 {
			t.Fatalf("expected 1 message on focused item 0, got %d", len(items[0].messages))
		}
		// Unfocused item should not receive the release.
		if len(items[1].messages) != 0 {
			t.Fatalf("item 1 should not receive release, got %d messages", len(items[1].messages))
		}
	})
}
