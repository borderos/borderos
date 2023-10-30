package config

import (
	_ "embed"
	"fmt"
	"io"
	"os"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
)

//go:embed "schema.cue"
var schema string

type Config struct {
	System     System               `json:"system"`
	Interfaces map[string]Interface `json:"interfaces"`
	Routing    Routing              `json:"routing"`
	NAT        NAT                  `json:"nat"`
	Firewall   Firewall             `json:"firewall"`
	Services   Services             `json:"services"`
}

type System struct {
	Hostname  string   `json:"hostname"`
	Resolvers []string `json:"resolvers"`
}

type Interface struct {
	Addresses []string `json:"addresses"`
	Gateway   string   `json:"gateway"`
	LLDP      struct{} `json:"lldp"`
}

type Routing struct {
	Default string         `json:"default"`
}

type NAT struct {
	InsideAddr   string `json:"insideaddr"`
	OutInterface string `json:"outinterface"`
	OutsideAddr  string `json:"outsideaddr"`
	Protocol     string `json:"protocol"`
	Type         string `json:"type"`
}

// Each of IPv4, IPv6, Inet maps to a netfilter table.
// Each map key is a chain, comprised of rules.
type Firewall struct {
	IPv4 map[string][]FirewallRule
	IPv6 map[string][]FirewallRule
	Inet map[string][]FirewallRule
}

type FirewallRule struct {
	InIface  string
	OutIface string
	Source   string
	Dest     string
	Proto    string
	Type     string // Only with ICMP
	SrcPort  uint16
	DstPort  uint16
	State    string // Only with CT
	Action   string
	Comment  string
}

type Services struct {
	DHCP ServiceDHCP `json:"dhcp"`
	DNS  ServiceDNS  `json:"dns"`
	NTP  ServiceNTP  `json:"ntp"`
	SSH  ServiceSSH  `json:"ssh"`
}

// Load reads and parses the config at path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("couldn't open config: %w", err)
	}
	return Parse(f)
}

// Parse reads everything from r and parses the CUE syntax.
func Parse(r io.Reader) (*Config, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	ctx := cuecontext.New()
	s := ctx.CompileString(schema)
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("error in schema: %w", err)
	}

	v := ctx.CompileBytes(data, cue.Scope(s))
	if err := v.Err(); err != nil {
		return nil, fmt.Errorf("couldn't parse configuration: %w", err)
	}

	// Apply the schema to the configuration and validate against the schema
	v = s.Unify(v)
	if err := v.Validate(
		cue.Final(),
		cue.Concrete(true),
		cue.All(),
	); err != nil {
		return nil, fmt.Errorf("configuration isn't valid: %w", err)
	}

	var config Config
	err = v.Decode(&config)
	return &config, err
}
