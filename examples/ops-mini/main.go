package main

import "github.com/fastygo/panel"

type capability string

const (
	servicesRead capability = "services.read"
	incidentsAck capability = "incidents.ack"
)

type principal struct {
	caps map[capability]struct{}
}

func (p principal) Has(cap capability) bool {
	_, ok := p.caps[cap]
	return ok
}

func main() {
	ops, err := panel.NewPanel[principal, capability](panel.PanelOptions[capability]{
		ID:       "ops",
		Title:    "Operations",
		BasePath: "/ops",
	})
	if err != nil {
		panic(err)
	}

	if err := ops.AddPages(panel.Page[capability]{
		ID:         "services",
		Kind:       panel.PageRuntime,
		Title:      "Services",
		Path:       "/ops/services",
		Capability: servicesRead,
		Navigation: panel.MenuItem[capability]{
			ID:         "services",
			Label:      "Services",
			Path:       "/ops/services",
			Capability: servicesRead,
			Order:      10,
		},
		Widgets: []panel.Widget[capability]{
			{ID: "health", Kind: panel.WidgetHealth, Title: "Health"},
		},
		Actions: []panel.Action[capability]{
			{ID: "ack", Label: "Acknowledge", Placement: panel.ActionHeader, Capability: incidentsAck},
		},
	}); err != nil {
		panic(err)
	}
}
