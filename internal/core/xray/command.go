package xray

import (
	"bgscan/internal/core/fileutil"
	"bgscan/internal/core/process"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

var xrayPaths = []string{
	"assets/xray",
	"xray",
}

// findXrayBinary attempts to locate the Xray executable.
func FindXrayBinary() (string, error) {
	return process.FindBinaryInPaths("xray", xrayPaths)
}

func XrayVersion() (string, error) {
	xrayBin, err := FindXrayBinary()
	if err != nil {
		return "", fmt.Errorf("xray binary not found: %w", err)
	}

	cmd := exec.Command(xrayBin, "-version")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("xray version check failed: %w\n%s", err, output)
	}

	version := strings.TrimSpace(string(output))
	return version, nil
}

// ValidateConfig verifies that a configuration file is valid
// by executing:
//
//	xray -c <config> --test
//
// If the configuration is invalid, the error returned will contain
// the full output produced by Xray to help diagnose the issue.
func ValidateConfig(configPath string) error {
	if !fileutil.CheckFileExists(configPath) {
		return fmt.Errorf("config file does not exist: %s", configPath)
	}

	xrayBin, err := FindXrayBinary()
	if err != nil {
		return err
	}

	cmd := exec.Command(xrayBin, "-c", configPath, "--test")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("xray config validation failed: %s", output)
	}

	return nil
}

// StartXray launches an Xray process using the provided configuration.
//
// The process is started asynchronously and returned as an XrayProcess
// instance so the caller can manage its lifecycle.
//
// The provided context controls the lifetime of the process. If the
// context is canceled, the Xray process will be terminated automatically.
func StartXray(ctx context.Context, configPath string) (*process.Process, error) {

	if !fileutil.CheckFileExists(configPath) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	xrayBin, err := FindXrayBinary()
	if err != nil {
		return nil, err
	}

	return process.Start(ctx, xrayBin, "-c", configPath)
}
