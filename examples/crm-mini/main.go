package main

import "github.com/fastygo/panel"

type capability string

const (
	leadsRead capability = "leads.read"
	leadsEdit capability = "leads.edit"
)

type principal struct {
	caps map[capability]struct{}
}

func (p principal) Has(cap capability) bool {
	_, ok := p.caps[cap]
	return ok
}

func main() {
	crm, err := panel.NewPanel[principal, capability](panel.PanelOptions[capability]{
		ID:       "crm",
		Title:    "CRM",
		BasePath: "/crm",
	})
	if err != nil {
		panic(err)
	}

	if err := crm.AddResources(panel.Resource[capability]{
		ID:       "leads",
		Label:    "Leads",
		Singular: "Lead",
		Plural:   "Leads",
		BasePath: "/crm/leads",
		Navigation: panel.MenuItem[capability]{
			ID:         "leads",
			Label:      "Leads",
			Path:       "/crm/leads",
			Capability: leadsRead,
			Order:      10,
		},
		Table: panel.TableSchema[capability]{
			Columns: []panel.Column{
				{ID: "name", Label: "Name", Type: panel.ColumnText, Searchable: true},
				{ID: "stage", Label: "Stage", Type: panel.ColumnBadge, Sortable: true},
			},
			RowActions: []panel.Action[capability]{
				{ID: "qualify", Label: "Qualify", Placement: panel.ActionRow, Capability: leadsEdit},
			},
		},
		Form: panel.FormSchema{
			Fields: []panel.Field{
				{ID: "name", Label: "Name", Type: panel.FieldText, Required: true},
				{ID: "email", Label: "Email", Type: panel.FieldText},
			},
		},
	}); err != nil {
		panic(err)
	}
}
