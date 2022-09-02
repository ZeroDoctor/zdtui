package main

import (
	"context"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zerodoctor/zdtui/ui"
)

func work() ui.ProgressWork {
	return func(p *ui.ProgressBar) error {
		time.Sleep(1 * time.Second)
		p.SendTick((0.3))
		time.Sleep(200 * time.Millisecond)
		p.SendTick((0.1))
		time.Sleep(1 * time.Second)
		p.SendTick((0.1))

		return nil
	}
}

func main() {
	var p *tea.Program
	progress := ui.NewProgress(context.Background(), work())

	go func() {
		p.Send(progress.Start())
	}()

	p = tea.NewProgram(progress)
	if err := p.Start(); err != nil {
		log.Fatalf("failed to start program [error=%s]", err.Error())
	}
}
