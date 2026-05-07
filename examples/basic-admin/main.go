package main

import "github.com/fastygo/panel"

type capability string

type principal struct {
	caps map[capability]struct{}
}

func (p principal) Has(cap capability) bool {
	_, ok := p.caps[cap]
	return ok
}

func main() {
	admin, err := panel.NewPanel[principal, capability](panel.PanelOptions[capability]{
		ID:       "admin",
		Title:    "Admin",
		BasePath: "/admin",
	})
	if err != nil {
		panic(err)
	}

	if err := admin.AddPages(panel.Page[capability]{
		ID:    "dashboard",
		Kind:  panel.PageDashboard,
		Title: "Dashboard",
		Path:  "/admin",
		Navigation: panel.MenuItem[capability]{
			ID:    "dashboard",
			Label: "Dashboard",
			Path:  "/admin",
			Order: 0,
		},
	}); err != nil {
		panic(err)
	}
}
