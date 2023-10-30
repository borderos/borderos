package config

import "net"

type ServiceDHCP struct{}

type ServiceDNS struct {
	Forwarders []net.IP `json:"forwarders"`
	Cache      bool     `json:"cache"`
}

type ServiceNTP struct {
	Enabled bool `json:"enabled"`
}

type ServiceSSH struct {
	Enabled bool     `json:"enabled"`
	Listen  []net.IP `json:"listen"`
	Port    uint16   `json:"port"`
}
