package ui

import (
	"errors"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

const (
	PADDING   = 2
	MAX_WIDTH = 80
)

const (
	GRAY_COL  = "#D7D7D7"
	GREEN_COL = "#ADFDAD"
	BLUE_COL  = "#AADDFF"
	RED_COL   = "#FF77AA"
)

var (
	defStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(GRAY_COL)).Render
	errStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(RED_COL)).Render
)

type DefProgErr error

var (
	EOF_ProgErr    DefProgErr = errors.New("EOF")
	Cancel_ProgErr DefProgErr = errors.New("progress was cancelled")
)

type DefProgStatus int

const (
	ProgRunning DefProgStatus = iota
	ProgComplete
	ProgFinished
)

type ProgressWork func(p *ProgressBar) error

type ProgressOption func(*ProgressBar)

func ProgEnableFocus() ProgressOption {
	return func(p *ProgressBar) {
		p.enableFocus = true
	}
}

func ProgDontQuitOnErr() ProgressOption {
	return func(p *ProgressBar) {
		p.quitOnErr = false
	}
}

func ProgSetID(id int) ProgressOption {
	return func(p *ProgressBar) {
		p.id = id
	}
}

type ProgMsg struct {
	id        int
	amount    float64
	err       DefProgErr
	nextFrame tea.Cmd
}

type ProgressBar struct {
	id          int
	enableFocus bool
	quitOnErr   bool
	focus       bool
	status      DefProgStatus
	percentage  float64

	work       ProgressWork
	msgChan    chan tea.Msg
	displayMsg string
	err        error

	Progress progress.Model
}

func NewProgress(work ProgressWork, opts ...ProgressOption) *ProgressBar {
	prog := &ProgressBar{
		id:         -1,
		work:       work,
		quitOnErr:  true,
		displayMsg: "working...",
		Progress: progress.New(
			progress.WithScaledGradient(GREEN_COL, BLUE_COL),
			progress.WithColorProfile(termenv.ANSI256),
		),
		msgChan: make(chan tea.Msg, 50),
	}

	for _, opt := range opts {
		opt(prog)
	}

	return prog
}

func (p *ProgressBar) Start() tea.Msg {
	if err := p.work(p); err != nil {
		return ProgMsg{id: p.id, err: DefProgErr(err)}
	}

	return ProgMsg{id: p.id, err: DefProgErr(EOF_ProgErr)}
}

func (p *ProgressBar) SendTick(tick float64) {
	p.msgChan <- ProgMsg{id: p.id, amount: tick}
}

func (p *ProgressBar) waitForActivity() tea.Cmd {
	return func() tea.Msg {
		return <-p.msgChan
	}
}

func (p *ProgressBar) Focus(b bool) tea.Cmd {
	p.focus = b
	return func() tea.Msg {
		return nil
	}
}

func (p *ProgressBar) Init() tea.Cmd {
	return p.waitForActivity()
}

func (p *ProgressBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !p.focus && p.enableFocus {
			return p, nil
		}

		for i := range msg.Runes {
			switch msg.Runes[i] {
			case 'c':
				return p, func() tea.Msg { return Cancel_ProgErr }
			}
		}

	case tea.WindowSizeMsg:
		p.Progress.Width = msg.Width - PADDING*2 - 4
		if p.Progress.Width > MAX_WIDTH {
			p.Progress.Width = MAX_WIDTH
		}

		return p, nil

	case ProgMsg:
		var cmds []tea.Cmd
		if msg.nextFrame != nil {
			cmds = append(cmds, msg.nextFrame)
		}

		if msg.err != nil {
			p.err = msg.err

			if p.quitOnErr {
				cmds = append(cmds, tea.Sequentially(
					tea.Tick(time.Millisecond*750, func(_ time.Time) tea.Msg { // pause a bit before quiting
						return nil
					}),
					tea.Quit,
				))
			}

			if msg.err == EOF_ProgErr {
				p.err = nil
				p.displayMsg = "done!"

				p.percentage = 1.0
			}

			p.status = ProgComplete

			return p, tea.Batch(cmds...)
		}

		p.percentage += msg.amount
		return p, tea.Batch(
			p.waitForActivity(),
		)
	}

	return p, nil
}

func (p *ProgressBar) View() string {
	fn := func() string { return defStyle(p.displayMsg) }
	if p.err != nil {
		fn = func() string { return errStyle(p.err.Error()) }
	}

	if p.status == ProgComplete {
		p.status = ProgFinished
	}

	pad := strings.Repeat(" ", PADDING)
	return pad + fn() + "\n" +
		pad + p.Progress.ViewAs(p.percentage) + "\n"
}
