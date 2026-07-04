package iplist

import (
	"bgscan/internal/core/iplist"
	"bgscan/internal/ui/components/basic/crud"
	"bgscan/internal/ui/components/basic/picker"
	"bgscan/internal/ui/shared/ui"

	tea "charm.land/bubbletea/v2"
)

func (m *Model) Update(msg tea.Msg) (ui.Component, tea.Cmd) {
	switch msg := msg.(type) {
	case crud.MsgActionTrigger:
		if msg.ActionType == "add" {
			return m, picker.OpenFilePickerCmd(
				m.layout,
				"Select IP File (.txt)",
				"",
				[]string{".txt"},
				m.handleFileSelect,
			)
		}
	}

	updated, cmd := m.crudTable.Update(msg)

	if table, ok := updated.(*crud.Model[iplist.IPFileInfo]); ok {
		m.crudTable = table
	}

	return m, cmd
}
