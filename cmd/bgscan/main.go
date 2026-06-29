package main

import (
	"fmt"
	"os"

	"bgscan/internal/logger"
	"bgscan/internal/startup"
	"bgscan/internal/ui/main/app"

	tea "charm.land/bubbletea/v2"
)

func main() {
	startup.RunHealthChecks()

	defer logger.CloseAll()

	p := tea.NewProgram(app.New())

	if _, err := p.Run(); err != nil {
		fmt.Printf("BubbleTea runtime error:%s", err.Error())
		os.Exit(1)
	}
}
