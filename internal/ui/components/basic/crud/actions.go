package crud

import (
	"bgscan/internal/ui/components/basic/table"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type MsgActionTrigger struct{ ActionType string }

func (m *Model[T]) configureKeymaps() {
	var keys []table.ActionKey

	// Select (Enter)
	if _, ok := m.provider.OnSelect(m.zeroValue()); ok {
		keys = append(keys, table.NewKey(
			[]string{tea.KeyEnter.String()},
			"select",
			fmt.Sprintf("Select %s", m.name),
			func() tea.Msg { return MsgActionTrigger{ActionType: "select"} },
		))
	}

	// Add (only if supported)
	if m.CanAdd() {
		keys = append(keys, table.NewKey(
			[]string{"a"},
			"add",
			fmt.Sprintf("Add %s", m.name),
			func() tea.Msg { return MsgActionTrigger{ActionType: "add"} },
		))
	}

	// Delete
	if _, ok := m.provider.OnDelete(m.zeroValue()); ok {
		keys = append(keys, table.NewKey(
			[]string{"x"},
			"delete",
			fmt.Sprintf("Delete %s", m.name),
			func() tea.Msg { return MsgActionTrigger{ActionType: "delete"} },
		))
	}

	// Rename
	if _, ok := m.provider.OnRename(m.zeroValue(), ""); ok {
		keys = append(keys, table.NewKey(
			[]string{"r"},
			"rename",
			fmt.Sprintf("Rename %s", m.name),
			func() tea.Msg { return MsgActionTrigger{ActionType: "rename"} },
		))
	}

	m.table.SetKeys(keys...)
}

// Helper to avoid creating zero value repeatedly
func (m *Model[T]) zeroValue() T {
	var zero T
	return zero
}
