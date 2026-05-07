package panel

import (
	"fmt"
	"strings"
)

type PanelID string

type PanelOptions[C comparable] struct {
	ID          PanelID
	Title       string
	BasePath    string
	Description string
	Icon        string
	Assets      []Asset
	Navigation  []NavigationGroup[C]
}

type Panel[P Principal[C], C comparable] struct {
	id          PanelID
	title       string
	basePath    string
	description string
	icon        string
	registry    *Registry[P, C]
	resources   []Resource[C]
	pages       []Page[C]
	widgets     []Widget[C]
}

func NewPanel[P Principal[C], C comparable](options PanelOptions[C]) (*Panel[P, C], error) {
	if strings.TrimSpace(string(options.ID)) == "" {
		return nil, fmt.Errorf("panel id is required")
	}
	if strings.TrimSpace(options.Title) == "" {
		return nil, fmt.Errorf("panel %q title is required", options.ID)
	}
	if strings.TrimSpace(options.BasePath) == "" {
		return nil, fmt.Errorf("panel %q base path is required", options.ID)
	}
	registry := NewRegistry[P, C]()
	registry.AddAssets(options.Assets...)
	registry.AddNavigationGroups(options.Navigation...)
	return &Panel[P, C]{
		id:          options.ID,
		title:       options.Title,
		basePath:    options.BasePath,
		description: options.Description,
		icon:        options.Icon,
		registry:    registry,
	}, nil
}

func (p *Panel[P, C]) ID() PanelID {
	return p.id
}

func (p *Panel[P, C]) Title() string {
	return p.title
}

func (p *Panel[P, C]) BasePath() string {
	return p.basePath
}

func (p *Panel[P, C]) Description() string {
	return p.description
}

func (p *Panel[P, C]) Icon() string {
	return p.icon
}

func (p *Panel[P, C]) Registry() *Registry[P, C] {
	return p.registry
}

func (p *Panel[P, C]) AddResources(resources ...Resource[C]) error {
	for _, resource := range resources {
		if err := resource.Validate(); err != nil {
			return err
		}
		p.resources = append(p.resources, resource)
		if resource.Navigation.Path != "" {
			p.registry.AddMenuItems(resource.Navigation)
		}
	}
	return nil
}

func (p *Panel[P, C]) Resources() []Resource[C] {
	return append([]Resource[C](nil), p.resources...)
}

func (p *Panel[P, C]) AddPages(pages ...Page[C]) error {
	for _, page := range pages {
		if err := page.Validate(); err != nil {
			return err
		}
		p.pages = append(p.pages, page)
		if page.Navigation.Path != "" {
			p.registry.AddMenuItems(page.Navigation)
		}
	}
	return nil
}

func (p *Panel[P, C]) Pages() []Page[C] {
	return append([]Page[C](nil), p.pages...)
}

func (p *Panel[P, C]) AddWidgets(widgets ...Widget[C]) error {
	for _, widget := range widgets {
		if err := widget.Validate(); err != nil {
			return err
		}
		p.widgets = append(p.widgets, widget)
	}
	return nil
}

func (p *Panel[P, C]) Widgets() []Widget[C] {
	return append([]Widget[C](nil), p.widgets...)
}
