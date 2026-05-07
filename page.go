package panel

import (
	"fmt"
	"strings"
)

type PageID string

type PageKind string

const (
	PageDashboard PageKind = "dashboard"
	PageSettings  PageKind = "settings"
	PageReport    PageKind = "report"
	PageWizard    PageKind = "wizard"
	PageRuntime   PageKind = "runtime"
	PageCustom    PageKind = "custom"
)

type Page[C comparable] struct {
	ID          PageID
	Kind        PageKind
	Title       string
	Description string
	Path        string
	Icon        string
	Navigation  MenuItem[C]
	Capability  C
	Actions     []Action[C]
	Widgets     []Widget[C]
	Form        FormSchema
	Table       TableSchema[C]
}

func (p Page[C]) Validate() error {
	if strings.TrimSpace(string(p.ID)) == "" {
		return fmt.Errorf("page id is required")
	}
	if strings.TrimSpace(p.Title) == "" {
		return fmt.Errorf("page %q title is required", p.ID)
	}
	if strings.TrimSpace(p.Path) == "" {
		return fmt.Errorf("page %q path is required", p.ID)
	}
	if err := p.Form.Validate(); err != nil {
		return err
	}
	if err := p.Table.Validate(); err != nil {
		return err
	}
	for _, action := range p.Actions {
		if err := action.Validate(); err != nil {
			return err
		}
	}
	for _, widget := range p.Widgets {
		if err := widget.Validate(); err != nil {
			return err
		}
	}
	return nil
}
