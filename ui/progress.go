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

type ProgressWork func() (*tea.Program, error)

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

type DefTick float64
type DefProgErr error

type Progress struct {
	Progress    progress.Model
	enableFocus bool
	focus       bool
	quitOnErr   bool
	err         error
}

func NewProgress(work ProgressWork, opts ...ProgressOption) *Progress {
	prog := &Progress{
		quitOnErr: true,
		Progress: progress.New(
			progress.WithScaledGradient(GREEN_COL, BLUE_COL),
			progress.WithColorProfile(termenv.ANSI256),
		),
	}

	for _, opt := range opts {
		opt(prog)
	}

	go func() {
		var p *tea.Program
		var err error

		if p, err = work(); err != nil {
			p.Send(DefProgErr(err))
			return
		}

		p.Send(DefProgErr(EOF_ProgErr))
	}()

	return prog
}

func (p *Progress) Focus(b bool) tea.Cmd {
	p.focus = b
	return func() tea.Msg {
		return nil
	}
}

func (p *Progress) Init() tea.Cmd { return nil }

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

	case DefProgErr:
		var quitCmd tea.Cmd
		if p.quitOnErr {
			quitCmd = tea.Sequentially(
				tea.Tick(time.Millisecond*750, func(_ time.Time) tea.Msg { // pause a bit before quiting
					return nil
				}),
				tea.Quit,
			)
		}

		if msg == EOF_ProgErr {
			return p, tea.Batch(quitCmd, p.Progress.IncrPercent(1.0))
		}

		p.err = msg

		return p, quitCmd

	case DefTick:
		return p, p.Progress.IncrPercent(float64(msg))

	case progress.FrameMsg:
		pm, cmd := p.Progress.Update(msg)
		p.Progress = pm.(progress.Model)

		return p, cmd
	}

	return p, nil
}

func (p *Progress) View() string {
	fn := func() string { return defStyle("working...") }
	if p.err != nil {
		fn = func() string { return errStyle(p.err.Error()) }
	}

	pad := strings.Repeat(" ", PADDING)
	return pad + fn() + "\n" +
		pad + p.Progress.View() + "\n"
}
