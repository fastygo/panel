package main

import (
	"context"
	"fmt"

	"github.com/fastygo/codex/schema"
	"github.com/fastygo/panel"
)

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

func main() { run() }

func run() {
	cms, err := panel.NewPanel[principal, capability](panel.PanelOptions[capability]{
		ID:       "cms",
		Title:    "GoCMS",
		BasePath: "/go-admin",
	})
	if err != nil {
		panic(err)
	}
	if err := cms.AddResources(postsResource(), pagesResource()); err != nil {
		panic(err)
	}
	validator := panel.RecordValidator{Record: pagesRecord()}
	issues := validator.Validate(context.Background(), panel.Command{Operation: panel.OperationCreate, Data: panel.Record{"title": "About"}})
	fmt.Println("pages columns:", len(pagesResource().Table.Columns), "issues:", len(issues))
}

// pagesRecord is a Codex record: the single source of truth for fields.
func pagesRecord() schema.RecordType {
	return schema.RecordType{
		ID: "pages", Label: "Pages", Scope: schema.ScopeWorkspace,
		Fields: []schema.Field{
			{ID: "title", Label: "Title", Type: schema.FieldText, Required: true, Localized: true, Searchable: true, Indexed: true},
			{ID: "status", Label: "Status", Type: schema.FieldSelect, Indexed: true, Options: []schema.Option{
				{Value: "draft", Label: "Draft"}, {Value: "published", Label: "Published"},
			}},
			{ID: "body", Label: "Body", Type: schema.FieldRichText},
		},
	}
}

// pagesResource is projected from Codex: no hand-written columns or fields.
func pagesResource() panel.Resource[capability] {
	return panel.ResourceFromRecord(pagesRecord(), panel.ResourceOptions[capability]{
		BasePath:     "/go-admin/pages",
		Capabilities: map[panel.ResourceOperation]capability{panel.OperationList: contentRead, panel.OperationEdit: contentEdit},
	})
}
