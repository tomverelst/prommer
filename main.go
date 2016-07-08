package main

import (
	//"encoding/json"
	"flag"
	"fmt"

	"github.com/docker/engine-api/client"
	"github.com/tomverelst/prommer/lib"
)

var (
	monitoringLabel = flag.String("monitoring-label", "prometheus-target", "Containers with this label will be added as target")
	targetFile      = flag.String("target-file", "/etc/prometheus/target-groups.json", "Target file to store the Prometheus target groups")
)

func main() {
	flag.Parse()

	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)

	if err != nil {
		panic(err)
	}

	serviceProvider, err := prommer.CreateServiceProvider(cli, *monitoringLabel)
	if err != nil {
		panic(err)
	}

	prometheus, err := prommer.NewPrometheusMonitor(*targetFile)
	if err != nil {
		panic(err)
	}

	services, err := serviceProvider.GetServices()

	for _, s := range services {
		fmt.Println(s.Name)
		for _, instance := range s.Instances {
			fmt.Println(instance.HostIP)
			fmt.Println(instance.HostPort)
		}
	}

	if err != nil {
		panic(err)
	}

	prometheus.Monitor(services)

}
