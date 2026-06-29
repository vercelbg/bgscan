package table

import (
	"bgscan/internal/logger"

	"charm.land/lipgloss/v2"
)

func (m *Model) View() string {
	width := m.Layout.Body.Width
	logger.DebugInfo("View: rows=%d cols=%v", len(m.BubbleTable.Rows()), m.BubbleTable.Columns())
	tableView := tableViewStyle(width).Render(m.BubbleTable.View())

	raw := m.BubbleTable.View()
	logger.DebugInfo("View raw (%d bytes): %q", len(raw), raw)

	return lipgloss.NewStyle().
		Width(width).
		Render(lipgloss.JoinVertical(
			lipgloss.Center,
			m.renderTitle(),
			tableView,
			m.renderHelpView(),
		))
}

func (m *Model) renderHelpView() string {
	width := m.Layout.Body.Width
	helpView := ""
	if m.FullHelp {
		helpView = helpViewStyle(width).
			Render(
				m.Help.FullHelpView(m.Keys.FullHelp()),
			)
	} else {
		helpView = helpViewStyle(width).
			Render(
				m.Help.ShortHelpView(m.Keys.ShortHelp()),
			)
	}
	return helpView
}

func (m *Model) renderTitle() string {
	width := m.Layout.Body.Width
	title := ""
	if m.Title != "" {
		title = titleStyles(width).Render(m.Title)
	}
	return title
}
