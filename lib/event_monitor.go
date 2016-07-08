package prommer

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/docker/docker/pkg/jsonlog"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	eventtypes "github.com/docker/engine-api/types/events"
	"github.com/docker/engine-api/types/filters"
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

	if event.TimeNano != 0 {
		fmt.Fprintf(output, "%s ", time.Unix(0, event.TimeNano).Format(jsonlog.RFC3339NanoFixed))
	} else if event.Time != 0 {
		fmt.Fprintf(output, "%s ", time.Unix(event.Time, 0).Format(jsonlog.RFC3339NanoFixed))
	}

	fmt.Fprintf(output, "%s %s %s", event.Type, event.Action, event.Actor.ID)

	if len(event.Actor.Attributes) > 0 {
		var attrs []string
		var keys []string
		for k := range event.Actor.Attributes {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := event.Actor.Attributes[k]
			attrs = append(attrs, fmt.Sprintf("%s=%s", k, v))
		}
		fmt.Fprintf(output, " (%s)", strings.Join(attrs, ", "))
	}
	fmt.Fprint(output, "\n")

}
