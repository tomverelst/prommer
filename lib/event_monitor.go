package prommer

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	eventtypes "github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

// EventMonitor monitors the Docker event stream
type EventMonitor struct {
	docker  *client.Client
	filters filters.Args
	stream  chan eventtypes.Message
}

// Start monitoring
func (m *EventMonitor) Start() error {
	options := types.EventsOptions{
		Filters: m.filters,
	}

	ctx := context.Background()

	response, err := m.docker.Events(ctx, options)
	if err != nil {
		panic(err)
	}
	defer response.Close()
	return m.streamEvents(response)
}

func (m *EventMonitor) streamEvents(input io.Reader) error {
	return m.process(input, func(event eventtypes.Message, err error) error {
		if err != nil {
			return err
		}
		//	printOutput(event, output)
		m.stream <- event
		return nil
	})
}

type eventProcessor func(event eventtypes.Message, err error) error

func (m *EventMonitor) process(input io.Reader, processEvent eventProcessor) error {
	decoder := json.NewDecoder(input)
	for {
		var event eventtypes.Message
		err := decoder.Decode(&event)
		if err != nil && err == io.EOF {
			break
		}

		if processErr := processEvent(event, err); processErr != nil {
			return processErr
		}
	}
	return nil
}

// printOutput prints all types of event information.
// Each output includes the event type, actor id, name and action.
// Actor attributes are printed at the end if the actor has any.
func printOutput(event eventtypes.Message, output io.Writer) {

	// docker run: create -> start
	// docker stop: kill -> die -> stop
	// docker rm: (kill die stop) -> destroy
	if event.Type == "container" {
		fmt.Println("--------------")
		fmt.Println("Status: " + event.Status)
		fmt.Println("ID: " + event.ID)
		fmt.Println("From: " + event.From)
		fmt.Println("Type: " + event.Type)
		fmt.Println("Action: " + event.Action)
	}

}
