package ui

import (
	//table "github.com/calyptia/go-bubble-table"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var (
	style = lipgloss.NewStyle().Padding(1)
)

type Table struct {
	table *table.Table
}

func NewTable(header []string, data [][]interface{}, w, h int) (Table, error) {
	var rows [][]string
	for _, d := range data {
		var temp []string
		for _, k := range d {
			temp = append(temp, fmt.Sprintf("%+v", k))
		}
		rows = append(rows, temp)
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("FF"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 0:
				return lipgloss.NewStyle().Bold(true)
			}
			return lipgloss.NewStyle()
		}).
		Headers(header...).
		Rows(rows...)

	return Table{table: t}, nil
}

func (t Table) Init() tea.Cmd { return nil }

func (t Table) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return nil, nil }

func (t Table) View() string {
	return style.Render(t.table.String())
}
