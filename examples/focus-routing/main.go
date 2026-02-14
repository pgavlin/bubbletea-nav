// Example focus-routing demonstrates FocusManager message routing
// with a simple form containing three text input fields. Tab and
// Shift+Tab cycle focus; typed characters go to the focused field.
package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	nav "github.com/pgavlin/bubbletea-nav"
)

// -- textField: a Focusable text input --

type textField struct {
	label   string
	value   string
	focused bool
}

func newTextField(label string) *textField {
	return &textField{label: label}
}

func (f *textField) Init() tea.Cmd { return nil }

func (f *textField) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyRunes:
			f.value += string(msg.Runes)
		case tea.KeyBackspace:
			if len(f.value) > 0 {
				f.value = f.value[:len(f.value)-1]
			}
		}
	}
	return f, nil
}

func (f *textField) View() string {
	if f.focused {
		return fmt.Sprintf("> %s: [%s]", f.label, f.value)
	}
	return fmt.Sprintf("  %s:  %s", f.label, f.value)
}

func (f *textField) Focus() tea.Cmd { f.focused = true; return nil }
func (f *textField) Blur()          { f.focused = false }

// -- formModel: top-level tea.Model --

type formModel struct {
	focus  nav.FocusManager
	fields []*textField
}

func newFormModel() formModel {
	fields := []*textField{
		newTextField("Name"),
		newTextField("Email"),
		newTextField("City"),
	}
	return formModel{
		focus:  nav.NewFocusManager(fields[0], fields[1], fields[2]),
		fields: fields,
	}
}

func (m formModel) Init() tea.Cmd { return nil }

func (m formModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle quit keys before delegating to FocusManager.
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEscape:
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.focus, cmd = m.focus.Update(msg)
	return m, cmd
}

func (m formModel) View() string {
	var b strings.Builder
	b.WriteString("Focus Routing Demo\n\n")
	for _, f := range m.fields {
		b.WriteString(f.View() + "\n")
	}
	b.WriteString("\ntab: next field | shift+tab: prev field | esc: quit")
	return b.String()
}

// -- Main --

func main() {
	p := tea.NewProgram(newFormModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
