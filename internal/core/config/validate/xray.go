package validate

import (
	"fmt"
	"time"

	"bgscan/internal/core/config"
)

var allowedPreScanTypes = []string{"tcp", "icmp", "none", "http"}

// ValidateXray validates an XrayConfig strictly.
// Returns a map of field name → error for every invalid field.
// Used by the UI OnValidate hook and SaveXrayConfig.
func ValidateXray(cfg *config.XrayConfig) map[string]error {
	errs := map[string]error{}

	if err := checkInt("Workers", cfg.Workers, 1, 1000); err != nil {
		errs["Workers"] = err
	}

	if !cfg.ConnectivityTestType.IsValid() {
		errs["ConnectivityTestType"] = errInvalidConnectivityTest()
	}

	if err := checkInt("DownloadSpeed", cfg.DownloadSpeed, 0, 10000); err != nil {
		errs["DownloadSpeed"] = err
	}

	if err := checkInt("UploadSpeed", cfg.UploadSpeed, 0, 10000); err != nil {
		errs["UploadSpeed"] = err
	}

	if err := checkDuration("Timeout", cfg.Timeout.Duration(),
		100*time.Millisecond, 60*time.Second); err != nil {
		errs["Timeout"] = err
	}

	if err := checkEnum("PreScanType", cfg.PreScanType, allowedPreScanTypes); err != nil {
		errs["PreScanType"] = err
	}

	if err := checkPrefix("PrefixOutput", cfg.PrefixOutput); err != nil {
		errs["PrefixOutput"] = err
	}

	return errs
}

// NormalizeXray auto-fixes invalid fields to their defaults.
// Returns a list of Warnings describing every correction made.
// Used only at TOML load time.
func NormalizeXray(cfg *config.XrayConfig) []Warning {
	def := config.DefaultXrayConfig()
	var warns []Warning

	fixInt("Workers", &cfg.Workers, 1, 1000, def.Workers, &warns)

	if !cfg.ConnectivityTestType.IsValid() {
		warns = append(warns, Warning{
			Field:  "ConnectivityTestType",
			OldVal: cfg.ConnectivityTestType,
			NewVal: def.ConnectivityTestType,
			Reason: "invalid → default",
		})
		cfg.ConnectivityTestType = def.ConnectivityTestType
	}

	fixInt("DownloadSpeed", &cfg.DownloadSpeed, 0, 10000, def.DownloadSpeed, &warns)
	fixInt("UploadSpeed", &cfg.UploadSpeed, 0, 10000, def.UploadSpeed, &warns)

	fixDurationMS("Timeout", &cfg.Timeout,
		100*time.Millisecond, 60*time.Second, def.Timeout, &warns)

	fixEnum("PreScanType", &cfg.PreScanType, allowedPreScanTypes, def.PreScanType, &warns)
	fixString("PrefixOutput", &cfg.PrefixOutput, def.PrefixOutput, &warns)

	return warns
}

func errInvalidConnectivityTest() error {
	return fmt.Errorf("must be one of: ConnectivityOnly, DownloadSpeedOnly, UploadSpeedOnly, Both")
}
