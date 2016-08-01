package prommer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/prometheus/common/log"
)

// NewPrometheusMonitor initializes a new PrometheusMonitor
func NewPrometheusMonitor(targetFilePath string) (*PrometheusMonitor, error) {

	if &targetFilePath == nil || targetFilePath == "" {
		return nil, errors.New("No path defined")
	}

	path, err := filepath.Abs(targetFilePath)
	if err != nil {
		return nil, errors.New("Could not resolve absolute path to target file path")
	}

	monitor := &PrometheusMonitor{
		targetFilePath: path,
	}

	return monitor, nil
}

// PrometheusMonitor
type PrometheusMonitor struct {
	targetFilePath string
}

type targetGroup struct {
	Targets []string          `json:"targets,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
}

// Monitor the given services.
// The target group configuration is updated and overrides the old configuration
func (m *PrometheusMonitor) Monitor(services []*Service) {
	targetGroups := m.createTargetGroups(services)

	content, err := json.Marshal(targetGroups)

	if err != nil {
		log.Errorln(err)
		return
	}

	tempFile := m.targetFilePath + ".tmp"

	fmt.Println(string(content[:]))

	if err := ioutil.WriteFile(tempFile, content, 0644); err != nil {
		log.Errorln(err)
	}

	err = os.Rename(tempFile, m.targetFilePath)
	if err != nil {
		log.Errorln(err)
	}

}

func (m *PrometheusMonitor) createTargetGroups(services []*Service) []*targetGroup {
	var groups []*targetGroup
	// Write files for current services.
	for _, service := range services {
		var targets []string
		for _, instance := range service.Instances {
			targets = append(targets, instance.HostIP+":"+instance.HostPort)
		}

		if len(targets) > 0 {
			labels := make(map[string]string)
			labels["job"] = service.Name

			groups = append(groups, &targetGroup{
				Targets: targets,
				Labels:  labels,
			})
		}
	}

	if groups == nil {
		groups = make([]*targetGroup, 0)
	}
	return groups
}
