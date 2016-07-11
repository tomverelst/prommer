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

// TargetGroup
type TargetGroup struct {
	Targets []string          `json:"targets,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`
}

// Monitor the given services.
// The target group configuration is updated and overrides the old configuration
func (m *PrometheusMonitor) Monitor(services []*Service) {

	var targetGroups []*TargetGroup
	// Write files for current services.
	for _, service := range services {
		var targets []string
		for _, instance := range service.Instances {
			targets = append(targets, instance.HostIP+":"+instance.HostPort)
		}

		targetGroups = append(targetGroups, &TargetGroup{
			Targets: targets,
		})
	}

	if targetGroups == nil {
		targetGroups = make([]*TargetGroup, 0)
	}

	content, err := json.Marshal(targetGroups)

	if err != nil {
		log.Errorln(err)
		return
	}

	f, err := m.getTargetFile()

	if err != nil {
		log.Errorln(err)
		return
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Errorln(err)
		}
	}()

	fmt.Println(string(content[:]))

	if _, writeError := f.Write(content); err != nil {
		log.Errorln(writeError)
	}

	if err := ioutil.WriteFile(f.Name(), content, 0644); err != nil {
		log.Errorln(err)
	}

}

func (m *PrometheusMonitor) getTargetFile() (*os.File, error) {
	var (
		file *os.File
	)

	var _, err = os.Stat(m.targetFilePath)

	if os.IsNotExist(err) {
		_, createError := os.Create(m.targetFilePath)
		if createError != nil {
			return nil, createError
		}
	}

	file, err = os.OpenFile(m.targetFilePath, os.O_WRONLY, 0644)

	if err != nil {
		return nil, err
	}

	return file, nil
}
