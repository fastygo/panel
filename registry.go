package panel

import (
	"net/http"
	"sort"
	"strings"
)

type Surface string

const (
	SurfaceAdmin  Surface = "admin"
	SurfaceREST   Surface = "rest"
	SurfacePublic Surface = "public"
)

type NavItem struct {
	Label string
	Path  string
	Icon  string
	Order int
}

type NavigationGroup[C comparable] struct {
	ID         string
	Label      string
	Icon       string
	Order      int
	Collapsed  bool
	Capability C
}

type MenuItem[C comparable] struct {
	ID         string
	Label      string
	Path       string
	Icon       string
	GroupID    string
	Order      int
	Capability C
}

type AssetKind string

const (
	AssetScript     AssetKind = "script"
	AssetStylesheet AssetKind = "stylesheet"
	AssetModule     AssetKind = "module"
)

type Asset struct {
	ID      string
	Kind    AssetKind
	Surface Surface
	Path    string
}

type EditorProviderRegistration struct {
	ID          string
	Label       string
	Description string
	Priority    int
	Formats     []string
}

type Route[P Principal[C], C comparable] struct {
	Pattern          string
	Surface          Surface
	Capability       C
	Protected        bool
	Handler          http.HandlerFunc
	ProtectedHandler func(http.ResponseWriter, *http.Request, P)
}

type Registry[P Principal[C], C comparable] struct {
	navigationGroups []NavigationGroup[C]
	menuItems        []MenuItem[C]
	editorProviders  []EditorProviderRegistration
	routes           []Route[P, C]
	assets           []Asset
}

func NewRegistry[P Principal[C], C comparable]() *Registry[P, C] {
	return &Registry[P, C]{}
}

func (r *Registry[P, C]) AddNavigationGroups(groups ...NavigationGroup[C]) {
	r.navigationGroups = append(r.navigationGroups, groups...)
}

func (r *Registry[P, C]) NavigationGroups(principal P) []NavigationGroup[C] {
	var zero C
	groups := make([]NavigationGroup[C], 0, len(r.navigationGroups))
	for _, group := range r.navigationGroups {
		if group.Capability != zero && !principal.Has(group.Capability) {
			continue
		}
		groups = append(groups, group)
	}
	sort.SliceStable(groups, func(i, j int) bool { return groups[i].Order < groups[j].Order })
	return groups
}

func (r *Registry[P, C]) AddMenuItems(items ...MenuItem[C]) {
	r.menuItems = append(r.menuItems, items...)
}

func (r *Registry[P, C]) MenuItems(principal P) []MenuItem[C] {
	var zero C
	items := make([]MenuItem[C], 0, len(r.menuItems))
	for _, item := range r.menuItems {
		if item.Capability != zero && !principal.Has(item.Capability) {
			continue
		}
		items = append(items, item)
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].Order < items[j].Order })
	return items
}

func (r *Registry[P, C]) NavItems(principal P) []NavItem {
	menu := r.MenuItems(principal)
	items := make([]NavItem, 0, len(menu))
	for _, item := range menu {
		items = append(items, NavItem{Label: item.Label, Path: item.Path, Icon: item.Icon, Order: item.Order})
	}
	return items
}

func (r *Registry[P, C]) AddRoutes(routes ...Route[P, C]) {
	r.routes = append(r.routes, routes...)
}

func (r *Registry[P, C]) RoutesForSurface(surface Surface) []Route[P, C] {
	result := []Route[P, C]{}
	for _, route := range r.routes {
		if route.Surface == surface {
			result = append(result, route)
		}
	}
	return result
}

func (r *Registry[P, C]) AddAssets(assets ...Asset) {
	r.assets = append(r.assets, assets...)
}

func (r *Registry[P, C]) AssetsForSurface(surface Surface) []Asset {
	result := []Asset{}
	for _, asset := range r.assets {
		if asset.Surface == surface {
			result = append(result, asset)
		}
	}
	return result
}

func (r *Registry[P, C]) AddEditorProviders(providers ...EditorProviderRegistration) {
	for _, provider := range providers {
		if strings.TrimSpace(provider.ID) == "" {
			continue
		}
		replaced := false
		for i := range r.editorProviders {
			if r.editorProviders[i].ID == provider.ID {
				r.editorProviders[i] = provider
				replaced = true
				break
			}
		}
		if !replaced {
			r.editorProviders = append(r.editorProviders, provider)
		}
	}
}

func (r *Registry[P, C]) EditorProviders() []EditorProviderRegistration {
	items := append([]EditorProviderRegistration(nil), r.editorProviders...)
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Priority == items[j].Priority {
			return items[i].Label < items[j].Label
		}
		return items[i].Priority < items[j].Priority
	})
	return items
}

func (r *Registry[P, C]) ResolveEditorProvider(id string) (EditorProviderRegistration, bool) {
	items := r.EditorProviders()
	for _, item := range items {
		if item.ID == id {
			return item, true
		}
	}
	if len(items) == 0 {
		return EditorProviderRegistration{}, false
	}
	return items[0], true
}
