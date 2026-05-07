package panel

import (
	"fmt"
	"strings"
)

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
)

type SchemaLayout string

const (
	SchemaStack SchemaLayout = "stack"
	SchemaGrid  SchemaLayout = "grid"
	SchemaTabs  SchemaLayout = "tabs"
)

type ValidationRule struct {
	Name    string
	Message string
	Args    map[string]string
}

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
