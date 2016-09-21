package prommer

import (
	"strconv"

	dockertypes "github.com/docker/docker/api/types"
)

// FindPortOptions holds the optional parameters for FindPort
type FindPortOptions struct {
	Label *string
}

// FindPort finds a port for a container
type FindPort func(dockertypes.Container, *FindPortOptions) int

// FindPortFromContainer attempts to find the correct port from the container
func FindPortFromContainer(c dockertypes.Container, options *FindPortOptions) int {
	var (
		portFromLabels int
		//publicPorts    []dockertypes.Port
		amountOfPorts int
		port          *dockertypes.Port
	)

	portFromLabels = findPortFromLabels(c, options)

	if portFromLabels != 0 {
		return portFromLabels
	}

	//publicPorts = onlyPublicPorts(c.Ports)
	amountOfPorts = len(c.Ports)

	if amountOfPorts == 0 {
		return 0
	}

	if amountOfPorts == 1 {
		return c.Ports[0].PrivatePort
	}

	// If there are multiple ports
	// prefer the one that forwards to port 80
	port = findPortEighty(c.Ports)

	if port == nil {
		port = &c.Ports[0]
	}

	if port == nil {
		return 0
	}

	return port.PrivatePort
}

func findPortFromLabels(c dockertypes.Container, options *FindPortOptions) int {
	if options != nil && options.Label != nil {
		portString := c.Labels[*options.Label]
		if &portString != nil {
			port, err := strconv.Atoi(portString)
			if err != nil {
				return 0
			}
			return port
		}
	}
	return 0
}

func onlyPublicPorts(ports []dockertypes.Port) []dockertypes.Port {
	var publicPorts []dockertypes.Port
	for _, port := range ports {
		if port.PublicPort != 0 {
			publicPorts = append(publicPorts, port)
		}
	}
	return publicPorts
}

func findPortEighty(ports []dockertypes.Port) *dockertypes.Port {
	for _, port := range ports {
		if port.PrivatePort == 80 {
			return &port
		}
	}
	return nil
}
