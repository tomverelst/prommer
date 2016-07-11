package prommer

type PrommerOptions struct {
	TargetFilePath  string
	MonitoringLabel string
}

// Service to be monitored
type Service struct {
	Name      string
	Instances []*Instance
}

// Instance of a Service
type Instance struct {
	HostIP   string
	HostPort string
}

// Monitor interface
type Monitor interface {
	Monitor(services []*Service)
}
