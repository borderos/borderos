package network

import (
	"log"
	"net"

	"github.com/borderos/borderos/config"

	"github.com/vishvananda/netlink"
)

// TODO: Configure sysctls for forwarding

func Configure(ifaces map[string]config.Interface, routing config.Routing) error {
	SetupInterfaces(ifaces)
	return SetDefaultRoute(routing.Default)
}

func SetupInterfaces(interfaces map[string]config.Interface) {
	for name, config := range interfaces {
		SetupInterface(name, config)
	}
}

func SetupInterface(name string, config config.Interface) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		log.Printf("couldn't configure interface %s: %v\n", name, err)
		return
	}
	for _, ip := range config.Addresses {
		addr, err := netlink.ParseAddr(ip)
		if err != nil {
			log.Printf("%s is not a valid addr: %v\n", ip, err)
			continue
		}
		if err := netlink.AddrAdd(link, addr); err != nil {
			log.Printf("couldn't add %s to link %s: %v\n", addr, name, err)
			continue
		}
	}
	if err := netlink.LinkSetUp(link); err != nil {
		log.Printf("couldn't set link %s up: %v\n", name, err)
	}
}

func SetDefaultRoute(addr string) error {
	route := &netlink.Route{
		Gw: net.ParseIP(addr),
	}
	return netlink.RouteAdd(route)
}
