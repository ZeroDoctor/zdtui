package main

import (
	"context"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zerodoctor/zdtui/ui"
)

func work1(p *tea.Program) ui.ProgressWork {
	return func() (*tea.Program, error) {
		for p == nil {
			time.Sleep(500 * time.Millisecond)
		}

		time.Sleep(1 * time.Second)
		p.Send(ui.DefTick(0.3))
		time.Sleep(200 * time.Millisecond)
		p.Send(ui.DefTick(0.1))
		time.Sleep(1 * time.Second)
		p.Send(ui.DefTick(0.1))

		return p, nil
	}
}

func work2(p *tea.Program) ui.ProgressWork {
	return func() (*tea.Program, error) {
		for p == nil {
			time.Sleep(500 * time.Millisecond)
		}

		time.Sleep(1 * time.Second)
		p.Send(ui.DefTick(0.3))
		time.Sleep(200 * time.Millisecond)
		p.Send(ui.DefTick(0.1))
		time.Sleep(1 * time.Second)
		p.Send(ui.DefTick(0.1))
		time.Sleep(1 * time.Second)
		p.Send(ui.DefTick(0.1))
		time.Sleep(1 * time.Second)
		p.Send(ui.DefTick(0.1))

		return p, nil
	}
}

func main() {
	var p *tea.Program
	var mp *ui.MultiProgress

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mp, ctx = ui.NewMultiProgress(ctx, []ui.ProgressWork{work1(p), work2(p), work1(p)})

	p = tea.NewProgram(mp)
	if err := p.Start(); err != nil {
		log.Fatalf("failed to start program [error=%s]", err.Error())
	}
}
