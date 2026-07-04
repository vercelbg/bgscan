package validate

import (
	"time"

	"bgscan/internal/core/config"
)

// ValidateWriter validates a WriterConfig strictly.
// Returns a map of field name → error for every invalid field.
// Used by the UI OnValidate hook and SaveWriterConfig.
func ValidateWriter(cfg *config.WriterConfig) map[string]error {
	errs := map[string]error{}

	if err := checkDuration("MergeFlushInterval", cfg.MergeFlushInterval.Duration(),
		100*time.Millisecond, 5*time.Minute); err != nil {
		errs["MergeFlushInterval"] = err
	}

	if err := checkInt("ChanSize", cfg.ChanSize, 1, 1_000_000); err != nil {
		errs["ChanSize"] = err
	}

	if err := checkInt("BatchSize", cfg.BatchSize, 1, 1_000_000); err != nil {
		errs["BatchSize"] = err
	}

	return errs
}

// NormalizeWriter auto-fixes invalid fields to their defaults.
// Returns a list of Warnings describing every correction made.
// Used only at TOML load time.
func NormalizeWriter(cfg *config.WriterConfig) []Warning {
	def := config.DefaultWriterConfig()
	var warns []Warning

	fixDurationMS("MergeFlushInterval", &cfg.MergeFlushInterval,
		100*time.Millisecond, 5*time.Minute, def.MergeFlushInterval, &warns)

	fixInt("ChanSize", &cfg.ChanSize, 1, 1_000_000, def.ChanSize, &warns)
	fixInt("BatchSize", &cfg.BatchSize, 1, 1_000_000, def.BatchSize, &warns)

	return warns
}
