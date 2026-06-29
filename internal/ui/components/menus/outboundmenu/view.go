package outboundmenu

import "charm.land/lipgloss/v2"

func (m *Model) View() string {
	return lipgloss.NewStyle().Padding(0, 5).Render(m.menu.View())
}
