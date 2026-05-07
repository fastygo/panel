package panel

import (
	"fmt"
	"strings"
)

type ResourceID string

type ResourceOperation string

const (
	OperationList   ResourceOperation = "list"
	OperationCreate ResourceOperation = "create"
	OperationEdit   ResourceOperation = "edit"
	OperationView   ResourceOperation = "view"
	OperationDelete ResourceOperation = "delete"
	OperationBulk   ResourceOperation = "bulk"
	OperationExport ResourceOperation = "export"
	OperationImport ResourceOperation = "import"
)

type ResourceRouteRole string

const (
	RouteIndex  ResourceRouteRole = "index"
	RouteNew    ResourceRouteRole = "new"
	RouteCreate ResourceRouteRole = "create"
	RouteEdit   ResourceRouteRole = "edit"
	RouteUpdate ResourceRouteRole = "update"
	RouteView   ResourceRouteRole = "view"
	RouteDelete ResourceRouteRole = "delete"
)

type ResourceRoute[C comparable] struct {
	Role       ResourceRouteRole
	Pattern    string
	Capability C
}

type ResourceCapability[C comparable] struct {
	Operation  ResourceOperation
	Capability C
}

type ResourceRelation struct {
	ID         string
	Label      string
	ResourceID ResourceID
	Type       string
}

type Resource[C comparable] struct {
	ID           ResourceID
	Label        string
	Singular     string
	Plural       string
	BasePath     string
	Description  string
	Icon         string
	Navigation   MenuItem[C]
	Capabilities []ResourceCapability[C]
	Routes       []ResourceRoute[C]
	Table        TableSchema[C]
	Form         FormSchema
	Detail       DetailSchema
	Relations    []ResourceRelation
	Actions      []Action[C]
}

func (r Resource[C]) Validate() error {
	if strings.TrimSpace(string(r.ID)) == "" {
		return fmt.Errorf("resource id is required")
	}
	if strings.TrimSpace(r.Label) == "" {
		return fmt.Errorf("resource %q label is required", r.ID)
	}
	if strings.TrimSpace(r.BasePath) == "" {
		return fmt.Errorf("resource %q base path is required", r.ID)
	}
	if err := r.Table.Validate(); err != nil {
		return err
	}
	if err := r.Form.Validate(); err != nil {
		return err
	}
	if err := r.Detail.Validate(); err != nil {
		return err
	}
	for _, action := range r.Actions {
		if err := action.Validate(); err != nil {
			return err
		}
	}
	return nil
}
