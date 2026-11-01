package panel

import (
	"testing"

	"github.com/fastygo/codex/schema"
	"github.com/fastygo/formset"
)

func TestValidateRejectsInvalidNestedDescriptors(t *testing.T) {
	badField := Field{Type: FieldText}
	badAction := Action[testCapability]{ID: "a"}
	badColumn := Column{ID: "c"}
	badFilter := Filter{ID: "f"}
	cases := map[string]interface{ Validate() error }{
		"field-type":      Field{ID: "x"},
		"field-nested":    Field{ID: "x", Type: FieldObject, Fields: []Field{badField}},
		"section-id":      SchemaSection{},
		"section-field":   SchemaSection{ID: "s", Fields: []Field{badField}},
		"form-section":    FormSchema{Sections: []SchemaSection{{}}},
		"form-field":      FormSchema{Fields: []Field{badField}},
		"detail-section":  DetailSchema{Sections: []SchemaSection{{}}},
		"detail-field":    DetailSchema{Fields: []Field{badField}},
		"column-type":     badColumn,
		"filter-id":       Filter{Type: FilterText},
		"filter-type":     badFilter,
		"table-column":    TableSchema[testCapability]{Columns: []Column{badColumn}},
		"table-filter":    TableSchema[testCapability]{Filters: []Filter{badFilter}},
		"table-row":       TableSchema[testCapability]{RowActions: []Action[testCapability]{badAction}},
		"table-bulk":      TableSchema[testCapability]{BulkActions: []Action[testCapability]{badAction}},
		"action-id":       Action[testCapability]{Label: "A"},
		"widget-id":       Widget[testCapability]{},
		"widget-title":    Widget[testCapability]{ID: "w"},
		"widget-action":   Widget[testCapability]{ID: "w", Title: "W", Actions: []Action[testCapability]{badAction}},
		"page-id":         Page[testCapability]{},
		"page-title":      Page[testCapability]{ID: "p"},
		"page-path":       Page[testCapability]{ID: "p", Title: "P"},
		"page-form":       Page[testCapability]{ID: "p", Title: "P", Path: "/p", Form: FormSchema{Fields: []Field{badField}}},
		"page-table":      Page[testCapability]{ID: "p", Title: "P", Path: "/p", Table: TableSchema[testCapability]{Columns: []Column{badColumn}}},
		"page-action":     Page[testCapability]{ID: "p", Title: "P", Path: "/p", Actions: []Action[testCapability]{badAction}},
		"page-widget":     Page[testCapability]{ID: "p", Title: "P", Path: "/p", Widgets: []Widget[testCapability]{{}}},
		"resource-id":     Resource[testCapability]{},
		"resource-label":  Resource[testCapability]{ID: "r"},
		"resource-path":   Resource[testCapability]{ID: "r", Label: "R"},
		"resource-table":  Resource[testCapability]{ID: "r", Label: "R", BasePath: "/r", Table: TableSchema[testCapability]{Columns: []Column{badColumn}}},
		"resource-form":   Resource[testCapability]{ID: "r", Label: "R", BasePath: "/r", Form: FormSchema{Fields: []Field{badField}}},
		"resource-detail": Resource[testCapability]{ID: "r", Label: "R", BasePath: "/r", Detail: DetailSchema{Fields: []Field{badField}}},
		"resource-action": Resource[testCapability]{ID: "r", Label: "R", BasePath: "/r", Actions: []Action[testCapability]{badAction}},
	}
	for name, descriptor := range cases {
		if err := descriptor.Validate(); err == nil {
			t.Errorf("%s: expected error", name)
		}
	}
	valid := Page[testCapability]{ID: "p", Title: "P", Path: "/p", Widgets: []Widget[testCapability]{{ID: "w", Title: "W", Actions: []Action[testCapability]{{ID: "a", Label: "A"}}}}, Actions: []Action[testCapability]{{ID: "a", Label: "A"}}}
	if err := valid.Validate(); err != nil {
		t.Fatal(err)
	}
	if err := (DetailSchema{Sections: []SchemaSection{{ID: "s"}}, Fields: []Field{{ID: "f", Type: FieldText}}}).Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestPanelAccessorsAndOptions(t *testing.T) {
	for _, options := range []PanelOptions[testCapability]{{}, {ID: "x"}, {ID: "x", Title: "X"}} {
		if _, err := NewPanel[testPrincipal, testCapability](options); err == nil {
			t.Fatalf("expected error for %#v", options)
		}
	}
	admin, err := NewPanel[testPrincipal, testCapability](PanelOptions[testCapability]{ID: "admin", Title: "Admin", BasePath: "/admin", Description: "d", Icon: "i"})
	if err != nil {
		t.Fatal(err)
	}
	if admin.ID() != "admin" || admin.Title() != "Admin" || admin.BasePath() != "/admin" || admin.Description() != "d" || admin.Icon() != "i" {
		t.Fatal("accessors")
	}
	if admin.AddResources(Resource[testCapability]{}) == nil || admin.AddPages(Page[testCapability]{}) == nil || admin.AddWidgets(Widget[testCapability]{}) == nil {
		t.Fatal("invalid descriptors must be rejected")
	}
	allow := PolicyFunc[testPrincipal, testCapability](func(_ testPrincipal, c testCapability) bool { return c == "ok" })
	if !allow.Can(newTestPrincipal(), "ok") || allow.Can(newTestPrincipal(), "no") {
		t.Fatal("policy func")
	}
}

func TestWidgetForHints(t *testing.T) {
	cases := map[FieldType]schema.Field{
		FieldTextarea: {Type: schema.FieldString, UIHint: formset.UIHintTextarea},
		FieldRichText: {Type: schema.FieldString, UIHint: formset.UIHintTipTap},
		FieldText:     {Type: schema.FieldEncrypted},
	}
	for want, field := range cases {
		if got := WidgetFor(field); got != want {
			t.Errorf("%#v: %q want %q", field, got, want)
		}
	}
}

func TestNavigationGroupsFilterByCapability(t *testing.T) {
	registry := NewRegistry[testPrincipal, testCapability]()
	registry.AddNavigationGroups(NavigationGroup[testCapability]{ID: "b", Order: 2}, NavigationGroup[testCapability]{ID: "a", Order: 1}, NavigationGroup[testCapability]{ID: "x", Capability: "secret"})
	groups := registry.NavigationGroups(newTestPrincipal())
	if len(groups) != 2 || groups[0].ID != "a" {
		t.Fatalf("groups: %#v", groups)
	}
}
