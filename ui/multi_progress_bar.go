package ui

import (
	"context"
	"strings"

	"github.com/alitto/pond"
	tea "github.com/charmbracelet/bubbletea"
)

type MultiProgressOption func(mp *MultiProgressBar)

func MultiProgWorkerOption(workers int) MultiProgressOption {
	return func(mp *MultiProgressBar) {
		mp.workers = workers
	}
}

func MultiProgPondOptions(opts ...pond.Option) MultiProgressOption {
	return func(mp *MultiProgressBar) {
		mp.pondOpt = append(mp.pondOpt, opts...)
	}
}

func MultiProgMaxCapicity(maxCap int) MultiProgressOption {
	return func(mp *MultiProgressBar) {
		mp.maxCap = maxCap
	}
}

type MultiProgressBar struct {
	workload []ProgressWork
	workers  int
	maxCap   int
	bars     []*ProgressBar
	progChan chan *ProgressBar
	err      error
	pondOpt  []pond.Option
	pool     *pond.TaskGroupWithContext
	ctx      context.Context
}

func NewMultiProgress(ctx context.Context, workload []ProgressWork, opts ...MultiProgressOption) (*MultiProgressBar, context.Context) {
	m := &MultiProgressBar{
		workers:  5,
		maxCap:   len(workload),
		workload: workload,
		ctx:      ctx,
		progChan: make(chan *ProgressBar, len(workload)),
	}

	for _, opt := range opts {
		opt(m)
	}

	pd := pond.New(m.workers, m.maxCap, m.pondOpt...)
	m.pool, m.ctx = pd.GroupContext(m.ctx)

	for i := range workload {
		work := workload[i]
		m.bars = append(m.bars, NewProgress(
			work,
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

		model, cmd := m.bars[msg.id].Update(msg)
		m.bars[msg.id] = model.(*ProgressBar)

		return m, cmd
	}

	cmds := make([]tea.Cmd, len(m.bars))
	for i := range m.bars {
		var model tea.Model
		model, cmds[i] = m.bars[i].Update(msg)
		m.bars[i].Progress = model.(*ProgressBar).Progress
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
