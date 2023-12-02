package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type TAOption func(*TextArea)

func WithPlaceHolder(placeHolder string) TAOption {
	return func(ta *TextArea) {
		ta.textarea.Placeholder = placeHolder
	}
}

func WithTitle(title string) TAOption {
	return func(ta *TextArea) {
		ta.title = title
	}
}

type TextArea struct {
	title    string
	textarea textarea.Model
	err      error
}

func NewTextArea(options ...TAOption) *TextArea {
	area := textarea.New()
	ta := &TextArea{textarea: area}

	for _, opt := range options {
		opt(ta)
	}

	return ta
}

func (ta *TextArea) Value() string {
	return ta.textarea.Value()
}

func (ta *TextArea) Init() tea.Cmd {
	return textarea.Blink
}

func (ta *TextArea) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return ta, tea.Quit
		default:
			if !ta.textarea.Focused() {
				cmd = ta.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	case error:
		ta.err = msg
		return ta, nil
	}

	ta.textarea, cmd = ta.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return ta, tea.Batch(cmds...)
}

func (ta *TextArea) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		ta.title,
		ta.textarea.View(),
		"(ctrl+c or esc to quit)",
	) + "\n\n"
}
