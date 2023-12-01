package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Item struct {
	title       string
	desc        string
	filterValue func() string
}

func NewItem(title string, desc string, filterValue func() string) *Item {
	if filterValue == nil {
		filterValue = func() string { return title }
	}

	return &Item{
		title:       title,
		desc:        desc,
		filterValue: filterValue,
	}
}

func (i *Item) FilterValue() string { return i.filterValue() }
func (i *Item) Title() string       { return i.title }
func (i *Item) Description() string { return i.desc }

type List struct {
	List      list.Model
	WasCancel bool
	err       error
}

func NewList(items []list.Item, width, height int) *List {
	l := list.New(items, list.NewDefaultDelegate(), width, height)
	li := &List{List: l}

	return li
}

func (l *List) Init() tea.Cmd { return nil }

func (l *List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		l.List.SetWidth(msg.Width)
		l.List.SetHeight(msg.Height)

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			l.WasCancel = true
			return l, tea.Quit
		case tea.KeyEnter:
			return l, tea.Quit
		}
	case error:
		l.err = msg
		return l, nil
	}

	var cmd tea.Cmd
	l.List, cmd = l.List.Update(msg)

	return l, cmd
}

func (l *List) View() string {
	return fmt.Sprintf("\n%s", l.List.View())
}
