package table

import (
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

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

	colsWidth []int
	paddingY  int

	// Pending fields for Functional Options Pattern
	pendingCols []table.Column
	pendingRows []table.Row
}

// --- Functional Options Pattern ---

// Option configures the table Model.
type Option func(*Model)

// WithTitle sets the table title.
func WithTitle(title string) Option {
	return func(m *Model) { m.Title = title }
}

// WithColumns sets the initial columns.
func WithColumns(cols []table.Column) Option {
	return func(m *Model) { m.pendingCols = cols }
}

// WithRows sets the initial rows.
func WithRows(rows []table.Row) Option {
	return func(m *Model) { m.pendingRows = rows }
}

// WithPaddingY sets the vertical padding.
func WithPaddingY(padding int) Option {
	return func(m *Model) { m.paddingY = padding }
}

// WithKeyBindings appends custom key bindings to the default ones.
func WithKeyBindings(keys ...ActionKey) Option {
	return func(m *Model) { m.Keys = defaultKeys(keys...) }
}

// New creates a new table model using the Functional Options Pattern.
func New(lay *layout.Layout, opts ...Option) *Model {
	m := &Model{
		id:     ui.NewComponentID(),
		name:   "table",
		Layout: lay,
		Help:   help.New(),
		Keys:   defaultKeys(),
	}

	// Apply options
	for _, opt := range opts {
		opt(m)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Fallback to empty slices if not provided
	cols := m.pendingCols
	if cols == nil {
		cols = []table.Column{}
	}
	rows := m.pendingRows
	if rows == nil {
		rows = []table.Row{}
	}

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

// --- Component Interface ---

func (m *Model) Init() tea.Cmd      { return nil }
func (m *Model) ID() ui.ComponentID { return m.id }
func (m *Model) Name() string       { return m.name }
func (m *Model) OnClose() tea.Cmd   { return nil }
func (m *Model) Mode() env.Mode     { return env.NormalMode }

// --- Public Mutators (Thread-safe) ---

func (m *Model) SetPaddingY(padding int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.paddingY = padding
	m.updateTableSizeLocked()
}

func (m *Model) AppendRow(row table.Row) {
	m.mu.Lock()
	defer m.mu.Unlock()
	rows := append(slices.Clone(m.BubbleTable.Rows()), slices.Clone(row))
	m.BubbleTable.SetRows(rows)
	m.updateTableSizeLocked()
}

func (m *Model) SetRows(rows []table.Row) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.BubbleTable.SetRows(cloneRows(rows))
	m.updateTableSizeLocked()
}

// --- Helper Functions ---

func NewRowTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format("2006-01-02 15:04:05")
}

func cloneRows(rows []table.Row) []table.Row {
	out := make([]table.Row, len(rows))
	for i, row := range rows {
		out[i] = slices.Clone(row)
	}
	return out
}

func columnWidths(cols []table.Column) []int {
	widths := make([]int, len(cols))
	for i, c := range cols {
		widths[i] = c.Width
	}
	return widths
}

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

func (m *Model) updateTableSizeLocked() {
	if m.Layout == nil || m.Layout.Body.Height == 0 || m.Layout.Body.Width == 0 {
		return
	}
	if len(m.colsWidth) == 0 {
		return
	}

	cols := m.scaledColumns(m.BubbleTable.Columns())
	m.BubbleTable.SetColumns(cols)

	// Calculate available height
	helpHeight := lipgloss.Height(m.renderHelpView())
	titleHeight := lipgloss.Height(m.renderTitle())
	height := max(1, m.Layout.Body.Height-helpHeight-titleHeight-m.paddingY)

	m.BubbleTable.SetHeight(height)
	m.BubbleTable.SetWidth(m.tableWidthLocked())
}

func (m *Model) tableWidthLocked() int {
	if m.Layout == nil || m.Layout.Body.Width == 0 {
		return 80
	}
	return min(80, m.Layout.Body.Width-gutterWidth)
}

// --- Key Bindings ---

type ActionKey struct {
	Keys      []string
	ShortHelp string
	FullHelp  string
	Cmd       tea.Cmd
}

var arrowSymbols = map[string]string{
	"up": "↑", "down": "↓", "left": "←", "right": "→",
}

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

	return ActionKey{Keys: keys, ShortHelp: shortHelp, FullHelp: fullHelp, Cmd: cmd}
}

type KeyMap struct{ Actions []ActionKey }

func (k *KeyMap) Add(a ActionKey) { k.Actions = append(k.Actions, a) }

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

func (k KeyMap) FullHelp(width int) [][]key.Binding {
	if len(k.Actions) == 0 {
		return nil
	}
	colCount := 3
	if width > 90 {
		colCount = 4
	}

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
	km := KeyMap{}
	km.Add(NewKey([]string{"up", "k"}, "up", "Move up", nil))
	km.Add(NewKey([]string{"down", "j"}, "down", "Move down", nil))
	km.Add(NewKey([]string{"b", "pgup"}, "", "Page up", nil))
	km.Add(NewKey([]string{"f", "pgdown", " "}, "", "Page down", nil))
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

// SetKeys replaces the current key bindings with the provided action keys.
// It merges them with the default navigation keys.
func (m *Model) SetKeys(keys ...ActionKey) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Keys = defaultKeys(keys...)
}
