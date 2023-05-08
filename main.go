package main

import (
	"flag"
	"fmt"
	"gitlab.tools.vipshop.com/terry.zeng/snmpmonitor/g"
	"gitlab.tools.vipshop.com/terry.zeng/snmpmonitor/http"
	"gitlab.tools.vipshop.com/terry.zeng/snmpmonitor/monitor"
	"os"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")

	flag.Parse()

	if *version {
		version := g.Version
		fmt.Printf("%s", version)
		os.Exit(0)
	}

	g.ParseConfig(*cfg)
	g.Logger = g.InitLog()

	go monitor.GlobalStart()

	go http.Start()

	select {}

}