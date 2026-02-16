package nav

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Bounded is an optional interface for models that support mouse-click
// targeting. Implement in addition to tea.Model.
type Bounded interface {
	// Bounds returns the component's bounding rectangle in terminal
	// coordinates (zero-based, top-left origin).
	Bounds() (x, y, width, height int)
}

// FocusMsg is sent to a model's Update method when it gains keyboard
// focus. Handle it in Update to activate input state (e.g., show cursor).
type FocusMsg struct{}

// BlurMsg is sent to a model's Update method when it loses keyboard
// focus. Handle it in Update to deactivate input state (e.g., hide cursor).
type BlurMsg struct{}

// FocusChangedMsg is emitted by FocusManager when focus changes.
type FocusChangedMsg struct {
	Previous int // index of previously focused item (-1 if none)
	Current  int // index of newly focused item (-1 if none)
}

// FocusManager manages sequential focus order within a screen.
// It is a value type — screens hold it as a field and delegate
// message handling to it.
type FocusManager struct {
	items      []tea.Model
	focusIndex int
}

// NewFocusManager creates a focus manager with the given items.
// The first item receives a FocusMsg through its Update method.
// If items is empty, no item is focused. Returns any command
// produced by the initial focus delivery.
func NewFocusManager(items ...tea.Model) (FocusManager, tea.Cmd) {
	fm := FocusManager{
		items:      items,
		focusIndex: -1,
	}
	if len(items) > 0 {
		fm.focusIndex = 0
		updated, cmd := items[0].Update(FocusMsg{})
		fm.items[0] = updated
		return fm, cmd
	}
	return fm, nil
}

// Update processes a message. Tab and Shift+Tab cycle focus and are
// not forwarded. Mouse clicks on bounded items move focus then
// forward the click. All other messages are forwarded to the focused
// item via its Update method.
func (fm FocusManager) Update(msg tea.Msg) (FocusManager, tea.Cmd) {
	if len(fm.items) == 0 {
		return fm, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			return fm.advanceFocus(1)
		case tea.KeyShiftTab:
			return fm.advanceFocus(-1)
		}

	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress {
			for i, item := range fm.items {
				b, ok := item.(Bounded)
				if !ok {
					continue
				}
				bx, by, bw, bh := b.Bounds()
				if msg.X >= bx && msg.X < bx+bw && msg.Y >= by && msg.Y < by+bh {
					if i != fm.focusIndex {
						var focusCmd tea.Cmd
						fm, focusCmd = fm.setFocus(i)
						routeCmd := fm.routeMessage(msg)
						return fm, tea.Batch(focusCmd, routeCmd)
					}
					return fm, fm.routeMessage(msg)
				}
			}
		}
	}

	return fm, fm.routeMessage(msg)
}

// routeMessage forwards a message to the focused item via its
// Update method. Returns the command from the item.
func (fm FocusManager) routeMessage(msg tea.Msg) tea.Cmd {
	if fm.focusIndex < 0 || fm.focusIndex >= len(fm.items) {
		return nil
	}
	updated, cmd := fm.items[fm.focusIndex].Update(msg)
	fm.items[fm.focusIndex] = updated
	return cmd
}

// advanceFocus moves focus by delta (1 for forward, -1 for backward)
// with cyclic wrapping.
func (fm FocusManager) advanceFocus(delta int) (FocusManager, tea.Cmd) {
	prev := fm.focusIndex
	n := len(fm.items)
	next := (prev + delta + n) % n
	return fm.setFocusWithPrev(prev, next)
}

// setFocus moves focus to the given index, blurring the old item.
func (fm FocusManager) setFocus(index int) (FocusManager, tea.Cmd) {
	return fm.setFocusWithPrev(fm.focusIndex, index)
}

// setFocusWithPrev performs the blur/focus transition via BlurMsg and
// FocusMsg messages, and emits FocusChangedMsg. Blur is delivered
// before focus per FR-010.
func (fm FocusManager) setFocusWithPrev(prev, next int) (FocusManager, tea.Cmd) {
	var cmds []tea.Cmd

	// Deliver BlurMsg to the old item.
	if prev >= 0 && prev < len(fm.items) {
		updated, cmd := fm.items[prev].Update(BlurMsg{})
		fm.items[prev] = updated
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	fm.focusIndex = next

	// Deliver FocusMsg to the new item.
	if next >= 0 && next < len(fm.items) {
		updated, cmd := fm.items[next].Update(FocusMsg{})
		fm.items[next] = updated
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	changeMsg := FocusChangedMsg{Previous: prev, Current: next}
	cmds = append(cmds, func() tea.Msg { return changeMsg })

	return fm, tea.Batch(cmds...)
}

// SetItems replaces the items list. The previously focused item
// receives BlurMsg and the first new item receives FocusMsg, both
// via Update. Focus resets to the first item (or no focus if empty).
func (fm FocusManager) SetItems(items ...tea.Model) (FocusManager, tea.Cmd) {
	var cmds []tea.Cmd

	// Blur the currently focused item.
	if fm.focusIndex >= 0 && fm.focusIndex < len(fm.items) {
		updated, cmd := fm.items[fm.focusIndex].Update(BlurMsg{})
		fm.items[fm.focusIndex] = updated
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	fm.items = items
	fm.focusIndex = -1
	if len(items) > 0 {
		fm.focusIndex = 0
		updated, cmd := items[0].Update(FocusMsg{})
		fm.items[0] = updated
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return fm, tea.Batch(cmds...)
}

// FocusedIndex returns the index of the currently focused item,
// or -1 if no item is focused.
func (fm FocusManager) FocusedIndex() int {
	return fm.focusIndex
}

// FocusIndex sets focus to the item at the given index. Index -1
// blurs all items. Out-of-range indices are clamped.
func (fm FocusManager) FocusIndex(index int) (FocusManager, tea.Cmd) {
	if index < -1 {
		index = -1
	}
	if index >= len(fm.items) {
		index = len(fm.items) - 1
	}

	prev := fm.focusIndex
	if prev == index {
		return fm, nil
	}

	if index == -1 {
		// Blur all.
		var cmds []tea.Cmd
		if prev >= 0 && prev < len(fm.items) {
			updated, cmd := fm.items[prev].Update(BlurMsg{})
			fm.items[prev] = updated
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
		fm.focusIndex = -1
		cmds = append(cmds, func() tea.Msg {
			return FocusChangedMsg{Previous: prev, Current: -1}
		})
		return fm, tea.Batch(cmds...)
	}

	return fm.setFocusWithPrev(prev, index)
}

// String returns a human-readable representation of the focus state.
func (fm FocusManager) String() string {
	return fmt.Sprintf("FocusManager[%d items, focused: %d]", len(fm.items), fm.focusIndex)
}
