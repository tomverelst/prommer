package main

import "flag"
import "github.com/tomverelst/prommer/lib"

var (
	monitoringLabel = flag.String("monitoring-label", "prometheus-target", "Containers with this label will be added as target")
	targetFile      = flag.String("target-file", "/etc/prometheus/target-groups.json", "Target file to store the Prometheus target groups")
)

func main() {
	flag.Parse()

	options := &prommer.PrommerOptions{
		MonitoringLabel: *monitoringLabel,
		TargetFilePath:  *targetFile,
	}

	p := &prommer.Prommer{
		Options: options,
	}

	p.Start()
}
