# FastyGo Panel

`github.com/fastygo/panel` is a reusable Go control-plane package for admin panels and business application consoles.

It is inspired by Filament's panel/resource/page/widget model and by Frappe's metadata-driven application platform, while staying idiomatic for Go: the panel package describes control-plane contracts, and each application owns its own data plane.

## Status

Early `v0` design. APIs are expected to evolve while GoCMS and non-CMS examples validate the contracts.

## Install

```bash
go get github.com/fastygo/panel
```

## Control Plane And Data Plane

Panel owns:

- panels, resources, pages, widgets, actions, forms, tables, navigation, routes, assets, editor providers, policies, and optional data-source interfaces.

Applications own:

- CMS content services, CRM leads/deals, operations metrics/incidents, chat conversations, storage adapters, external APIs, and business workflows.

## Minimal Example

```go
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

	_ = admin.AddPages(panel.Page[capability]{
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
	})
}
```

## Package Map

- `Panel`: a mounted panel with resources, pages, widgets, assets, and navigation.
- `Registry`: route, menu, asset, and editor provider aggregation.
- `Resource`: CRUD-oriented control-plane descriptor.
- `Page`: arbitrary screen descriptor.
- `Widget`: dashboard or page component descriptor.
- `Action`: UI operation descriptor for header, row, bulk, modal, and inline actions.
- `FormSchema`, `TableSchema`, `DetailSchema`: portable schemas for app-specific delivery layers.
- `DataSource`, `RecordProvider`, `CommandHandler`: optional Frappe-like data-plane contracts.

## Documentation

- [Architecture](docs/architecture.md)
- [Concepts](docs/concepts.md)
- [Filament Comparison](docs/filament-comparison.md)
- [Frappe Comparison](docs/frappe-comparison.md)
- [GoCMS Integration](docs/cms-integration.md)
- [CRM Example](docs/crm-example.md)
- [Operations Example](docs/ops-example.md)
- [Chat Example](docs/chat-example.md)
- [API Stability](docs/api-stability.md)

## License

MIT. See [LICENSE](LICENSE).
