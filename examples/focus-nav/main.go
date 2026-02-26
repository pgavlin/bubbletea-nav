// Example focus-nav demonstrates Stack navigation combined with
// FocusManager message routing. A list screen shows contacts;
// pressing Enter pushes an edit screen with focusable text fields.
package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	nav "github.com/pgavlin/bubbletea-nav"
)

// -- textField: handles FocusMsg/BlurMsg for focus state --

type textField struct {
	label   string
	value   string
	focused bool
}

func newTextField(label, value string) *textField {
	return &textField{label: label, value: value}
}

func (f *textField) Init() tea.Cmd { return nil }

func (f *textField) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case nav.FocusMsg:
		f.focused = true
	case nav.BlurMsg:
		f.focused = false
	case tea.KeyPressMsg:
		if msg.Text != "" {
			f.value += msg.Text
		} else if msg.Code == tea.KeyBackspace {
			if len(f.value) > 0 {
				f.value = f.value[:len(f.value)-1]
			}
		}
	}
	return f, nil
}

func (f *textField) View() tea.View {
	if f.focused {
		return tea.NewView(fmt.Sprintf("> %s: [%s]", f.label, f.value))
	}
	return tea.NewView(fmt.Sprintf("  %s:  %s", f.label, f.value))
}

// -- contact data --

type contact struct {
	name, email, phone string
}

var contacts = []contact{
	{"Alice", "alice@example.com", "555-0100"},
	{"Bob", "bob@example.com", "555-0200"},
	{"Charlie", "charlie@example.com", "555-0300"},
}

// -- listScreen: shows contacts, Enter pushes editScreen --

type listScreen struct {
	selected int
}

func newListScreen() listScreen {
	return listScreen{}
}

func (s listScreen) Init() tea.Cmd { return nil }

func (s listScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.Code {
		case tea.KeyUp:
			if s.selected > 0 {
				s.selected--
			}
		case tea.KeyDown:
			if s.selected < len(contacts)-1 {
				s.selected++
			}
		case tea.KeyEnter:
			c := contacts[s.selected]
			return s, nav.Push(newEditScreen(s.selected, c))
		}
		if msg.String() == "ctrl+c" {
			return s, tea.Quit
		}
	}
	return s, nil
}

func (s listScreen) View() tea.View {
	var b strings.Builder
	b.WriteString("Contacts\n\n")
	for i, c := range contacts {
		cursor := "  "
		if i == s.selected {
			cursor = "> "
		}
		b.WriteString(cursor + c.name + "\n")
	}
	b.WriteString("\nenter: edit | ctrl+c: quit")
	return tea.NewView(b.String())
}

// -- editScreen: form with FocusManager --

type editScreen struct {
	index  int
	focus  nav.FocusManager
	fields []*textField
}

func newEditScreen(index int, c contact) editScreen {
	fields := []*textField{
		newTextField("Name", c.name),
		newTextField("Email", c.email),
		newTextField("Phone", c.phone),
	}
	fm, _ := nav.NewFocusManager(fields[0], fields[1], fields[2])
	return editScreen{
		index:  index,
		focus:  fm,
		fields: fields,
	}
}

func (s editScreen) Init() tea.Cmd { return nil }

func (s editScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle screen-level keys before delegating to FocusManager.
	if msg, ok := msg.(tea.KeyPressMsg); ok {
		switch msg.Code {
		case tea.KeyEscape:
			// Save edits back and pop.
			contacts[s.index] = contact{
				name:  s.fields[0].value,
				email: s.fields[1].value,
				phone: s.fields[2].value,
			}
			return s, nav.Pop()
		}
		if msg.String() == "ctrl+c" {
			return s, tea.Quit
		}
	}

	var cmd tea.Cmd
	s.focus, cmd = s.focus.Update(msg)
	return s, cmd
}

func (s editScreen) View() tea.View {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Edit Contact #%d\n\n", s.index+1))
	for _, f := range s.fields {
		b.WriteString(f.View().Content + "\n")
	}
	b.WriteString("\ntab: next field | shift+tab: prev | esc: save & back | ctrl+c: quit")
	return tea.NewView(b.String())
}

// -- Main --

func main() {
	stack := nav.NewStack(newListScreen())
	p := tea.NewProgram(stack)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
