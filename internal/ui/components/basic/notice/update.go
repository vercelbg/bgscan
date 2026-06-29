package notice

import (
	"bgscan/internal/ui/shared/ui"

	tea "charm.land/bubbletea/v2"
)

// Update handles incoming BubbleTea messages and updates the notice state.
//
// Behavior:
//   - Enter closes the notice dialog.
//   - All messages are forwarded to the internal viewport to support
//     scrolling for long messages.
func (m *Model) Update(msg tea.Msg) (ui.Component, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.Code == tea.KeyEnter {
			return m, m.CloseCmd()
		}
	}

	// Delegate message handling to the viewport
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)

	return m, cmd
}
