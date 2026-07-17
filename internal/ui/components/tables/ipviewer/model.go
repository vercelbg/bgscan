package ipviewer

import (
	"bgscan/internal/core/result"
	"bgscan/internal/ui/components/basic/table"
	"bgscan/internal/ui/shared/env"
	"bgscan/internal/ui/shared/layout"
	"bgscan/internal/ui/shared/ui"

	tea "charm.land/bubbletea/v2"
)

// Model represents the main logs menu component.
type Model struct {
	id     ui.ComponentID
	name   string
	maxRow uint32
	table  ui.Component
	rows   []table.Row
	schema result.ResultSchema
	layout *layout.Layout
}

// ID returns the component's unique identifier.
func (m *Model) ID() ui.ComponentID {
	return m.id
}

// Name returns the component's display name.
func (m *Model) Name() string {
	return m.name
}

// OnClose is called when the component is closed.
func (m *Model) OnClose() tea.Cmd {
	return nil
}

// New creates and returns a new logs menu model.
func New(l *layout.Layout, name string, rows []result.Result, schema result.ResultSchema) *Model {
	cols := make([]table.Column, 0, 2)
	for _, col := range schema.Columns {
		cols = append(cols, table.Column{Title: col.Name, Width: col.Width})
	}
	t := table.New(l, table.WithColumns(cols), table.WithRows([]table.Row{}), table.WithMaxWidth(90))

	m := &Model{
		id:     ui.NewComponentID(),
		name:   name,
		maxRow: 10_000,
		layout: l,
		table:  t,
		schema: schema,
	}

	m.updateRows(rows)
	return m
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) SetRows(rows []result.Result) {
	m.updateRows(rows)
}

func (m *Model) Table() *table.Model {
	if t, ok := m.table.(*table.Model); ok {
		return t
	}
	return nil
}

func (m *Model) Mode() env.Mode {
	return env.NormalMode
}
