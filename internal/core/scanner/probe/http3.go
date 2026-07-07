package probe

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/quic-go/quic-go/http3"

	"bgscan/internal/core/config"
	"bgscan/internal/core/result"
	"bgscan/internal/core/scanner/netutil"
	"bgscan/internal/logger"
)

// HTTP3Probe performs HTTP/3 (QUIC-based) probing against an IP address.
//
// It sends a single HTTP HEAD request over QUIC and validates the response
// status code against an optional allow-list.
//
// The probe forces QUIC connections to the target IP while preserving the
// original hostname in the Host header and TLS SNI for virtual-host routing.
type HTTP3Probe struct {
	req       HTTPRequest
	filter    statusFilter
	transport *http3.Transport
}

// NewHTTP3Probe creates a new HTTP/3 probe instance.
//
// acceptedCodes defines an allow-list of HTTP status codes considered valid.
// If empty, all status codes are accepted.
//
// The probe uses QUIC (HTTP/3 over TLS 1.3) and configures the TLS settings
// based on the request parameters.
func NewHTTP3Probe(req HTTPRequest, acceptedCodes []int) (Probe, error) {
	// HTTP/3 always runs over TLS 1.3 (RFC 9001); Min/MaxTLSVersion are
	// not configurable at this layer — enforce that here.
	tlsCfg := newTLSConfig(req)

	return &HTTP3Probe{
		req:    req,
		filter: newStatusFilter(acceptedCodes, totalHTTPStatusCodes),
		transport: &http3.Transport{
			TLSClientConfig: tlsCfg,
		},
	}, nil
}

// Init prepares the probe for execution. HTTP/3 probe requires no initialization.
func (p *HTTP3Probe) Init(_ context.Context) error { return nil }

// Close releases underlying QUIC transport resources.
func (p *HTTP3Probe) Close() error {
	if err := p.transport.Close(); err != nil {
		return fmt.Errorf("close http3 transport: %w", err)
	}
	return nil
}

// Run executes a single HTTP/3 HEAD probe against the given IP address.
//
// Behavior:
//   - Sends one HTTP HEAD request over QUIC.
//   - Overrides URL host to the target IP (direct dial).
//   - Preserves original Host header + SNI for virtual hosting.
//   - Validates HTTP status code against the configured allow-list.
//
// Returns an IPScanResult containing latency on success.
func (p *HTTP3Probe) Run(ctx context.Context, ip string) (*result.IPScanResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, p.req.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	// Force QUIC to dial raw IP while preserving Host + SNI.
	req.URL.Host = net.JoinHostPort(ip, req.URL.Port())
	req.Host = p.req.Host

	client := &http.Client{
		Transport: p.transport,
		Timeout:   p.req.Timeout,
	}

	start := time.Now()

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			logger.CoreError("close response body: %v", cerr)
		}
	}()

	if !p.filter.isAccepted(resp.StatusCode) {
		return nil, fmt.Errorf("status %d not accepted", resp.StatusCode)
	}

	return &result.IPScanResult{
		IP:      ip,
		Latency: time.Since(start),
	}, nil
}

// NewHTTP3RequestFromConfig builds an HTTPRequest suitable for HTTP/3 probing.
//
// Notes:
//   - HTTP/3 always uses TLS 1.3 (RFC 9001).
//   - Any TLS version fields in cfg are intentionally ignored.
//   - The scheme is always forced to https.
//   - Host and SNI are normalized for QUIC compatibility.
func NewHTTP3RequestFromConfig(cfg config.HTTPConfig) (*HTTPRequest, error) {
	const scheme = "https://"

	host, err := netutil.ExtractTLSServerName(cfg.Host)
	if err != nil {
		return nil, fmt.Errorf("extract host: %w", err)
	}

	urlHost, err := netutil.NormalizeHostWithSuffix(cfg.Host)
	if err != nil {
		return nil, fmt.Errorf("normalize host: %w", err)
	}

	port, err := defaultPort(cfg.Port, true)
	if err != nil {
		return nil, err
	}

	sni, err := resolveSNI(cfg.ServerName, true)
	if err != nil {
		return nil, fmt.Errorf("resolve SNI: %w", err)
	}

	return &HTTPRequest{
		URL:           fmt.Sprintf("%s%s:%d", scheme, urlHost, port),
		Host:          host,
		SNI:           sni,
		UseTLS:        true,
		SkipTLSVerify: !cfg.TLSValidation,
		Timeout:       cfg.Timeout.Duration(),
	}, nil
}
