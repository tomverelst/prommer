package prommer

// Service to be monitored
type Service struct {
	Name      string
	Instances []*Instance
}

// Instances of a Service
type Instance struct {
	HostIP   string
	HostPort string
}

type Monitor interface {
	Monitor(services []*Service)
}
