package storage

import (
	"context"
	"strings"
	"time"

	grafana "github.com/kihamo/go-grafana-api"
	"github.com/mrsmtvd/shadow/components/annotations"
)

type Grafana struct {
	client     *grafana.Client
	dashboards []int64
}

func NewGrafana(address, apiKey, username, password string, dashboards []int64, logger grafana.Logger) *Grafana {
	client := grafana.New(address)

	if apiKey != "" {
		client = client.WithApiKey(apiKey)
	} else {
		client = client.WithBasicAuth(username, password)
	}

	if logger != nil {
		client = client.WithLogger(logger)
	}

	return &Grafana{
		client:     client,
		dashboards: dashboards,
	}
}

func (s *Grafana) Create(annotation annotations.Annotation) (err error) {
	text := []string{
		annotation.Title(),
		annotation.Text(),
	}

	input := &grafana.CreateAnnotationInput{
		Time: grafana.Int64(annotation.Time().UnixNano() / int64(time.Millisecond)),
		Text: grafana.String(strings.Join(text, "\n")),
		Tags: grafana.StringSlice(annotation.Tags()),
	}

	if annotation.TimeEnd() != nil {
		input.TimeEnd = grafana.Int64(annotation.TimeEnd().UnixNano() / int64(time.Millisecond))
		input.IsRegion = grafana.Bool(true)
	}

	for _, dashboardID := range s.dashboards {
		input.DashboardId = grafana.Int64(dashboardID)

		if _, err = s.client.CreateAnnotation(context.Background(), input); err != nil {
			break
		}
	}

	return err
}
