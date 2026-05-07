package panel

import "testing"

func TestPanelKeepsResourcesPagesAndWidgetsSeparate(t *testing.T) {
	admin, err := NewPanel[testPrincipal, testCapability](PanelOptions[testCapability]{
		ID:       "admin",
		Title:    "Admin",
		BasePath: "/admin",
	})
	if err != nil {
		t.Fatalf("NewPanel() error = %v", err)
	}

	if err := admin.AddResources(Resource[testCapability]{
		ID:       "leads",
		Label:    "Leads",
		BasePath: "/admin/leads",
		Navigation: MenuItem[testCapability]{
			ID:    "leads",
			Label: "Leads",
			Path:  "/admin/leads",
		},
	}); err != nil {
		t.Fatalf("AddResources() error = %v", err)
	}
	if err := admin.AddPages(Page[testCapability]{
		ID:    "settings",
		Title: "Settings",
		Path:  "/admin/settings",
	}); err != nil {
		t.Fatalf("AddPages() error = %v", err)
	}
	if err := admin.AddWidgets(Widget[testCapability]{
		ID:    "pipeline",
		Title: "Pipeline",
	}); err != nil {
		t.Fatalf("AddWidgets() error = %v", err)
	}

	if len(admin.Resources()) != 1 || len(admin.Pages()) != 1 || len(admin.Widgets()) != 1 {
		t.Fatalf("panel counts = resources:%d pages:%d widgets:%d", len(admin.Resources()), len(admin.Pages()), len(admin.Widgets()))
	}
	if got := admin.Registry().MenuItems(newTestPrincipal()); len(got) != 1 || got[0].ID != "leads" {
		t.Fatalf("registered menu items = %v, want leads", got)
	}
}

func TestPanelValidatesRequiredFields(t *testing.T) {
	if _, err := NewPanel[testPrincipal, testCapability](PanelOptions[testCapability]{}); err == nil {
		t.Fatalf("expected invalid empty panel options")
	}
	admin, err := NewPanel[testPrincipal, testCapability](PanelOptions[testCapability]{
		ID:       "admin",
		Title:    "Admin",
		BasePath: "/admin",
	})
	if err != nil {
		t.Fatalf("NewPanel() error = %v", err)
	}
	if err := admin.AddResources(Resource[testCapability]{}); err == nil {
		t.Fatalf("expected invalid empty resource")
	}
	if err := admin.AddPages(Page[testCapability]{}); err == nil {
		t.Fatalf("expected invalid empty page")
	}
	if err := admin.AddWidgets(Widget[testCapability]{}); err == nil {
		t.Fatalf("expected invalid empty widget")
	}
}
