// BRD is the BorderRouter daemon.
package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/borderos/borderos/config"
	"github.com/borderos/borderos/network"

	"github.com/gokrazy/gokrazy"
	"github.com/pterm/pterm"
	"golang.org/x/sys/unix"
)

var buildTime = "INFINITY"

func mountfs() error {
	if err := os.Mkdir("/perm/etc", 0750); !os.IsExist(err) {
		return err
	}
	if err := os.Mkdir("/perm/.w", 0750); !os.IsExist(err) {
		return err
	}
	var flags uintptr = unix.MOUNT_ATTR_NOATIME | unix.MOUNT_ATTR_NODEV | unix.MOUNT_ATTR_NOEXEC
	if err := unix.Mount("overlay", "/etc", "overlay", flags, "lowerdir=/etc,upperdir=/perm/etc,workdir=/perm/.w"); err != nil {
		log.Printf("couldn't mount overlay for /etc: %v\n", err)
	}
	return nil
}

func main() {
	init := os.Getpid() == 1

	if init {
		pterm.Bold.Println(time.Now().Format(time.DateTime), "System starting...")
		if bi, ok := debug.ReadBuildInfo(); ok {
			fmt.Printf("Build: %s @ %s [%s]\n\n", bi.Main.Version, buildTime, bi.GoVersion)
		}
		pterm.DefaultBox.WithLeftPadding(2).WithRightPadding(2).Println("BorderOS")
		fmt.Println()

		if err := gokrazy.Boot(buildTime); err != nil {
			log.Fatal(err)
		}
		// Re-mount /etc with an overlay.
		if err := mountfs(); err != nil {
			log.Printf("couldn't re-mount /etc rw: %v\n", err)
		}
	} else {
		// Possibly running as a gokrazy service, don't use formatting.
		fmt.Println("BorderRouter starting")
	}

	c, err := config.Load("/config/config.cue")
	if err != nil {
		log.Fatal(err)
	}

	{
		if err := network.Configure(c.Interfaces, c.Routing); err != nil {
			log.Fatalf("couldn't configure network: %v", err)
		}
		if err := network.SetupResolver(c.System.Resolvers); err != nil {
			log.Printf("couldn't configure DNS resolvers: %v\n", err)
		}
	}
	fmt.Println("✅ Networking")

	if init {
		// N.B. package gokrazy stores hostname as a global in gokrazy.Boot().
		// Setting it here won't update that, so the gokrazy status page still shows the old hostname.
		if c.System.Hostname != "" {
			log.Println("Setting hostname to ", c.System.Hostname)
			if err := os.WriteFile("/etc/hostname", []byte(c.System.Hostname), 0644); err != nil {
				log.Printf("couldn't update /etc/hostname: %v\n", err)
			}
			if err := unix.Sethostname([]byte(c.System.Hostname)); err != nil {
				log.Printf("couldn't set hostname: %v\n", err)
			}
		}

		// Start the gokrazy supervisor just to get the web server and listener for IP changes.
		// TODO: Add things like SSH as a service here.
		if err := gokrazy.SuperviseServices(nil); err != nil {
			log.Printf("failed to start supervisor: %v\n", err)
		}

		pterm.Bold.Println(time.Now().Format(time.DateTime), "System ready")

		select {}
	}
}
