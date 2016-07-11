package prommer

import (
	"fmt"

	"github.com/docker/engine-api/client"
)

type Prommer struct {
	Options *PrommerOptions
}

// Start Prommer
func (p *Prommer) Start() {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)

	if err != nil {
		panic(err)
	}

	serviceProvider, err := CreateServiceProvider(cli, p.Options.MonitoringLabel)
	if err != nil {
		panic(err)
	}

	prometheus, err := NewPrometheusMonitor(p.Options.TargetFilePath)
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
