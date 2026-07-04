package selectinput

import (
	"bgscan/internal/ui/components/basic/input"

	"charm.land/lipgloss/v2"
)

// View renders the select component.
//
// The view consists of:
//   - An optional title displayed above the select field
//   - The select component (list of options with cursor)
//   - An optional validation error message
//   - Key hints for user interaction
//
// The content is vertically stacked and wrapped inside a styled
// container whose width is determined by the layout.
func (m *Model[T]) View() string {
	content := make([]string, 0, 4)

	// Optional message
	if m.title != "" {
		content = append(content, input.MessageStyle().Render(m.title))
	}

	// Select field
	content = append(content, m.huhInput.View())

	// Validation error (if present)
	if m.errorMsg != "" {
		content = append(content, input.ErrorStyle().Render("✗ "+m.errorMsg))
	}

	// Key hints
	hints := input.KeyHintStyle().Render("↑/↓ to move • Enter to confirm • Esc/b to cancel")
	content = append(content, hints)

	body := lipgloss.JoinVertical(
		lipgloss.Top,
		content...,
	)
	return input.ContainerStyle(m.Width()).Render(body)
}
