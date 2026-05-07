package panel

import (
	"fmt"
	"strings"
)

type ActionPlacement string

const (
	ActionHeader ActionPlacement = "header"
	ActionRow    ActionPlacement = "row"
	ActionBulk   ActionPlacement = "bulk"
	ActionModal  ActionPlacement = "modal"
	ActionInline ActionPlacement = "inline"
)

type ActionStyle string

const (
	ActionButton     ActionStyle = "button"
	ActionLink       ActionStyle = "link"
	ActionIconButton ActionStyle = "icon-button"
	ActionBadge      ActionStyle = "badge"
)

type Action[C comparable] struct {
	ID                   string
	Label                string
	Description          string
	Icon                 string
	Placement            ActionPlacement
	Style                ActionStyle
	Color                string
	URL                  string
	Capability           C
	RequiresConfirmation bool
	Disabled             bool
	ModalTitle           string
	ModalDescription     string
	Form                 FormSchema
}

func (a Action[C]) Validate() error {
	if strings.TrimSpace(a.ID) == "" {
		return fmt.Errorf("action id is required")
	}
	if strings.TrimSpace(a.Label) == "" {
		return fmt.Errorf("action %q label is required", a.ID)
	}
	return nil
}
