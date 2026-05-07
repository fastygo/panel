package panel

import (
	"fmt"
	"strings"
)

type ColumnType string

const (
	ColumnText     ColumnType = "text"
	ColumnNumber   ColumnType = "number"
	ColumnBoolean  ColumnType = "boolean"
	ColumnDateTime ColumnType = "datetime"
	ColumnBadge    ColumnType = "badge"
	ColumnImage    ColumnType = "image"
	ColumnProgress ColumnType = "progress"
	ColumnActions  ColumnType = "actions"
)

type Column struct {
	ID          string
	Label       string
	Type        ColumnType
	Sortable    bool
	Searchable  bool
	Toggleable  bool
	Description string
	Width       string
}

func (c Column) Validate() error {
	if strings.TrimSpace(c.ID) == "" {
		return fmt.Errorf("column id is required")
	}
	if strings.TrimSpace(string(c.Type)) == "" {
		return fmt.Errorf("column %q type is required", c.ID)
	}
	return nil
}

type FilterType string

const (
	FilterText      FilterType = "text"
	FilterSelect    FilterType = "select"
	FilterBoolean   FilterType = "boolean"
	FilterDateRange FilterType = "date-range"
)

type Filter struct {
	ID          string
	Label       string
	Type        FilterType
	Options     []Option
	Default     string
	Description string
}

func (f Filter) Validate() error {
	if strings.TrimSpace(f.ID) == "" {
		return fmt.Errorf("filter id is required")
	}
	if strings.TrimSpace(string(f.Type)) == "" {
		return fmt.Errorf("filter %q type is required", f.ID)
	}
	return nil
}

type SortDirection string

const (
	SortAsc  SortDirection = "asc"
	SortDesc SortDirection = "desc"
)

type Sort struct {
	ColumnID  string
	Direction SortDirection
}

type TableSchema[C comparable] struct {
	ID          string
	Columns     []Column
	Filters     []Filter
	DefaultSort Sort
	Searchable  bool
	Exportable  bool
	PerPage     []int
	RowActions  []Action[C]
	BulkActions []Action[C]
}

func (s TableSchema[C]) Validate() error {
	for _, column := range s.Columns {
		if err := column.Validate(); err != nil {
			return err
		}
	}
	for _, filter := range s.Filters {
		if err := filter.Validate(); err != nil {
			return err
		}
	}
	for _, action := range s.RowActions {
		if err := action.Validate(); err != nil {
			return err
		}
	}
	for _, action := range s.BulkActions {
		if err := action.Validate(); err != nil {
			return err
		}
	}
	return nil
}
