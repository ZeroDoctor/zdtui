package ui

import (
	"context"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type MultiProgressOption func(mp *MultiProgressBar)

type MPQuit bool

type MultiProgressBar struct {
	workload          []ProgressWork
	bars              []*ProgressBar
	err               error
	ctx               context.Context // not really used
	totalFinishedJobs int
}

func NewMultiProgress(ctx context.Context, workload []ProgressWork, opts ...MultiProgressOption) (*MultiProgressBar, context.Context) {
	m := &MultiProgressBar{
		workload: workload,
		ctx:      ctx,
	}

	for _, opt := range opts {
		opt(m)
	}

	for i := range workload {
		m.bars = append(m.bars, NewProgress(
			m.ctx,
			workload[i],
			ProgSetID(i),
			ProgDontQuitOnErr(),
		))
	}

	return m, m.ctx
}

func (m *MultiProgressBar) Init() tea.Cmd {
	var cmds []tea.Cmd
	for i := range m.bars {
		n := i
		cmds = append(cmds,
			m.bars[n].Init(),
			func() tea.Msg {
				return m.bars[n].Start()
			},
		)
	}

	return tea.Batch(cmds...)
}

func (m *MultiProgressBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		for i := range msg.Runes {
			switch msg.Runes[i] {
			case 'q':
				return m, tea.Quit
			}
		}
	case ProgMsg:
		if msg.id == -1 {
			break
		}

		var cmds []tea.Cmd

		model, cmd := m.bars[msg.id].Update(msg)
		cmds = append(cmds, cmd)
		m.bars[msg.id] = model.(*ProgressBar)

		return m, tea.Batch(cmds...)

	case MPQuit:
		m.totalFinishedJobs++
		if m.totalFinishedJobs >= len(m.bars) {

			return m, tea.Sequentially(
				tea.Tick(1*time.Second, func(t time.Time) tea.Msg { return nil }),
				tea.Quit,
			)
		}

		return m, nil
	}

	cmds := make([]tea.Cmd, len(m.bars))
	for i := range m.bars {
		var model tea.Model
		model, cmds[i] = m.bars[i].Update(msg)
		m.bars[i] = model.(*ProgressBar)
	}

	return m, tea.Batch(cmds...)
}

func (m *MultiProgressBar) View() string {
	var builder strings.Builder

	for i := range m.bars {
		builder.WriteString(m.bars[i].View() + "\n")
	}

	return builder.String()
}
