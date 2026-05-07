package panel

import "testing"

type testCapability string

type testPrincipal struct {
	capabilities map[testCapability]struct{}
}

func newTestPrincipal(capabilities ...testCapability) testPrincipal {
	principal := testPrincipal{capabilities: map[testCapability]struct{}{}}
	for _, capability := range capabilities {
		principal.capabilities[capability] = struct{}{}
	}
	return principal
}

func (p testPrincipal) Has(capability testCapability) bool {
	_, ok := p.capabilities[capability]
	return ok
}

func TestRegistryFiltersAndSortsMenuItems(t *testing.T) {
	registry := NewRegistry[testPrincipal, testCapability]()
	registry.AddMenuItems(
		MenuItem[testCapability]{ID: "settings", Label: "Settings", Path: "/settings", Order: 20, Capability: "settings.manage"},
		MenuItem[testCapability]{ID: "dashboard", Label: "Dashboard", Path: "/", Order: 0},
		MenuItem[testCapability]{ID: "posts", Label: "Posts", Path: "/posts", Order: 10, Capability: "content.read"},
	)

	items := registry.MenuItems(newTestPrincipal("content.read"))
	if len(items) != 2 {
		t.Fatalf("menu items = %v, want dashboard and posts", items)
	}
	if items[0].ID != "dashboard" || items[1].ID != "posts" {
		t.Fatalf("menu order = %v, want dashboard then posts", items)
	}

	navItems := registry.NavItems(newTestPrincipal("settings.manage"))
	if len(navItems) != 2 || navItems[1].Path != "/settings" {
		t.Fatalf("nav items = %v, want dashboard and settings", navItems)
	}
}

func TestRegistryFiltersRoutesAndAssetsBySurface(t *testing.T) {
	registry := NewRegistry[testPrincipal, testCapability]()
	registry.AddRoutes(
		Route[testPrincipal, testCapability]{Pattern: "GET /admin", Surface: SurfaceAdmin},
		Route[testPrincipal, testCapability]{Pattern: "GET /api", Surface: SurfaceREST},
	)
	registry.AddAssets(
		Asset{ID: "admin-js", Surface: SurfaceAdmin, Path: "/static/admin.js"},
		Asset{ID: "public-js", Surface: SurfacePublic, Path: "/static/public.js"},
	)

	adminRoutes := registry.RoutesForSurface(SurfaceAdmin)
	if len(adminRoutes) != 1 || adminRoutes[0].Pattern != "GET /admin" {
		t.Fatalf("admin routes = %v, want GET /admin", adminRoutes)
	}
	publicAssets := registry.AssetsForSurface(SurfacePublic)
	if len(publicAssets) != 1 || publicAssets[0].ID != "public-js" {
		t.Fatalf("public assets = %v, want public-js", publicAssets)
	}
}

func TestRegistryResolveEditorProviderUsesPriorityFallback(t *testing.T) {
	registry := NewRegistry[testPrincipal, testCapability]()
	registry.AddEditorProviders(
		EditorProviderRegistration{ID: "secondary", Label: "Secondary", Priority: 20},
		EditorProviderRegistration{ID: "primary", Label: "Primary", Priority: 10},
	)

	provider, ok := registry.ResolveEditorProvider("missing")
	if !ok {
		t.Fatalf("expected fallback editor provider")
	}
	if provider.ID != "primary" {
		t.Fatalf("provider = %q, want primary", provider.ID)
	}

	registry.AddEditorProviders(EditorProviderRegistration{ID: "primary", Label: "Primary Updated", Priority: 30})
	provider, ok = registry.ResolveEditorProvider("primary")
	if !ok || provider.Label != "Primary Updated" {
		t.Fatalf("provider = %v, %v, want updated primary", provider, ok)
	}
}
