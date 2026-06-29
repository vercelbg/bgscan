package logview

import (
	"context"
	"sync"
	"time"

	"bgscan/internal/core/config"
	"bgscan/internal/logger"
	"bgscan/internal/ui/shared/env"
	"bgscan/internal/ui/shared/layout"
	"bgscan/internal/ui/shared/ui"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Model implements a scrollable log viewer component.
//
// It subscribes to the application's logger and streams log
// messages into a BubbleTea viewport. Messages are buffered
// and updated periodically to avoid excessive UI refreshes.
type Model struct {
	// Component identity
	id    ui.ComponentID
	name  string
	title string

	// Layout
	layout *layout.Layout

	padding           int
	containerMaxWidth int
	containerWidth    int
	showBorder        bool

	// Viewport
	viewport viewport.Model

	// Logger integration
	logger     *logger.Logger
	loggerChan chan string
	maxMessage int

	// Thread‑safe message buffer
	mu         sync.Mutex
	messages   []string
	needUpdate bool

	// Lifecycle management
	ctx    context.Context
	cancel context.CancelFunc
}

// New creates a new log viewer component.
func New(l *layout.Layout, log *logger.Logger, title string) *Model {
	ctx, cancel := context.WithCancel(context.Background())

	m := &Model{
		id:     ui.NewComponentID(),
		name:   title,
		title:  title,
		layout: l,

		logger:     log,
		maxMessage: 200,

		viewport: viewport.New(),

		padding:           5,
		containerMaxWidth: l.Body.Width - 10,
		showBorder:        true,

		ctx:    ctx,
		cancel: cancel,
	}

	m.setSize()

	return m
}

// SetContainerWidth limits the maximum width of the log container.
func (m *Model) SetContainerWidth(width int) {
	m.containerMaxWidth = width
	m.setSize()
}

// SetShowBorder enables or disables the container border.
func (m *Model) SetShowBorder(border bool) {
	m.showBorder = border
	m.setSize()
}

// Init starts the log subscription and background reader.
func (m *Model) Init() tea.Cmd {
	go m.readLogs()
	return m.tick()
}

// readLogs listens for new log messages from the logger.
func (m *Model) readLogs() {
	m.loggerChan = m.logger.Subscribe(200, m.maxMessage)

	for {
		select {

		case <-m.ctx.Done():
			return

		case logMsg, ok := <-m.loggerChan:
			if !ok {
				return
			}

			m.mu.Lock()

			m.messages = append(m.messages, logMsg)

			if len(m.messages) > m.maxMessage {
				m.messages = m.messages[len(m.messages)-m.maxMessage:]
			}

			m.needUpdate = true

			m.mu.Unlock()
		}
	}
}

// tick schedules periodic UI refreshes.
func (m *Model) tick() tea.Cmd {
	return tea.Tick(
		config.Get().General.StatusInterval.Duration(),
		func(time.Time) tea.Msg {
			return LogUpdateTickMsg{}
		},
	)
}

// setSize recalculates the viewport and container dimensions.
func (m *Model) setSize() {
	maxViewportWidth := 80

	m.containerWidth = min(
		m.containerMaxWidth,
		m.layout.Body.Width-10,
	)

	width := min(maxViewportWidth, m.containerWidth-2)

	m.viewport.SetWidth(width)
	helpHeight := lipgloss.Height(
		helpStyle(m.viewport.Width()).Render(helpView()),
	)

	height := m.layout.Body.Height -
		m.padding -
		lipgloss.Height(m.title) -
		helpHeight
	m.viewport.SetHeight(height)
}

// ID returns the component identifier.
func (m *Model) ID() ui.ComponentID {
	return m.id
}

// Name returns the component name.
func (m *Model) Name() string {
	return m.name
}

// Mode defines the input mode used by this component.
func (m *Model) Mode() env.Mode {
	return env.NormalMode
}

// OnClose cleans up resources when the component is removed.
func (m *Model) OnClose() tea.Cmd {
	m.cancel()

	if m.loggerChan != nil {
		m.logger.Unsubscribe(m.loggerChan)
	}

	return nil
}
