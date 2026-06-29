package table

import (
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"bgscan/internal/logger"
	"bgscan/internal/ui/shared/env"
	"bgscan/internal/ui/shared/layout"
	"bgscan/internal/ui/shared/ui"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// Aliases for upstream table types
type (
	Column = table.Column
	Row    = table.Row
)

// gutterWidth is the horizontal space reserved for borders/padding when
// computing the table's available width.
const gutterWidth = 10

// Model wraps a Bubble Tea table with responsive layout, key bindings, and concurrency safety.
type Model struct {
	mu sync.RWMutex

	id     ui.ComponentID
	name   string
	Title  string
	Layout *layout.Layout

	Help     help.Model
	FullHelp bool

	BubbleTable table.Model
	Keys        KeyMap

	colsWidth []int // original (unscaled) column widths, used as scaling reference
	paddingY  int
}

// Init implements tea.Model.
func (m *Model) Init() tea.Cmd {
	return nil
}

// ID returns the component ID.
func (m *Model) ID() ui.ComponentID {
	return m.id
}

// Name returns the component name.
func (m *Model) Name() string {
	return m.name
}

// OnClose performs cleanup when the component is removed.
func (m *Model) OnClose() tea.Cmd {
	return nil
}

// Mode implements ui.Component.
func (m *Model) Mode() env.Mode {
	return env.NormalMode
}

// New creates a new table model.
func New(title string, cols []table.Column, rows []table.Row, lay *layout.Layout) *Model {
	m := &Model{
		id:     ui.NewComponentID(),
		name:   "table",
		Title:  title,
		Layout: lay,
		Help:   help.New(),
		Keys:   defaultKeys(),
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.colsWidth = columnWidths(cols)
	m.BubbleTable = table.New(
		table.WithColumns(m.scaledColumns(cols)),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(max(1, len(rows))),
		table.WithWidth(m.tableWidthLocked()),
	)
	m.BubbleTable.SetStyles(tableStyles())
	m.updateTableSizeLocked()

	return m
}

// SetPaddingY sets vertical padding and updates table size.
func (m *Model) SetPaddingY(padding int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.paddingY = padding
	m.updateTableSizeLocked()
}

// SetKeys replaces the current key bindings.
func (m *Model) SetKeys(keys ...ActionKey) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Keys = defaultKeys(keys...)
}

// AppendRow adds a new row to the table safely.
func (m *Model) AppendRow(row table.Row) {
	m.mu.Lock()
	defer m.mu.Unlock()

	rows := append(slices.Clone(m.BubbleTable.Rows()), slices.Clone(row))
	m.BubbleTable.SetRows(rows)
	m.updateTableSizeLocked()
}

// SetRow replaces a row at the given index. No-op if index is out of range.
func (m *Model) SetRow(index int, row table.Row) {
	m.mu.Lock()
	defer m.mu.Unlock()

	rows := slices.Clone(m.BubbleTable.Rows())
	if index < 0 || index >= len(rows) {
		return
	}

	rows[index] = slices.Clone(row)
	m.BubbleTable.SetRows(rows)
	m.updateTableSizeLocked()
}

// SetRows replaces all rows in the table.
func (m *Model) SetRows(rows []table.Row) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.BubbleTable.SetRows(cloneRows(rows))
	m.updateTableSizeLocked()
}

// NewRowTime formats a timestamp for display.
func NewRowTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format("2006-01-02 15:04:05")
}

// NewRowBool returns "yes" or "no" for a boolean value.
func NewRowBool(v bool) string {
	if v {
		return "yes"
	}
	return "no"
}

// NewTimeDurationRow returns a human-readable duration from a past time.
func NewTimeDurationRow(from time.Time) string {
	if from.IsZero() {
		return "-"
	}
	d := time.Since(from)
	switch {
	case d < time.Second:
		return "just now"
	case d < time.Minute:
		return d.Truncate(time.Second).String()
	case d < time.Hour:
		return d.Truncate(time.Minute).String()
	default:
		return d.Truncate(time.Hour).String()
	}
}

// cloneRows returns a deep copy of rows, safe to store independently of the caller's slice.
func cloneRows(rows []table.Row) []table.Row {
	out := make([]table.Row, len(rows))
	for i, row := range rows {
		out[i] = slices.Clone(row)
	}
	return out
}

// columnWidths extracts the width of each column, in order.
func columnWidths(cols []table.Column) []int {
	widths := make([]int, len(cols))
	for i, c := range cols {
		widths[i] = c.Width
	}
	return widths
}

// scaledColumns returns a copy of cols with widths scaled to fit the
// available table width, proportional to m.colsWidth. Caller must hold m.mu.
func (m *Model) scaledColumns(cols []table.Column) []table.Column {
	out := slices.Clone(cols)

	total := 0
	for _, w := range m.colsWidth {
		total += w
	}
	if total <= 0 {
		return out
	}

	ratio := float64(m.tableWidthLocked()) / float64(total)
	for i := range out {
		out[i].Width = int(ratio * float64(m.colsWidth[i]))
	}
	return out
}

// updateTableSizeLocked recalculates column widths and table height to fit
// the current layout. Caller must hold m.mu.
func (m *Model) updateTableSizeLocked() {
	if m.Layout == nil || m.Layout.Body.Height == 0 || m.Layout.Body.Width == 0 {
		return
	}
	if len(m.colsWidth) == 0 {
		return
	}

	cols := m.scaledColumns(m.BubbleTable.Columns())
	m.BubbleTable.SetColumns(cols)

	helpHeight := lipgloss.Height(m.renderHelpView())
	titleHeight := lipgloss.Height(m.renderTitle())
	height := max(1, m.Layout.Body.Height-helpHeight-titleHeight-m.paddingY)
	m.BubbleTable.SetHeight(height)
	m.BubbleTable.SetWidth(m.tableWidthLocked())

	logger.DebugInfo("Updated table size: %dx%d", m.tableWidthLocked(), height)
}

// tableWidthLocked returns the available table width given the current
// layout. Caller must hold m.mu (read or write).
func (m *Model) tableWidthLocked() int {
	if m.Layout == nil || m.Layout.Body.Width == 0 {
		return 80
	}
	return min(80, m.Layout.Body.Width-gutterWidth)
}

//
// ──────────────────────────────────────────────────────────────
//  Key bindings
// ──────────────────────────────────────────────────────────────
//

// ActionKey defines a key binding and its action.
type ActionKey struct {
	Keys      []string
	ShortHelp string
	FullHelp  string
	Cmd       tea.Cmd
}

// arrowSymbols maps key names to their display glyphs.
var arrowSymbols = map[string]string{
	"up":    "↑",
	"down":  "↓",
	"left":  "←",
	"right": "→",
}

// NewKey creates a new action key definition.
func NewKey(keys []string, shortHelp, fullHelp string, cmd tea.Cmd) ActionKey {
	ks := make([]string, len(keys))
	for i, k := range keys {
		if symbol, ok := arrowSymbols[k]; ok {
			ks[i] = symbol
		} else {
			ks[i] = k
		}
	}

	kstr := strings.Join(ks, "/")
	if kstr == "" {
		kstr = "?"
	}
	if shortHelp != "" {
		shortHelp = fmt.Sprintf("%s %s", kstr, shortHelp)
	}
	if fullHelp != "" {
		fullHelp = fmt.Sprintf("%s %s", kstr, fullHelp)
	}

	return ActionKey{
		Keys:      keys,
		ShortHelp: shortHelp,
		FullHelp:  fullHelp,
		Cmd:       cmd,
	}
}

// KeyMap stores registered key bindings.
type KeyMap struct {
	Actions []ActionKey
}

// Add appends a new key binding.
func (k *KeyMap) Add(a ActionKey) {
	k.Actions = append(k.Actions, a)
}

// Check returns the command associated with a key message.
func (k KeyMap) Check(msg tea.Msg) tea.Cmd {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil
	}

	keyStr := keyMsg.String()
	for _, a := range k.Actions {
		if slices.Contains(a.Keys, keyStr) {
			return a.Cmd
		}
	}
	return nil
}

// ShortHelp returns key bindings for the condensed help view.
func (k KeyMap) ShortHelp() []key.Binding {
	var bindings []key.Binding
	for _, a := range k.Actions {
		if a.ShortHelp == "" {
			continue
		}
		bindings = append(bindings, key.NewBinding(
			key.WithKeys(a.Keys...),
			key.WithHelp(a.ShortHelp, ""),
		))
	}
	return bindings
}

// FullHelp returns key bindings for the expanded help view in columns.
func (k KeyMap) FullHelp() [][]key.Binding {
	if len(k.Actions) == 0 {
		return nil
	}
	const colCount = 4
	cols := make([][]key.Binding, colCount)
	for i, a := range k.Actions {
		if a.FullHelp == "" {
			continue
		}
		binding := key.NewBinding(
			key.WithKeys(a.Keys...),
			key.WithHelp("", a.FullHelp),
		)
		cols[i%colCount] = append(cols[i%colCount], binding)
	}
	return cols
}

func defaultKeys(extra ...ActionKey) KeyMap {
	const spacebar = " "
	km := KeyMap{}
	km.Add(NewKey([]string{"up", "k"}, "up", "Move up", nil))
	km.Add(NewKey([]string{"down", "j"}, "down", "Move down", nil))
	km.Add(NewKey([]string{"b", "pgup"}, "", "Page up", nil))
	km.Add(NewKey([]string{"f", "pgdown", spacebar}, "", "Page down", nil))
	km.Add(NewKey([]string{"u", "ctrl+u"}, "", "½ page up", nil))
	km.Add(NewKey([]string{"d", "ctrl+d"}, "", "½ page down", nil))
	km.Add(NewKey([]string{"home", "g"}, "", "Go to start", nil))
	km.Add(NewKey([]string{"end", "G"}, "", "Go to end", nil))
	for _, k := range extra {
		km.Add(k)
	}
	km.Add(NewKey([]string{"?"}, "help", "Toggle help", nil))
	km.Add(NewKey([]string{"q", "esc"}, "quit", "Quit", nil))
	return km
}
