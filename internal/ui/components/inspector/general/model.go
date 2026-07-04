package general

import (
	"strconv"
	"time"

	"bgscan/internal/core/config"
	"bgscan/internal/core/config/validate"
	"bgscan/internal/ui/components/basic/input"
	"bgscan/internal/ui/components/basic/input/selectinput"
	"bgscan/internal/ui/components/basic/input/textinput"
	"bgscan/internal/ui/components/basic/input/toggleinput"
	"bgscan/internal/ui/components/basic/inspector"
	"bgscan/internal/ui/components/basic/notice"
	"bgscan/internal/ui/shared/env"
	"bgscan/internal/ui/shared/layout"
	"bgscan/internal/ui/shared/ui"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
)

const (
	groupGeneral = "General"
	groupWriter  = "Writer"
)

type pipelineMode = string

const (
	pipelineSequential pipelineMode = "sequential"
	pipelineStreaming  pipelineMode = "streaming"
	pipelineBatch      pipelineMode = "batch"
)

const (
	descStatusInterval     = "The interval in milliseconds for pushing status updates to the UI."
	descStopAfterFound     = "The maximum number of successful results before halting the scan. Set to 0 to scan all targets."
	descMaxIPsToTest       = "The maximum number of IPs to read from the input source. Set to 0 to read all available IPs."
	descShuffled           = "Randomizes the target IP order before scanning to prevent subnet slamming and reduce firewall alerts."
	descPipelineMode       = "The execution mode for multi-stage scanning: 'sequential' (disk-based), 'streaming' (channel-based), or 'batch' (hybrid)."
	descMaxIPsPerStage     = "The maximum number of IPs a pipeline stage can hold in memory. Exceeding this limit blocks the previous stage."
	descBatchSize          = "The number of IPs processed per batch when using the 'batch' pipeline mode."
	descMergeFlushInterval = "The interval in milliseconds for merging delta results into the main result file."
	descChanSize           = "The capacity of the internal channel used by scanner workers to send IP scan results to the writer goroutine."
	descWriterBatchSize    = "The initial capacity of the in-memory batch used to accumulate IP scan results before flushing to disk."
)

type Model struct {
	layout    *layout.Layout
	name      string
	id        ui.ComponentID
	inspector ui.Component
}

func (m *Model) ID() ui.ComponentID { return m.id }
func (m *Model) Init() tea.Cmd      { return nil }
func (m *Model) Mode() env.Mode     { return env.NormalMode }
func (m *Model) Name() string       { return m.name }
func (m *Model) OnClose() tea.Cmd   { return nil }

func saveGeneral(l *layout.Layout, cfg *config.GeneralConfig) tea.Cmd {
	if err := config.SaveGeneralConfig(cfg); err != nil {
		return notice.NewNoticeCmd(l, "Failed to save General settings", err.Error(), notice.NOTICE_ERROR)
	}
	return nil
}

func saveWriter(l *layout.Layout, cfg *config.WriterConfig) tea.Cmd {
	if err := config.SaveWriterConfig(cfg); err != nil {
		return notice.NewNoticeCmd(l, "Failed to save Writer settings", err.Error(), notice.NOTICE_ERROR)
	}
	return nil
}

func intInput(l *layout.Layout, title string, value int, val func(string) error, set func(int), save func() tea.Cmd) input.Input[string] {
	return textinput.New(
		l, title,
		textinput.WithValue(strconv.Itoa(value)),
		textinput.WithValidation(val),
		textinput.WithFocus(),
		textinput.WithOnSubmit(func(v string) tea.Cmd {
			n, err := strconv.Atoi(v)
			if err != nil {
				return notice.NewNoticeCmd(l, "Invalid "+title, err.Error(), notice.NOTICE_ERROR)
			}
			set(n)
			return save()
		}),
	)
}

func durationMSInput(l *layout.Layout, title string, value time.Duration, val func(string) error, set func(time.Duration), save func() tea.Cmd) input.Input[string] {
	return textinput.New(
		l, title,
		textinput.WithValue(strconv.FormatInt(value.Milliseconds(), 10)),
		textinput.WithValidation(val),
		textinput.WithFocus(),
		textinput.WithOnSubmit(func(v string) tea.Cmd {
			n, err := strconv.Atoi(v)
			if err != nil {
				return notice.NewNoticeCmd(l, "Invalid "+title, err.Error(), notice.NOTICE_ERROR)
			}
			set(time.Duration(n) * time.Millisecond)
			return save()
		}),
	)
}

func New(l *layout.Layout, name string) *Model {
	cfgG := config.GetGeneral()
	cfgW := config.GetWriter()

	saveGen := func() tea.Cmd { return saveGeneral(l, cfgG) }
	saveWri := func() tea.Cmd { return saveWriter(l, cfgW) }

	// ── General ──────────────────────────────────────────────────────────────

	statusInterval := durationMSInput(l, "Enter Status Interval", cfgG.StatusInterval.Duration(),
		func(v string) error {
			n, err := strconv.Atoi(v)
			if err != nil {
				return err
			}
			tmp := *cfgG
			tmp.StatusInterval = config.NewDurationMS(time.Duration(n) * time.Millisecond)
			return fieldErr(validate.ValidateGeneral(&tmp), "StatusInterval")
		},
		func(d time.Duration) { cfgG.StatusInterval = config.NewDurationMS(d) }, saveGen)

	stopAfterFound := intInput(l, "Enter Stop After Found (0 = unlimited)", cfgG.StopAfterFound,
		func(v string) error {
			n, err := strconv.Atoi(v)
			if err != nil {
				return err
			}
			tmp := *cfgG
			tmp.StopAfterFound = n
			return fieldErr(validate.ValidateGeneral(&tmp), "StopAfterFound")
		},
		func(n int) { cfgG.StopAfterFound = n }, saveGen)

	maxIPsToTest := intInput(l, "Enter Max IPs To Test (0 = unlimited)", cfgG.MaxIPsToTest,
		func(v string) error {
			n, err := strconv.Atoi(v)
			if err != nil {
				return err
			}
			tmp := *cfgG
			tmp.MaxIPsToTest = n
			return fieldErr(validate.ValidateGeneral(&tmp), "MaxIPsToTest")
		},
		func(n int) { cfgG.MaxIPsToTest = n }, saveGen)

	shuffled := toggleinput.New(
		l, "Shuffle",
		toggleinput.WithValue(cfgG.Shuffled),
		toggleinput.WithFocus(),
		toggleinput.WithLabels("Enabled", "Disabled"),
		toggleinput.WithOnSubmit(func(v bool) tea.Cmd {
			cfgG.Shuffled = v
			return saveGen()
		}),
	)

	pipelineMode := selectinput.New(
		l, "Select Pipeline Mode",
		selectinput.WithValue(cfgG.PipelineMode),
		selectinput.WithFocus[string](),
		selectinput.WithOptions(
			huh.NewOption("Sequential", pipelineSequential),
			huh.NewOption("Streaming", pipelineStreaming),
			huh.NewOption("Batch", pipelineBatch),
		),
		selectinput.WithOnSubmit(func(v string) tea.Cmd {
			cfgG.PipelineMode = v
			return saveGen()
		}),
	)

	maxIPsPerStage := intInput(l, "Enter Max IPs Per Stage", cfgG.MaxIPsPerStage,
		func(v string) error {
			n, err := strconv.Atoi(v)
			if err != nil {
				return err
			}
			tmp := *cfgG
			tmp.MaxIPsPerStage = n
			return fieldErr(validate.ValidateGeneral(&tmp), "MaxIPsPerStage")
		},
		func(n int) { cfgG.MaxIPsPerStage = n }, saveGen)

	batchSize := intInput(l, "Batch Size", cfgG.BatchSize,
		func(v string) error {
			n, err := strconv.Atoi(v)
			if err != nil {
				return err
			}
			tmp := *cfgG
			tmp.BatchSize = n
			return fieldErr(validate.ValidateGeneral(&tmp), "BatchSize")
		},
		func(n int) { cfgG.BatchSize = n }, saveGen)

	// ── Writer ────────────────────────────────────────────────────────────────

	mergeFlushInterval := durationMSInput(l, "Enter Merge Flush Interval", cfgW.MergeFlushInterval.Duration(),
		func(v string) error {
			n, err := strconv.Atoi(v)
			if err != nil {
				return err
			}
			tmp := *cfgW
			tmp.MergeFlushInterval = config.NewDurationMS(time.Duration(n) * time.Millisecond)
			return fieldErr(validate.ValidateWriter(&tmp), "MergeFlushInterval")
		},
		func(d time.Duration) { cfgW.MergeFlushInterval = config.NewDurationMS(d) }, saveWri)

	chanSize := intInput(l, "Enter Channel Size", cfgW.ChanSize,
		func(v string) error {
			n, err := strconv.Atoi(v)
			if err != nil {
				return err
			}
			tmp := *cfgW
			tmp.ChanSize = n
			return fieldErr(validate.ValidateWriter(&tmp), "ChanSize")
		},
		func(n int) { cfgW.ChanSize = n }, saveWri)

	writerBatchSize := intInput(l, "Enter Batch Size", cfgW.BatchSize,
		func(v string) error {
			n, err := strconv.Atoi(v)
			if err != nil {
				return err
			}
			tmp := *cfgW
			tmp.BatchSize = n
			return fieldErr(validate.ValidateWriter(&tmp), "BatchSize")
		},
		func(n int) { cfgW.BatchSize = n }, saveWri)

	fields := []inspector.Field{
		{Name: "Status Interval", Description: descStatusInterval, Group: groupGeneral, Input: inspector.Adapt(statusInterval), Visible: alwaysVisible, Format: inspector.FormatDurationMS},
		{Name: "Stop After Found", Description: descStopAfterFound, Group: groupGeneral, Input: inspector.Adapt(stopAfterFound), Visible: alwaysVisible, Format: inspector.FormatIntOrUnlimited},
		{Name: "Max IPs To Test", Description: descMaxIPsToTest, Group: groupGeneral, Input: inspector.Adapt(maxIPsToTest), Visible: alwaysVisible, Format: inspector.FormatIntOrUnlimited},
		{Name: "Shuffled", Description: descShuffled, Group: groupGeneral, Input: inspector.Adapt(shuffled), Visible: alwaysVisible, Format: inspector.FormatBool},
		{Name: "Pipeline Mode", Description: descPipelineMode, Group: groupGeneral, Input: inspector.Adapt(pipelineMode), Visible: alwaysVisible},
		{Name: "Max IPs Per Stage", Description: descMaxIPsPerStage, Group: groupGeneral, Input: inspector.Adapt(maxIPsPerStage), Visible: visibleWhenMode(cfgG, pipelineStreaming), Format: inspector.FormatInt},
		{Name: "Batch Size", Description: descBatchSize, Group: groupGeneral, Input: inspector.Adapt(batchSize), Visible: visibleWhenMode(cfgG, pipelineBatch), Format: inspector.FormatInt},

		{Name: "Merge Flush Interval", Description: descMergeFlushInterval, Group: groupWriter, Input: inspector.Adapt(mergeFlushInterval), Visible: alwaysVisible, Format: inspector.FormatDurationMS},
		{Name: "Channel Size", Description: descChanSize, Group: groupWriter, Input: inspector.Adapt(chanSize), Visible: alwaysVisible, Format: inspector.FormatInt},
		{Name: "Batch Size", Description: descWriterBatchSize, Group: groupWriter, Input: inspector.Adapt(writerBatchSize), Visible: alwaysVisible, Format: inspector.FormatInt},
	}

	return &Model{
		layout:    l,
		name:      name,
		id:        ui.NewComponentID(),
		inspector: inspector.New(l, "general settings", fields),
	}
}

func alwaysVisible() bool { return true }

func visibleWhenMode(cfg *config.GeneralConfig, want pipelineMode) func() bool {
	return func() bool { return cfg.PipelineMode == want }
}

func fieldErr(errs map[string]error, field string) error {
	return errs[field]
}
