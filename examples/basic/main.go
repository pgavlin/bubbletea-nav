// Example basic demonstrates a two-screen push/pop application
// with focus management. The home screen has a list of items;
// pressing Enter pushes a detail screen. Pressing Escape pops back.
package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	nav "github.com/pgavlin/bubbletea-nav"
)

// -- Home Screen --

type homeScreen struct {
	items    []string
	selected int
}

func newHomeScreen() homeScreen {
	return homeScreen{
		items: []string{"Alpha", "Bravo", "Charlie"},
	}
}

func (s homeScreen) Init() tea.Cmd { return nil }

func (s homeScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.Code {
		case tea.KeyUp:
			if s.selected > 0 {
				s.selected--
			}
		case tea.KeyDown:
			if s.selected < len(s.items)-1 {
				s.selected++
			}
		case tea.KeyEnter:
			return s, nav.Push(newDetailScreen(s.items[s.selected]))
		}
		if msg.String() == "ctrl+c" {
			return s, tea.Quit
		}
	}
	return s, nil
}

func (s homeScreen) View() tea.View {
	var b strings.Builder
	b.WriteString("Home Screen\n\n")
	for i, item := range s.items {
		cursor := "  "
		if i == s.selected {
			cursor = "> "
		}
		b.WriteString(cursor + item + "\n")
	}
	b.WriteString("\nEnter: view detail | ctrl+c: quit")
	return tea.NewView(b.String())
}

// -- Detail Screen --

type detailScreen struct {
	item string
}

func newDetailScreen(item string) detailScreen {
	return detailScreen{item: item}
}

func (s detailScreen) Init() tea.Cmd { return nil }

func (s detailScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.Code {
		case tea.KeyEscape:
			return s, nav.Pop()
		}
		if msg.String() == "ctrl+c" {
			return s, tea.Quit
		}
	}
	return s, nil
}

func (s detailScreen) View() tea.View {
	return tea.NewView(fmt.Sprintf("Detail Screen\n\nViewing: %s\n\nEsc: go back | ctrl+c: quit", s.item))
}

// -- Main --

func main() {
	stack := nav.NewStack(newHomeScreen())
	p := tea.NewProgram(stack)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
