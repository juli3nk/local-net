package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
)

func (c *Config) Events(filterType events.Type) (<-chan events.Message, <-chan error) {
	var eventsOpts types.EventsOptions

	if filterType != "" {
		filterE := filters.NewArgs()
		filterE.Add("type", filterType)

		eventsOpts.Filters = filterE
	}

	return c.Client.Events(context.Background(), eventsOpts)
}
