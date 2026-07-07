package probe

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"bgscan/internal/core/result"
	"bgscan/internal/logger"
)

// HTTPProbe performs HTTP/HTTPS validation against a target IP while preserving
// Host and SNI semantics.
type HTTPProbe struct {
	req    HTTPRequest
	filter statusFilter
	dialer *net.Dialer
	tls    *tls.Config
}

// NewHTTPProbe creates a new HTTPProbe with optional accepted status codes.
// If acceptedCodes is empty or covers all known codes, all responses are accepted.
func NewHTTPProbe(req HTTPRequest, acceptedCodes []int) Probe {
	return &HTTPProbe{
		req:    req,
		filter: newStatusFilter(acceptedCodes, totalHTTPStatusCodes),
		dialer: &net.Dialer{Timeout: req.Timeout},
		tls:    newTLSConfig(req),
	}
}

// Init implements Probe lifecycle (no-op).
func (p *HTTPProbe) Init(context.Context) error { return nil }

// Close implements Probe lifecycle (no-op).
func (p *HTTPProbe) Close() error { return nil }

// Run executes a single HTTP HEAD request against the target IP.
func (p *HTTPProbe) Run(ctx context.Context, ip string) (*result.IPScanResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, p.req.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	start := time.Now()

	resp, err := p.buildClient(ip).Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.CoreError("close response body: %v", err)
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

// buildClient creates an HTTP client bound to a specific target IP.
func (p *HTTPProbe) buildClient(ip string) *http.Client {
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			_, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, fmt.Errorf("parse addr: %w", err)
			}
			return p.dialer.DialContext(ctx, network, net.JoinHostPort(ip, port))
		},
		DisableKeepAlives:     true,
		TLSHandshakeTimeout:   p.req.Timeout,
		ResponseHeaderTimeout: p.req.Timeout,
		TLSClientConfig:       p.tls,
		ForceAttemptHTTP2:     p.req.UseTLS && p.req.Version != HTTPVersionH1,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   p.req.Timeout,
	}
}

func isHTTPS(proto string) bool {
	p := strings.ToLower(proto)
	p = strings.TrimSpace(p)
	return strings.HasPrefix(p, "https")
}
