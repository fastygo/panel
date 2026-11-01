package panel

import (
	"fmt"
	"strings"

	"github.com/fastygo/codex/schema"
)

// FieldType is a panel UI widget. Data types come from Codex: see
// Field.DataType and ResourceFromRecord. Widgets that share a name with a
// Codex field type use the same string.
type FieldType string

const (
	FieldText         FieldType = "text"
	FieldTextarea     FieldType = "textarea"
	FieldRichText     FieldType = "richtext"
	FieldMarkdown     FieldType = "markdown"
	FieldNumber       FieldType = "number"
	FieldBoolean      FieldType = "boolean"
	FieldSelect       FieldType = "select"
	FieldMultiSelect  FieldType = "multiselect"
	FieldRelation     FieldType = "relation"
	FieldRepeater     FieldType = "repeater"
	FieldFile         FieldType = "file"
	FieldDateTime     FieldType = "datetime"
	FieldJSON         FieldType = "json"
	FieldCode         FieldType = "code"
	FieldHidden       FieldType = "hidden"
	FieldDisplay      FieldType = "display"
	FieldStatus       FieldType = "status"
	FieldTimeline     FieldType = "timeline"
	FieldRelationship FieldType = "relationship"
	FieldObject       FieldType = "object"
)

type SchemaLayout string

const (
	SchemaStack SchemaLayout = "stack"
	SchemaGrid  SchemaLayout = "grid"
	SchemaTabs  SchemaLayout = "tabs"
)

// ValidationRule is the Codex validation rule (fastygo.codex/* names).
type ValidationRule = schema.ValidationRule

type Option struct {
	Value       string
	Label       string
	Description string
	Icon        string
	Disabled    bool
}

type Relation struct {
	ID            string
	Label         string
	ResourceID    string
	DisplayColumn string
	Multiple      bool
	Searchable    bool
}

type Field struct {
	ID           string
	Label        string
	Type         FieldType
	Description  string
	Placeholder  string
	Required     bool
	ReadOnly     bool
	Hidden       bool
	DefaultValue string
	Options      []Option
	Relation     Relation
	Fields       []Field
	Rules        []ValidationRule
	ColumnSpan   int
	// DataType is the Codex field type this widget edits, when the field was
	// projected from a Codex record (empty for UI-only fields).
	DataType schema.FieldType
	// Localized mirrors the Codex field: the value is stored per locale.
	Localized bool
}

func (f Field) Validate() error {
	if strings.TrimSpace(f.ID) == "" {
		return fmt.Errorf("field id is required")
	}
	if strings.TrimSpace(string(f.Type)) == "" {
		return fmt.Errorf("field %q type is required", f.ID)
	}
	for _, child := range f.Fields {
		if err := child.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type SchemaSection struct {
	ID          string
	Label       string
	Description string
	Layout      SchemaLayout
	Columns     int
	Fields      []Field
}

func (s SchemaSection) Validate() error {
	if strings.TrimSpace(s.ID) == "" {
		return fmt.Errorf("schema section id is required")
	}
	for _, field := range s.Fields {
		if err := field.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type FormSchema struct {
	ID          string
	Operation   string
	Sections    []SchemaSection
	Fields      []Field
	Description string
}

func (s FormSchema) Validate() error {
	for _, section := range s.Sections {
		if err := section.Validate(); err != nil {
			return err
		}
	}
	for _, field := range s.Fields {
		if err := field.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type DetailSchema struct {
	ID       string
	Sections []SchemaSection
	Fields   []Field
}

func (s DetailSchema) Validate() error {
	for _, section := range s.Sections {
		if err := section.Validate(); err != nil {
			return err
		}
	}
	for _, field := range s.Fields {
		if err := field.Validate(); err != nil {
			return err
		}
	}
	return nil
}
