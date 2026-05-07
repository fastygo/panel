package panel

import (
	"context"
	"testing"
)

type memorySource struct {
	records map[RecordID]Record
}

func (s memorySource) List(context.Context, Query) (PageResult, error) {
	items := make([]Record, 0, len(s.records))
	for _, record := range s.records {
		items = append(items, record)
	}
	return PageResult{Records: items, Page: 1, PerPage: len(items), TotalItems: len(items), TotalPages: 1}, nil
}

func (s memorySource) Get(_ context.Context, id RecordID) (Record, error) {
	return s.records[id], nil
}

func (s memorySource) Handle(_ context.Context, command Command) (CommandResult, error) {
	return CommandResult{ID: command.ID, Record: command.Data}, nil
}

func TestDataSourceContract(t *testing.T) {
	var source DataSource = memorySource{records: map[RecordID]Record{
		"1": {"name": "Example"},
	}}

	result, err := source.List(context.Background(), Query{Page: 1, PerPage: 10})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if result.TotalItems != 1 {
		t.Fatalf("TotalItems = %d, want 1", result.TotalItems)
	}

	record, err := source.Get(context.Background(), "1")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if record["name"] != "Example" {
		t.Fatalf("record = %v, want Example", record)
	}
}
