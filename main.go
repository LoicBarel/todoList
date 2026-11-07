package main

import (
	"fmt"
	"os"
	"slices"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	//"log"
)

type model struct {
	choices []string
	cursor int
	selected map[int]struct{}
	todoTextarea textarea.Model 
	errorMsg string
	editSelected int
}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Saisir une tâche..."
	ta.ShowLineNumbers = false
	ta.SetHeight(1)
	//ta.Focus()
	return model {
		choices: []string{"Buy milk", "Buy coffee", "Change grinder's settings"},
		selected: make(map[int]struct{}),
		todoTextarea: ta,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		m.errorMsg = ""
		if m.todoTextarea.Focused() {
			switch msg.String() {
			case "tab":
				if !m.todoTextarea.Focused() {
					m.todoTextarea.Focus()
				} else {
					m.todoTextarea.Blur()
				}
			case "enter":
				text := m.todoTextarea.Value()
				if slices.Contains(m.choices, text) {
					m.errorMsg = "The task already exist\n\n"
				} else {
					m.choices = append(m.choices, text)
				}
				m.todoTextarea.Reset()
				m.todoTextarea.Blur()
			}
		} else {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			case "enter", "space":
				if _, exist := m.selected[m.cursor]; exist {
					delete(m.selected, m.cursor)
					break
				}
				m.selected[m.cursor] = struct{}{}	
			case "tab":
				m.todoTextarea.Focus()
			case "e":
				m.editSelected = m.cursor
				m.todoTextarea.Focus()
			}
		}
	}
	m.todoTextarea, cmd = m.todoTextarea.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	s := "What should I do ?\n\n"
	s += fmt.Sprintf("%s\n\n",m.todoTextarea.View())
	s += m.errorMsg
	for index, task := range m.choices {

		var cursor, selected string

		if m.cursor == index {
			cursor = ">"
		}

		if _, exist := m.selected[index]; exist {
			selected = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, selected, task)

	}
	s += "\nPress q to quit.\n"
	return tea.NewView(s)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("There's been an error: %v", err)
		os.Exit(1)
	}
}
