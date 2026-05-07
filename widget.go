package panel

import (
	"fmt"
	"strings"
)

type WidgetID string

type WidgetKind string

const (
	WidgetStatCard     WidgetKind = "stat-card"
	WidgetChart        WidgetKind = "chart"
	WidgetTable        WidgetKind = "table"
	WidgetActivityFeed WidgetKind = "activity-feed"
	WidgetHealth       WidgetKind = "health"
	WidgetCustom       WidgetKind = "custom"
)

type Widget[C comparable] struct {
	ID          WidgetID
	Kind        WidgetKind
	Title       string
	Description string
	Icon        string
	Order       int
	Width       string
	Capability  C
	Table       TableSchema[C]
	Actions     []Action[C]
}

func (w Widget[C]) Validate() error {
	if strings.TrimSpace(string(w.ID)) == "" {
		return fmt.Errorf("widget id is required")
	}
	if strings.TrimSpace(w.Title) == "" {
		return fmt.Errorf("widget %q title is required", w.ID)
	}
	for _, action := range w.Actions {
		if err := action.Validate(); err != nil {
			return err
		}
	}
	if err := w.Table.Validate(); err != nil {
		return err
	}
	return nil
}
