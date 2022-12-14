package main

import (
	"context"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zerodoctor/zdtui/ui"
)

func work1(msg string) ui.ProgressWork {
	return func(p *ui.ProgressBar) error {
		p.Display(msg)
		time.Sleep(1 * time.Second)
		p.SendTick((0.3))
		time.Sleep(200 * time.Millisecond)
		p.SendTick((0.1))
		time.Sleep(1 * time.Second)
		p.SendTick((0.1))

		return nil
	}
}

func work2() ui.ProgressWork {
	return func(p *ui.ProgressBar) error {
		time.Sleep(1 * time.Second)
		p.SendTick((0.2))
		time.Sleep(200 * time.Millisecond)
		p.SendTick((0.1))
		time.Sleep(1 * time.Second)
		p.Display("working harder...")
		p.SendTick((0.1))
		time.Sleep(1 * time.Second)
		p.SendTick((0.1))
		time.Sleep(1 * time.Second)
		p.SendTick((0.1))

		return nil
	}
}

func main() {
	var p *tea.Program
	var mp *ui.MultiProgressBar

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mp, ctx = ui.NewMultiProgress(ctx, []ui.ProgressWork{work1("this is work 1..."), work2(), work1("this is another work 1...")})

	p = tea.NewProgram(mp)
	if err := p.Start(); err != nil {
		log.Fatalf("failed to start program [error=%s]", err.Error())
	}
}
