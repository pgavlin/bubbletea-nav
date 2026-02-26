// Example focus-routing demonstrates FocusManager message routing
// with a simple form containing three text input fields. Tab and
// Shift+Tab cycle focus; typed characters go to the focused field.
package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
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
	fm, _ := nav.NewFocusManager(fields[0], fields[1], fields[2])
	return formModel{
		focus:  fm,
		fields: fields,
	}
}

func (m formModel) Init() tea.Cmd { return nil }

func (m formModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle quit keys before delegating to FocusManager.
	if msg, ok := msg.(tea.KeyPressMsg); ok {
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.focus, cmd = m.focus.Update(msg)
	return m, cmd
}

func (m formModel) View() tea.View {
	var b strings.Builder
	b.WriteString("Focus Routing Demo\n\n")
	for _, f := range m.fields {
		b.WriteString(f.View().Content + "\n")
	}
	b.WriteString("\ntab: next field | shift+tab: prev field | esc: quit")
	return tea.NewView(b.String())
}

// -- Main --

func main() {
	p := tea.NewProgram(newFormModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
