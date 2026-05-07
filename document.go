package panel

import (
	"fmt"
	"strings"
)

type RecordTypeID string

type RecordType struct {
	ID          RecordTypeID
	Label       string
	Description string
	Fields      []Field
	Form        FormSchema
	Table       TableSchema[string]
	Detail      DetailSchema
	Permissions []RecordPermission
}

func (t RecordType) Validate() error {
	if strings.TrimSpace(string(t.ID)) == "" {
		return fmt.Errorf("record type id is required")
	}
	if strings.TrimSpace(t.Label) == "" {
		return fmt.Errorf("record type %q label is required", t.ID)
	}
	for _, field := range t.Fields {
		if err := field.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type RecordPermission struct {
	Role       string
	Operation  ResourceOperation
	Capability string
}

type WorkflowState struct {
	ID          string
	Label       string
	Description string
	Initial     bool
	Terminal    bool
	Color       string
}

type WorkflowTransition[C comparable] struct {
	ID         string
	Label      string
	From       string
	To         string
	Capability C
	Action     Action[C]
}

type Workflow[C comparable] struct {
	ID          string
	Label       string
	Description string
	States      []WorkflowState
	Transitions []WorkflowTransition[C]
}

type ActivityKind string

const (
	ActivityComment ActivityKind = "comment"
	ActivityAudit   ActivityKind = "audit"
	ActivityStatus  ActivityKind = "status"
	ActivityMessage ActivityKind = "message"
)

type TimelineEvent struct {
	ID          string
	Kind        ActivityKind
	ActorID     string
	Title       string
	Description string
	CreatedAt   string
	Metadata    map[string]string
}

type Timeline struct {
	ID     string
	Label  string
	Events []TimelineEvent
}

type Report[C comparable] struct {
	ID          string
	Label       string
	Description string
	Filters     []Filter
	Columns     []Column
	Actions     []Action[C]
	Exportable  bool
}

type ImportExportDescriptor[C comparable] struct {
	ID          string
	Label       string
	Formats     []string
	Import      Action[C]
	Export      Action[C]
	Capability  C
	Description string
}

type PrintDescriptor[C comparable] struct {
	ID          string
	Label       string
	Description string
	Format      string
	Capability  C
}
