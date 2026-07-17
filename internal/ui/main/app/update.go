package app

import (
	"fmt"
	"runtime"
	"strings"

	"bgscan/internal/logger"
	"bgscan/internal/ui/shared/dialog"
	"bgscan/internal/ui/shared/env"
	"bgscan/internal/ui/shared/ui"

	tea "charm.land/bubbletea/v2"
)

// Update is the central message router for the application.
// It processes BubbleTea messages, manages overlay layers,
// and dispatches updates to UI components.
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	// Handle terminal resize
	case tea.WindowSizeMsg:
		m.layout.Update(msg.Width, msg.Height)

	// Handle keyboard input
	case tea.KeyPressMsg:

		// Immediate application quit
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// dump Goroutine for DebugInfo
		if msg.String() == env.KeyCtrlT {
			logger.DebugInfo("%s", dumpGoroutines())
		}

		// Overlay back/quit handling
		if len(m.dialog) > 0 {
			lastIdx := len(m.dialog) - 1
			top := m.dialog[lastIdx]

			if env.IsBackKey(msg, top.Mode()) || env.IsQuitKey(msg, top.Mode()) {

				// Execute overlay cleanup command
				cmds = append(cmds, top.OnClose())

				// Remove overlay placement metadata
				delete(m.dialogPlacements, top.ID())

				// Remove overlay from stack
				m.dialog[lastIdx] = nil
				m.dialog = m.dialog[:lastIdx]

				return m, tea.Batch(cmds...)
			}
		}

	// Add a new overlay component
	case dialog.OpenDialogMsg:
		m.dialog = append(m.dialog, msg.Component)

		m.dialogPlacements[msg.Component.ID()] = &dialogPosition{
			XPos:    msg.XPos,
			YPos:    msg.YPos,
			XOffset: msg.XOffset,
			YOffset: msg.YOffset,
		}

		return m, msg.Component.Init()

	// Close an existing overlay component
	case ui.CloseComponentMsg:
		for i, ov := range m.dialog {
			if ov.ID() == msg.ID {

				cmds = append(cmds, ov.OnClose())

				// Remove overlay safely from slice
				m.dialog = append(m.dialog[:i], m.dialog[i+1:]...)

				// Remove placement metadata
				delete(m.dialogPlacements, msg.ID)

				break
			}
		}
	}

	// --- Overlay Input Routing ---

	// If overlays exist, the top overlay consumes all input.
	if len(m.dialog) > 0 {
		lastIdx := len(m.dialog) - 1

		newLayer, cmd := m.dialog[lastIdx].Update(msg)
		m.dialog[lastIdx] = newLayer

		// Block key input from reaching background components
		if _, ok := msg.(tea.KeyMsg); ok {
			return m, cmd
		}

		cmds = append(cmds, cmd)
	}

	// --- Background Component Updates ---

	var hCmd, bCmd, fCmd tea.Cmd

	m.header, hCmd = m.header.Update(msg)
	m.body, bCmd = m.body.Update(msg)
	m.footer, fCmd = m.footer.Update(msg)

	cmds = append(cmds, hCmd, bCmd, fCmd)

	return m, tea.Batch(cmds...)
}

// dumpGoroutines captures all goroutines in suspicious wait states and returns
// the full formatted dump as a string so callers can log, send, or write it.
//
// Each entry in the returned string is a complete goroutine block, e.g.:
//
//	goroutine 42 [chan receive, 3 minutes]:
//	net/http.(*persistConn).readLoop(...)
//	    /usr/local/go/src/net/http/transport.go:2205
//
// Call this before and after a probe run and diff the output to find leaks.
func dumpGoroutines() string {
	buf := make([]byte, 1<<20)
	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			buf = buf[:n]
			break
		}
		buf = make([]byte, len(buf)*2)
	}
	suspiciousStates := []string{
		"[chan receive",
		"[chan send",
		"[select",
		"[IO wait",
		"[sleep",
		"[semacquire",
		"[sync.Mutex.Lock",
	}

	blocks := strings.Split(string(buf), "\n\n")
	var matched []string
	for _, block := range blocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}
		for _, state := range suspiciousStates {
			if strings.Contains(block, state) {
				matched = append(matched, block)
				break
			}
		}
	}

	total := runtime.NumGoroutine()
	header := fmt.Sprintf("=== Goroutine Dump: %d total, %d suspicious ===\n", total, len(matched))

	if len(matched) == 0 {
		return header + "(none)\n"
	}

	return header + strings.Join(matched, "\n\n") + "\n"
}
