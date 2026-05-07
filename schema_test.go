package panel

import "testing"

func TestSchemasValidateNestedFieldsAndActions(t *testing.T) {
	form := FormSchema{
		Sections: []SchemaSection{
			{
				ID: "main",
				Fields: []Field{
					{ID: "title", Type: FieldText},
					{ID: "items", Type: FieldRepeater, Fields: []Field{{ID: "name", Type: FieldText}}},
				},
			},
		},
	}
	if err := form.Validate(); err != nil {
		t.Fatalf("FormSchema.Validate() error = %v", err)
	}

	table := TableSchema[testCapability]{
		Columns: []Column{{ID: "title", Type: ColumnText}},
		Filters: []Filter{{ID: "status", Type: FilterSelect}},
		RowActions: []Action[testCapability]{
			{ID: "edit", Label: "Edit", Placement: ActionRow},
		},
	}
	if err := table.Validate(); err != nil {
		t.Fatalf("TableSchema.Validate() error = %v", err)
	}
}

func TestSchemasRejectInvalidDescriptors(t *testing.T) {
	if err := (Field{Type: FieldText}).Validate(); err == nil {
		t.Fatalf("expected invalid field without id")
	}
	if err := (Column{Type: ColumnText}).Validate(); err == nil {
		t.Fatalf("expected invalid column without id")
	}
	if err := (Action[testCapability]{ID: "missing-label"}).Validate(); err == nil {
		t.Fatalf("expected invalid action without label")
	}
}
