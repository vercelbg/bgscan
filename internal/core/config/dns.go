package config

// DNSConfig represents the top‑level DNS configuration, combining resolver,
// DNSTT, and SlipStream settings.
type DNSConfig struct {
	Resolver   *ResolverConfig   `toml:"resolver"`
	DNSTT      *DNSTTConfig      `toml:"dnstt"`
	SlipStream *SlipStreamConfig `toml:"slip_stream"`
}

///////////////////////////////////////////////////////////////////////////////
// Resolver
///////////////////////////////////////////////////////////////////////////////

// ResolverConfig defines settings for traditional DNS resolvers.
type ResolverConfig struct {
	Workers         int        `toml:"workers"`
	Protocol        string     `toml:"protocol"`
	Domain          string     `toml:"domain"`
	Port            uint16     `toml:"port"`
	CheckTypes      []string   `toml:"check_types"`
	EDNSBufSize     uint16     `toml:"ends_buffer_size"`
	Timeout         DurationMS `toml:"timeout"`
	Tries           int        `toml:"tries"`
	RandomSubdomain bool       `toml:"random_subdomain"`
	AcceptedRCodes  []string   `toml:"accepted_rcodes"`
	CheckDPI        bool       `toml:"check_dpi"`
	DPITimeout      DurationMS `toml:"dpi_timeout"`
	DPITries        int        `toml:"dpi_tries"`
	PrefixOutput    string     `toml:"prefix_output"`
}

///////////////////////////////////////////////////////////////////////////////
// DNSTT
///////////////////////////////////////////////////////////////////////////////

// DNSTTConfig defines configuration for DNSTT (DNS Tunnel Transport) scanning.
type DNSTTConfig struct {
	Enabled      bool       `toml:"enabled"`
	Workers      int        `toml:"workers"`
	Domain       string     `toml:"domain"`
	PublicKey    string     `toml:"public_key"`
	Timeout      DurationMS `toml:"timeout"`
	PrefixOutput string     `toml:"prefix_output"`
}

///////////////////////////////////////////////////////////////////////////////
// SlipStream
///////////////////////////////////////////////////////////////////////////////

// SlipStreamConfig defines configuration for SlipStream-based DNS scanning.
type SlipStreamConfig struct {
	Enabled      bool       `toml:"enabled"`
	Workers      int        `toml:"workers"`
	Domain       string     `toml:"domain"`
	CertPath     string     `toml:"cert_path"`
	Timeout      DurationMS `toml:"timeout"`
	PrefixOutput string     `toml:"prefix_output"`
}
