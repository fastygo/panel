package panel

import (
	"context"
	"testing"

	"github.com/fastygo/codex/schema"
	"github.com/fastygo/formset"
)

func articleRecord() schema.RecordType {
	return schema.RecordType{
		ID: "article", Label: "Articles", Description: "Blog", Scope: schema.ScopeWorkspace,
		Fields: []schema.Field{
			{ID: "title", Label: "Title", Type: schema.FieldText, Required: true, Localized: true, Searchable: true, Indexed: true},
			{ID: "body", Label: "Body", Type: schema.FieldRichText},
			{ID: "notes", Label: "Notes", Type: schema.FieldMarkdown, UIHint: formset.UIHintMarkdown},
			{ID: "summary", Label: "Summary", Type: schema.FieldTextarea},
			{ID: "price", Label: "Price", Type: schema.FieldNumber, Indexed: true, Rules: []schema.ValidationRule{{Name: schema.RuleMoney}}},
			{ID: "featured", Label: "Featured", Type: schema.FieldBoolean, Indexed: true},
			{ID: "status", Label: "Status", Type: schema.FieldSelect, Indexed: true, Options: []schema.Option{{Value: "draft", Label: "Draft"}}},
			{ID: "published", Label: "Published", Type: schema.FieldDateTime, Indexed: true},
			{ID: "cover", Label: "Cover", Type: schema.FieldRelation, UIHint: schema.UIHintMedia, Indexed: true},
			{ID: "author", Label: "Author", Type: schema.FieldRelation},
			{ID: "secret", Label: "Secret", Type: schema.FieldString, Sensitive: true, Indexed: true},
			{ID: "meta", Label: "Meta", Type: schema.FieldJSON, Rules: []schema.ValidationRule{{Name: schema.RuleJSONAny}}},
			{ID: "seo", Label: "SEO", Type: schema.FieldObject, Fields: []schema.Field{{ID: "title", Label: "SEO title", Type: schema.FieldString}}},
			{ID: "tags", Label: "Tags", Type: schema.FieldCollection, Items: &schema.Field{ID: "tag", Label: "Tag", Type: schema.FieldSelect, Options: []schema.Option{{Value: "go", Label: "Go"}}}},
			{ID: "links", Label: "Links", Type: schema.FieldCollection, Items: &schema.Field{ID: "link", Label: "Link", Type: schema.FieldObject, Fields: []schema.Field{{ID: "url", Label: "URL", Type: schema.FieldString}}}},
			{ID: "aliases", Label: "Aliases", Type: schema.FieldCollection, Items: &schema.Field{ID: "alias", Label: "Alias", Type: schema.FieldString}},
			{ID: "score", Label: "Score", Type: schema.FieldComputed},
			{ID: "slug", Label: "Slug", Type: schema.FieldString, Rules: []schema.ValidationRule{{Name: schema.RuleReadOnly}}},
		},
		Relations: []schema.Relation{
			{ID: "author", Label: "Author", Source: "article", Target: "user", Cardinality: schema.RelationOneToOne},
			{ID: "cover", Source: "article", Target: "media", Cardinality: schema.RelationOneToMany},
		},
	}
}

func TestResourceFromRecordProjectsCodexRecord(t *testing.T) {
	resource := ResourceFromRecord(articleRecord(), ResourceOptions[testCapability]{
		Capabilities: map[ResourceOperation]testCapability{OperationList: "read", OperationEdit: "write"},
		TrailColumns: []Column{{ID: "updated_at", Label: "Updated", Type: ColumnDateTime}},
		MaxColumns:   10,
	})
	if err := resource.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
	if resource.BasePath != "/article" || resource.Description != "Blog" || len(resource.Capabilities) != 2 {
		t.Fatalf("resource: %#v", resource)
	}
	widgets := map[string]FieldType{}
	byID := map[string]Field{}
	for _, field := range resource.Form.Fields {
		widgets[field.ID] = field.Type
		byID[field.ID] = field
	}
	want := map[string]FieldType{
		"title": FieldText, "body": FieldRichText, "notes": FieldMarkdown, "summary": FieldTextarea, "price": FieldNumber,
		"featured": FieldBoolean, "status": FieldSelect, "published": FieldDateTime, "cover": FieldFile,
		"author": FieldRelation, "meta": FieldJSON, "seo": FieldObject, "tags": FieldMultiSelect, "links": FieldRepeater,
		"aliases": FieldRepeater, "score": FieldDisplay, "slug": FieldText,
	}
	for id, widget := range want {
		if widgets[id] != widget {
			t.Errorf("%s widget = %q want %q", id, widgets[id], widget)
		}
	}
	if byID["price"].DataType != schema.FieldNumber || byID["price"].Rules[0].Name != schema.RuleMoney {
		t.Fatalf("price: %#v", byID["price"])
	}
	if !byID["title"].Localized || !byID["title"].Required || !byID["slug"].ReadOnly || !byID["score"].ReadOnly {
		t.Fatal("flags not projected")
	}
	if byID["author"].Relation.ResourceID != "user" || byID["author"].Relation.Multiple || !byID["cover"].Relation.Multiple {
		t.Fatalf("relations: %#v %#v", byID["author"].Relation, byID["cover"].Relation)
	}
	if len(byID["tags"].Options) != 1 || len(byID["links"].Fields) != 1 || len(byID["aliases"].Fields) != 1 || len(byID["seo"].Fields) != 1 {
		t.Fatal("nested projection")
	}
	columns := map[string]ColumnType{}
	for _, column := range resource.Table.Columns {
		columns[column.ID] = column.Type
	}
	if _, leaked := columns["secret"]; leaked {
		t.Fatal("sensitive field must not become a column")
	}
	wantColumns := map[string]ColumnType{"title": ColumnText, "price": ColumnNumber, "featured": ColumnBoolean, "status": ColumnBadge, "published": ColumnDateTime, "cover": ColumnImage, "updated_at": ColumnDateTime}
	for id, columnType := range wantColumns {
		if columns[id] != columnType {
			t.Errorf("column %s = %q want %q", id, columns[id], columnType)
		}
	}
	if !resource.Table.Searchable || len(resource.Table.Filters) != 3 || resource.Table.DefaultSort.ColumnID != "title" {
		t.Fatalf("table: %#v", resource.Table)
	}
	if len(resource.Relations) != 2 || resource.Relations[0].Type != "one-to-one" {
		t.Fatalf("relations: %#v", resource.Relations)
	}
}

func TestTableFromRecordFallbackAndLimit(t *testing.T) {
	record := schema.RecordType{ID: "note", Label: "Notes", Scope: schema.ScopeWorkspace, Fields: []schema.Field{
		{ID: "flag", Label: "Flag", Type: schema.FieldBoolean},
		{ID: "name", Label: "Name", Type: schema.FieldString},
	}}
	table := TableFromRecord(record, ResourceOptions[string]{LeadColumns: []Column{{ID: "id", Type: ColumnText}}})
	if len(table.Columns) != 2 || table.Columns[1].ID != "name" || table.DefaultSort.ColumnID != "id" {
		t.Fatalf("fallback: %#v", table.Columns)
	}
	resource := ResourceFromRecord(articleRecord(), ResourceOptions[string]{BasePath: "/a", MaxColumns: 1})
	if len(resource.Table.Columns) != 1 || resource.BasePath != "/a" {
		t.Fatalf("limit: %#v", resource.Table.Columns)
	}
	if empty := TableFromRecord(schema.RecordType{ID: "x"}, ResourceOptions[string]{}); len(empty.Columns) != 0 {
		t.Fatal("empty")
	}
}

func TestRecordValidatorUsesFormsetAndCodexRules(t *testing.T) {
	validator := RecordValidator{Record: schema.RecordType{ID: "offer", Label: "Offer", Scope: schema.ScopeWorkspace, Fields: []schema.Field{
		{ID: "title", Label: "Title", Type: schema.FieldText, Required: true},
		{ID: "price", Label: "Price", Type: schema.FieldNumber, Rules: []schema.ValidationRule{{Name: schema.RuleMoney}}},
	}}}
	ctx := context.Background()
	if issues := validator.Validate(ctx, Command{Operation: OperationCreate, Data: Record{"title": "A", "price": 100.0}}); len(issues) != 0 {
		t.Fatalf("valid: %#v", issues)
	}
	issues := validator.Validate(ctx, Command{Operation: OperationEdit, Data: Record{"price": 1.5}})
	if len(issues) != 2 {
		t.Fatalf("issues: %#v", issues)
	}
	if validator.Validate(ctx, Command{Operation: OperationDelete}) != nil {
		t.Fatal("delete is not validated")
	}
	broken := RecordValidator{Record: schema.RecordType{ID: "x"}, Locale: "ru"}
	if issues := broken.Validate(ctx, Command{Operation: OperationCreate}); len(issues) != 1 || issues[0].Field != "" {
		t.Fatalf("broken record: %#v", issues)
	}
}
