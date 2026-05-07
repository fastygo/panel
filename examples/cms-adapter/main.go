package main

import "github.com/fastygo/panel"

type capability string

const (
	contentRead capability = "content.read"
	contentEdit capability = "content.edit"
)

type principal struct {
	caps map[capability]struct{}
}

func (p principal) Has(cap capability) bool {
	_, ok := p.caps[cap]
	return ok
}

func postsResource() panel.Resource[capability] {
	return panel.Resource[capability]{
		ID:       "posts",
		Label:    "Posts",
		Singular: "Post",
		Plural:   "Posts",
		BasePath: "/go-admin/posts",
		Navigation: panel.MenuItem[capability]{
			ID:         "posts",
			Label:      "Posts",
			Path:       "/go-admin/posts",
			Capability: contentRead,
			Order:      1,
		},
		Table: panel.TableSchema[capability]{
			Columns: []panel.Column{
				{ID: "title", Label: "Title", Type: panel.ColumnText, Searchable: true},
				{ID: "status", Label: "Status", Type: panel.ColumnBadge, Sortable: true},
			},
			RowActions: []panel.Action[capability]{
				{ID: "edit", Label: "Edit", Placement: panel.ActionRow, Capability: contentEdit},
			},
		},
		Form: panel.FormSchema{
			Fields: []panel.Field{
				{ID: "title", Label: "Title", Type: panel.FieldText, Required: true},
				{ID: "content", Label: "Content", Type: panel.FieldRichText},
			},
		},
	}
}

func main() {
	cms, err := panel.NewPanel[principal, capability](panel.PanelOptions[capability]{
		ID:       "cms",
		Title:    "GoCMS",
		BasePath: "/go-admin",
	})
	if err != nil {
		panic(err)
	}
	if err := cms.AddResources(postsResource()); err != nil {
		panic(err)
	}
}
