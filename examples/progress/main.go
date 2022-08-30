package main

import (
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zerodoctor/zdtui/ui"
)

func main() {
	var p *tea.Program

	progress := ui.NewProgress(func() (*tea.Program, error) {
		time.Sleep(1 * time.Second)
		p.Send(ui.DefTick(0.3))
		time.Sleep(200 * time.Millisecond)
		p.Send(ui.DefTick(0.1))
		time.Sleep(1 * time.Second)
		p.Send(ui.DefTick(0.1))

		return p, nil
	})

	go progress.Start()

	p = tea.NewProgram(progress)
	if err := p.Start(); err != nil {
		log.Fatalf("failed to start program [error=%s]", err.Error())
	}
}
