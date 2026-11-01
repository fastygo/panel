package panel

import "github.com/fastygo/codex/schema"

// RecordTypeID and RecordType are the Codex record vocabulary. panel owns no
// record model; project one into a Resource with ResourceFromRecord.
type (
	RecordTypeID = schema.RecordTypeID
	RecordType   = schema.RecordType
)

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
