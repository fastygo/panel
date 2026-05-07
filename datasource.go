package panel

import "context"

type RecordID string

type Record map[string]any

type Query struct {
	Page     int
	PerPage  int
	Search   string
	Filters  map[string]any
	Sort     Sort
	Fields   []string
	Relation string
}

type PageResult struct {
	Records    []Record
	Page       int
	PerPage    int
	TotalItems int
	TotalPages int
}

type Command struct {
	Operation ResourceOperation
	ID        RecordID
	Data      Record
	Metadata  map[string]any
}

type ValidationError struct {
	Field   string
	Message string
}

type CommandResult struct {
	ID       RecordID
	Record   Record
	Warnings []string
}

type RecordProvider interface {
	List(context.Context, Query) (PageResult, error)
	Get(context.Context, RecordID) (Record, error)
}

type CommandHandler interface {
	Handle(context.Context, Command) (CommandResult, error)
}

type Validator interface {
	Validate(context.Context, Command) []ValidationError
}

type TransactionManager interface {
	WithinTransaction(context.Context, func(context.Context) error) error
}

type DataSource interface {
	RecordProvider
	CommandHandler
}
