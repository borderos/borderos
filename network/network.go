package network

import (
	"log"
	"net"

	"github.com/borderos/borderos/config"
	"github.com/borderos/borderos/sysctl"

	"github.com/vishvananda/netlink"
)

func Configure(ifaces map[string]config.Interface, routing config.Routing) error {
	for _, ctl := range []string{"net.ipv4.ip_forward", "net.ipv6.conf.all.forwarding"} {
		if err := sysctl.Write(ctl, "1"); err != nil {
			log.Printf("couldn't configure sysctl %s: %v\n", ctl, err)
		}
	}
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
