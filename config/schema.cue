import "net"

#System: {
	hostname!: string & !=""
	resolvers: [...net.IP & string]
}

#Interface: {
	addresses: [...net.IPCIDR & string]
	enabled: *false | bool
	lldp?: {
		enabled: *false | bool
	}
	vlan?: {
		[string]: #Interface
	}
}

#Routing: {
	default: net.IP & string
	static: [net.IPCIDR & string]: {
		interface?: string
		nexthop?:   string
	}
}

#NAT: {
	if type == "destination" {
		insideaddr!: string & net.IP
	}
	if type == "masquerade" {
		outinterface!: string
	}
	if type == "source" {
		outsideaddr!: string & net.IP & type == "source"
	}
	protocol: *"all" | "tcp" | "udp"
	type!:    "source" | "destination" | "masquerade"
}

#Firewall: {
	ipv4: [string]: [...#FirewallRule]
	ipv6: [string]: [...#FirewallRule]
	inet: [string]: [...#FirewallRule]
}

#FirewallRule: {
	source?: string & (net.IP | =~"^@")
	dest?:   string & (net.IP | =~"^@")
	action:  "accept" | "drop"
}

#Services: {
	dhcp: {}
	dns: {
		forwarders: [...net.IP & string]
		cache: bool | *false
	}
	ntp: {}
	radv: {}
	ssh: {
		enable: bool | *false
		listen: [...net.IP & string]
		port: uint16 | *22
	}
}

system: #System
interfaces: [string]: #Interface
routing:   #Routing
nat?:      #NAT
firewall?: #Firewall
services?: #Services
