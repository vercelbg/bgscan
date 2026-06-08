package crud

func (m *Model[T]) View() string {
	return m.table.View()
}
