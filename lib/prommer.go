package prommer

import (
	"time"

	"github.com/docker/docker/client"
)

// Prommer
type Prommer struct {
	Options         *PrommerOptions
	serviceProvider *ServiceProvider
	prometheus      *PrometheusMonitor
}

// CreatePrommer creates a new instance of Prommer
func CreatePrommer(options *PrommerOptions) *Prommer {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	docker, err := client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)

	if err != nil {
		panic(err)
	}

	serviceProvider, err := CreateServiceProvider(docker, options.MonitoringLabel)
	if err != nil {
		panic(err)
	}

	prometheus, err := NewPrometheusMonitor(options.TargetFilePath)
	if err != nil {
		panic(err)
	}

	p := &Prommer{
		Options:         options,
		serviceProvider: serviceProvider,
		prometheus:      prometheus,
	}

	return p
}

// Start Prommer
func (p *Prommer) Start() {
	for {
		services, err := p.serviceProvider.GetServices()

		if err != nil {
			panic(err)
		}

		p.prometheus.Monitor(services)

		time.Sleep(time.Duration(10) * time.Second)
	}
}

func (p *Prommer) Stop() {

}
