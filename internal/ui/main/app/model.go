package app

import (
	"bgscan/internal/ui/main/body"
	"bgscan/internal/ui/main/footer"
	"bgscan/internal/ui/main/header"
	"bgscan/internal/ui/shared/dialog"
	"bgscan/internal/ui/shared/layout"
	"bgscan/internal/ui/shared/ui"

	tea "charm.land/bubbletea/v2"
)

// dialogPosition defines how an overlay component
// should be positioned relative to the terminal layout.
type dialogPosition struct {
	XPos    dialog.DialogPosition
	YPos    dialog.DialogPosition
	XOffset int
	YOffset int
}

// model is the root BubbleTea application model.
// It coordinates layout updates, base UI components,
// and overlay layers.
type model struct {
	layout *layout.Layout

	// dialog contains active overlay components.
	dialog []ui.Component

	// dialogPlacements stores the position metadata
	// for each overlay component.
	dialogPlacements map[ui.ComponentID]*dialogPosition

	header ui.Component
	body   ui.Component
	footer ui.Component
}

// New creates and initializes the root application model.
func New() tea.Model {
	l := layout.New()

	return &model{
		layout:           l,
		dialog:           make([]ui.Component, 0, 5),
		dialogPlacements: make(map[ui.ComponentID]*dialogPosition),
		header:           header.New(l),
		body:             body.New(l),
		footer:           footer.New(l),
	}
}

// Init initializes all base UI components.
func (m *model) Init() tea.Cmd {
	return tea.Batch(
		m.header.Init(),
		m.body.Init(),
		m.footer.Init(),
	)
}

// getDialogPlacement returns the placement configuration
// for an overlay component. If no placement exists,
// a default centered placement is created and stored.
func (m *model) getDialogPlacement(id ui.ComponentID) *dialogPosition {
	if p, ok := m.dialogPlacements[id]; ok {
		return p
	}

	p := &dialogPosition{
		XPos:    dialog.Center,
		YPos:    dialog.Center,
		XOffset: 0,
		YOffset: 0,
	}

	m.dialogPlacements[id] = p
	return p
}
