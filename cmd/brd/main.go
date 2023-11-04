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
)

var buildTime = "INFINITY"

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
	} else {
		// Possibly running as a gokrazy service, don't use formatting.
		fmt.Println("BorderRouter starting")
	}

	c, err := config.Load("/config/config.cue")
	if err != nil {
		log.Fatal(err)
	}

	if err := network.Configure(c.Interfaces, c.Routing); err != nil {
		log.Fatalf("couldn't configure network: %v", err)
	}
	fmt.Println("✅ Networking")

	if init {
		// Start the gokrazy supervisor just to get the web server and listener for IP changes.
		// TODO: Add things like SSH as a service here.
		if err := gokrazy.SuperviseServices(nil); err != nil {
			log.Printf("failed to start supervisor: %v\n", err)
		}

		pterm.Bold.Println(time.Now().Format(time.DateTime), "System ready")

		select {}
	}
}
