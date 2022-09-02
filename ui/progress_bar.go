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

var EOF_ProgErr DefProgErr = errors.New("EOF")
var Cancel_ProgErr DefProgErr = errors.New("progress was cancelled")

type DefProgErr error

type ProgressWork func(p *Progress) error

type ProgressOption func(*Progress)

func ProgEnableFocus() ProgressOption {
	return func(p *Progress) {
		p.enableFocus = true
	}
}

func ProgDontQuitOnErr() ProgressOption {
	return func(p *Progress) {
		p.quitOnErr = false
	}
}

func ProgSetID(id int) ProgressOption {
	return func(p *Progress) {
		p.id = id
	}
}

type ProgMsg struct {
	id        int
	amount    float64
	nextFrame tea.Cmd
	err       DefProgErr
}

type Progress struct {
	id          int
	enableFocus bool
	quitOnErr   bool
	focus       bool
	finished    bool

	work       ProgressWork
	msgChan    chan tea.Msg
	displayMsg string
	err        error

	Progress progress.Model
}

func NewProgress(work ProgressWork, opts ...ProgressOption) *Progress {
	prog := &Progress{
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

func (p *Progress) Start() tea.Msg {
	if err := p.work(p); err != nil {
		return ProgMsg{err: DefProgErr(err)}
	}

	return ProgMsg{err: DefProgErr(EOF_ProgErr)}
}

func (p *Progress) SendTick(tick float64) {
	p.msgChan <- ProgMsg{id: p.id, amount: tick}
}

func (p *Progress) waitForActivity() tea.Cmd {
	return func() tea.Msg { return <-p.msgChan }
}

func (p *Progress) Focus(b bool) tea.Cmd {
	p.focus = b
	return func() tea.Msg {
		return nil
	}
}

func (p *Progress) Init() tea.Cmd {
	return p.waitForActivity()
}

func (p *Progress) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

		if msg.err != nil {
			if p.quitOnErr {
				cmds = append(cmds, tea.Sequentially(
					tea.Tick(time.Millisecond*750, func(_ time.Time) tea.Msg { // pause a bit before quiting
						return nil
					}),
					tea.Quit,
				))
			}

			if msg.err == EOF_ProgErr {
				p.displayMsg = "done!"

				if !p.quitOnErr {
					return p, tea.Batch(p.Progress.IncrPercent(1.0), p.waitForActivity())
				}

				cmds = append(cmds, p.Progress.IncrPercent(1.0))
				return p, tea.Batch(cmds...)
			}

			p.err = msg.err

			return p, tea.Batch(cmds...)
		}

		return p, tea.Batch(
			p.Progress.IncrPercent(float64(msg.amount)),
			p.waitForActivity(),
		)

	case progress.FrameMsg:
		pm, cmd := p.Progress.Update(msg)
		p.Progress = pm.(progress.Model)

		return p, cmd
	}

	return p, nil
}

func (p *Progress) View() string {
	fn := func() string { return defStyle(p.displayMsg) }
	if p.err != nil {
		fn = func() string { return errStyle(p.err.Error()) }
	}

	pad := strings.Repeat(" ", PADDING)
	return pad + fn() + "\n" +
		pad + p.Progress.View() + "\n"
}
