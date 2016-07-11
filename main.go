package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tomverelst/prommer/lib"
)

var (
	monitoringLabel = flag.String("monitoring-label", "prometheus-target", "Containers with this label will be added as target")
	targetFile      = flag.String("target-file", "/etc/prometheus/target-groups.json", "Target file to store the Prometheus target groups")
)

func main() {
	flag.Parse()

	p := prommer.CreatePrommer(&prommer.PrommerOptions{
		MonitoringLabel: *monitoringLabel,
		TargetFilePath:  *targetFile,
	})

	go p.Start()

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	p.Stop()
}
