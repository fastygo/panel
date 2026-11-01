package panel

import (
	"context"
	"strings"

	"github.com/fastygo/codex/schema"
	"github.com/fastygo/formset"
)

// ResourceOptions tunes ResourceFromRecord. Every field is optional.
type ResourceOptions[C comparable] struct {
	// BasePath defaults to "/" + record ID.
	BasePath string
	Icon     string
	// Capabilities maps an operation to the capability that guards it.
	Capabilities map[ResourceOperation]C
	// LeadColumns are prepended (for example id or status chrome columns).
	LeadColumns []Column
	// TrailColumns are appended (for example updated_at).
	TrailColumns []Column
	// MaxColumns caps projected field columns; 0 means 6.
	MaxColumns int
}

// ResourceFromRecord projects a Codex record type (or a formset record,
// which is the same type) into a panel Resource: form fields, table columns
// from Indexed/Searchable fields, filters and relations.
// The record is projected through formset.FromCodex so callers never share
// slices with their manifest.
func ResourceFromRecord[C comparable](record schema.RecordType, options ResourceOptions[C]) Resource[C] {
	record = formset.FromCodex(record)
	basePath := options.BasePath
	if basePath == "" {
		basePath = "/" + string(record.ID)
	}
	resource := Resource[C]{
		ID:          ResourceID(record.ID),
		Label:       record.Label,
		Singular:    record.Label,
		Plural:      record.Label,
		BasePath:    basePath,
		Description: record.Description,
		Icon:        options.Icon,
		Form:        FormSchema{ID: string(record.ID), Fields: FieldsFromRecord(record)},
		Detail:      DetailSchema{ID: string(record.ID), Fields: FieldsFromRecord(record)},
		Table:       TableFromRecord[C](record, options),
	}
	for _, operation := range []ResourceOperation{OperationList, OperationCreate, OperationEdit, OperationView, OperationDelete, OperationBulk, OperationExport, OperationImport} {
		if capability, ok := options.Capabilities[operation]; ok {
			resource.Capabilities = append(resource.Capabilities, ResourceCapability[C]{Operation: operation, Capability: capability})
		}
	}
	for _, relation := range record.Relations {
		resource.Relations = append(resource.Relations, ResourceRelation{
			ID: string(relation.ID), Label: relation.Label,
			ResourceID: ResourceID(relation.Target), Type: string(relation.Cardinality),
		})
	}
	return resource
}

// FieldsFromRecord projects Codex fields into panel form fields.
func FieldsFromRecord(record schema.RecordType) []Field {
	fields := make([]Field, 0, len(record.Fields))
	for _, field := range record.Fields {
		fields = append(fields, fieldFromCodex(field, relationFor(record, field.ID)))
	}
	return fields
}

func relationFor(record schema.RecordType, id schema.FieldID) *schema.Relation {
	for index := range record.Relations {
		if record.Relations[index].ID == schema.RelationID(id) {
			return &record.Relations[index]
		}
	}
	return nil
}

// WidgetFor returns the default panel widget for a Codex field.
func WidgetFor(field schema.Field) FieldType {
	switch field.UIHint {
	case schema.UIHintMedia:
		return FieldFile
	case formset.UIHintTipTap:
		return FieldRichText
	case formset.UIHintMarkdown:
		return FieldMarkdown
	case formset.UIHintTextarea:
		return FieldTextarea
	}
	switch field.Type {
	case schema.FieldTextarea:
		return FieldTextarea
	case schema.FieldRichText:
		return FieldRichText
	case schema.FieldMarkdown:
		return FieldMarkdown
	case schema.FieldNumber:
		return FieldNumber
	case schema.FieldBoolean:
		return FieldBoolean
	case schema.FieldSelect:
		return FieldSelect
	case schema.FieldDateTime:
		return FieldDateTime
	case schema.FieldJSON:
		return FieldJSON
	case schema.FieldObject:
		return FieldObject
	case schema.FieldRelation:
		return FieldRelation
	case schema.FieldCollection:
		if field.Items != nil && field.Items.Type == schema.FieldSelect {
			return FieldMultiSelect
		}
		return FieldRepeater
	case schema.FieldComputed:
		return FieldDisplay
	default:
		return FieldText
	}
}

func fieldFromCodex(field schema.Field, relation *schema.Relation) Field {
	_, readOnly := schema.FieldRule(field, schema.RuleReadOnly)
	projected := Field{
		ID:           string(field.ID),
		Label:        field.Label,
		Type:         WidgetFor(field),
		Description:  field.Description,
		Required:     field.Required,
		ReadOnly:     readOnly || field.Type == schema.FieldComputed,
		DefaultValue: field.DefaultValue,
		Rules:        field.Rules,
		DataType:     field.Type,
		Localized:    field.Localized,
	}
	for _, option := range field.Options {
		projected.Options = append(projected.Options, Option{Value: option.Value, Label: option.Label})
	}
	for _, nested := range field.Fields {
		projected.Fields = append(projected.Fields, fieldFromCodex(nested, nil))
	}
	if field.Items != nil {
		item := *field.Items
		switch item.Type {
		case schema.FieldSelect:
			for _, option := range item.Options {
				projected.Options = append(projected.Options, Option{Value: option.Value, Label: option.Label})
			}
		case schema.FieldObject:
			for _, nested := range item.Fields {
				projected.Fields = append(projected.Fields, fieldFromCodex(nested, nil))
			}
		default:
			projected.Fields = append(projected.Fields, fieldFromCodex(item, nil))
		}
	}
	if relation != nil {
		projected.Relation = Relation{
			ID: string(relation.ID), Label: relation.Label, ResourceID: string(relation.Target),
			DisplayColumn: "title", Multiple: relation.Cardinality != schema.RelationOneToOne,
			Searchable: true,
		}
	}
	return projected
}

// TableFromRecord projects list columns and filters. Columns are fields that
// are Indexed or Searchable (never Sensitive), in manifest order; when none
// qualify the first text-like field is used. Filters come from indexed
// select, boolean and datetime fields.
func TableFromRecord[C comparable](record schema.RecordType, options ResourceOptions[C]) TableSchema[C] {
	limit := options.MaxColumns
	if limit <= 0 {
		limit = 6
	}
	table := TableSchema[C]{ID: string(record.ID), Columns: append([]Column(nil), options.LeadColumns...)}
	projected := 0
	for _, field := range record.Fields {
		if field.Sensitive || field.Encrypted || (!field.Indexed && !field.Searchable) {
			continue
		}
		if field.Searchable {
			table.Searchable = true
		}
		if projected < limit {
			table.Columns = append(table.Columns, columnFromCodex(field))
			projected++
		}
		if filter, ok := filterFromCodex(field); ok && field.Indexed {
			table.Filters = append(table.Filters, filter)
		}
	}
	if projected == 0 {
		for _, field := range record.Fields {
			if !field.Sensitive && isTextual(field.Type) {
				table.Columns = append(table.Columns, columnFromCodex(field))
				break
			}
		}
	}
	table.Columns = append(table.Columns, options.TrailColumns...)
	if len(table.Columns) > 0 {
		table.DefaultSort = Sort{ColumnID: table.Columns[0].ID, Direction: SortAsc}
	}
	return table
}

func isTextual(fieldType schema.FieldType) bool {
	switch fieldType {
	case schema.FieldText, schema.FieldString, schema.FieldTextarea:
		return true
	default:
		return false
	}
}

func columnFromCodex(field schema.Field) Column {
	column := Column{ID: string(field.ID), Label: field.Label, Type: ColumnText, Sortable: field.Indexed, Searchable: field.Searchable, Description: field.Description}
	switch {
	case field.UIHint == schema.UIHintMedia:
		column.Type = ColumnImage
	case field.Type == schema.FieldNumber:
		column.Type = ColumnNumber
	case field.Type == schema.FieldBoolean:
		column.Type = ColumnBoolean
	case field.Type == schema.FieldDateTime:
		column.Type = ColumnDateTime
	case field.Type == schema.FieldSelect:
		column.Type = ColumnBadge
	}
	return column
}

func filterFromCodex(field schema.Field) (Filter, bool) {
	switch field.Type {
	case schema.FieldSelect:
		filter := Filter{ID: string(field.ID), Label: field.Label, Type: FilterSelect}
		for _, option := range field.Options {
			filter.Options = append(filter.Options, Option{Value: option.Value, Label: option.Label})
		}
		return filter, true
	case schema.FieldBoolean:
		return Filter{ID: string(field.ID), Label: field.Label, Type: FilterBoolean}, true
	case schema.FieldDateTime:
		return Filter{ID: string(field.ID), Label: field.Label, Type: FilterDateRange}, true
	default:
		return Filter{}, false
	}
}

// RecordValidator validates create/edit command data against a Codex record
// with formset (Codex rules included). It implements Validator.
type RecordValidator struct {
	Record schema.RecordType
	// Locale defaults to "en".
	Locale string
}

var _ Validator = RecordValidator{}

// Validate reports record or field issues for create and edit commands.
func (validator RecordValidator) Validate(_ context.Context, command Command) []ValidationError {
	if command.Operation != OperationCreate && command.Operation != OperationEdit {
		return nil
	}
	locale := strings.TrimSpace(validator.Locale)
	if locale == "" {
		locale = "en"
	}
	form, err := formset.BindLocale(validator.Record, locale, command.Data)
	if err != nil {
		return []ValidationError{{Message: err.Error()}}
	}
	if len(form.Issues) == 0 {
		return nil
	}
	issues := make([]ValidationError, 0, len(form.Issues))
	for _, issue := range form.Issues {
		issues = append(issues, ValidationError{Field: issue.Field, Message: issue.Message})
	}
	return issues
}
