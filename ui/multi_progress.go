package ui

import (
	"context"
	"strings"

	"github.com/alitto/pond"
	tea "github.com/charmbracelet/bubbletea"
)

type MultiProgressOption func(mp *MultiProgress)

func MultiProgWorkerOption(workers int) MultiProgressOption {
	return func(mp *MultiProgress) {
		mp.workers = workers
	}
}

func MultiProgPondOptions(opts ...pond.Option) MultiProgressOption {
	return func(mp *MultiProgress) {
		mp.pondOpt = append(mp.pondOpt, opts...)
	}
}

func MultiProgMaxCapicity(maxCap int) MultiProgressOption {
	return func(mp *MultiProgress) {
		mp.maxCap = maxCap
	}
}

type MultiProgress struct {
	workload []ProgressWork
	workers  int
	maxCap   int
	bars     []*Progress
	err      error
	pondOpt  []pond.Option
	pool     *pond.TaskGroupWithContext
	ctx      context.Context
}

func NewMultiProgress(ctx context.Context, workload []ProgressWork, opts ...MultiProgressOption) (*MultiProgress, context.Context) {
	m := &MultiProgress{
		workers:  5,
		maxCap:   0,
		workload: workload,
		ctx:      ctx,
	}

	for _, opt := range opts {
		opt(m)
	}

	pd := pond.New(m.workers, m.maxCap, m.pondOpt...)
	m.pool, m.ctx = pd.GroupContext(m.ctx)

	return m, m.ctx
}

func (m *MultiProgress) Init() tea.Cmd {
	for i := range m.workload {
		n := i
		m.pool.Submit(func() error {
			progress := NewProgress(m.workload[n])
			m.bars = append(m.bars, progress)

			return nil
		})
	}

	return func() tea.Msg {

		return nil
	}
}

func (m *MultiProgress) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		for i := range msg.Runes {
			switch msg.Runes[i] {
			case 'q':
				return m, tea.Quit
			}
		}
	}

	cmds := make([]tea.Cmd, len(m.bars))

	// for i := range m.bars {
	// 	var model tea.Model
	// 	model, cmds[i] = m.bars[i].Progress.Update(msg)
	// 	m.bars[i].Progress = model.(progress.Model)
	// }

	for i := range m.bars {
		var model tea.Model
		model, cmds[i] = m.bars[i].Update(msg)
		m.bars[i] = model.(*Progress)
	}

	return m, tea.Batch(cmds...)
}

func (m *MultiProgress) View() string {
	var builder strings.Builder

	for i := range m.bars {
		builder.WriteString(m.bars[i].View())

		if i < len(m.bars)-1 {
			builder.WriteRune('\n')
		}
	}

	return builder.String()
}
