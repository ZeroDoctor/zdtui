package ui

import (
	//table "github.com/calyptia/go-bubble-table"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

var (
	style = lipgloss.NewStyle().Padding(1)
)

type Table struct {
	table table.Model
}

func NewTable(header []string, data [][]interface{}, w, h int) (Table, error) {
	t := Table{}

	top, right, bottom, left := style.GetPadding()
	w = w - left - right
	h = h - top - bottom

	// tbl := table.New(header, w, h)
	//
	// var rows []table.Row
	// for _, d := range data {
	// 	var r table.SimpleRow
	// 	t := append(r, append(table.SimpleRow{}, d...))
	// 	rows = append(rows, t)
	// }
	// tbl.SetRows(rows)

	//t.table = tbl

	var maxWidth int

	var cols []table.Column
	for _, str := range header {
		cols = append(cols, table.NewFlexColumn(str, str, len(header)))
		if lipgloss.Width(str) > maxWidth {
			maxWidth = lipgloss.Width(str)
		}
	}

	var rows []table.Row
	for _, d := range data {
		rowData := make(table.RowData)

		for i, k := range d {
			rowData[header[i]] = k

			kstr := fmt.Sprintf("%+v", k)
			fmt.Println(kstr)
			if lipgloss.Width(kstr) > maxWidth {
				maxWidth = lipgloss.Width(kstr)
			}
		}

		rows = append(rows, table.NewRow(rowData))
	}

	t.table = table.New(cols).WithRows(rows).WithTargetWidth(maxWidth * len(data))

	return t, nil
}

func (t Table) Init() tea.Cmd { return nil }

func (t Table) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return nil, nil }

func (t Table) View() string {
	return style.Render(t.table.View())
}
