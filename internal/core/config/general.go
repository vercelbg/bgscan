package config

// GeneralConfig defines global scanner behavior and execution settings.
type GeneralConfig struct {
	StatusInterval DurationMS `toml:"status_interval"`
	StopAfterFound int        `toml:"stop_after_found"`
	MaxIPsToTest   int        `toml:"max_ips_to_test"`
	PipelineMode   string     `toml:"pipeline_mode"`
	MaxIPsPerStage int        `toml:"max_ips_per_stage"`
	BatchSize      int        `toml:"batch_size"`
	Shuffled       bool       `toml:"shuffled"`
}
