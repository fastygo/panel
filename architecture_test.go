package panel

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPanelDoesNotImportGoCMS(t *testing.T) {
	disallowed := "github.com/fastygo/" + "cms"
	err := filepath.WalkDir(".", func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			switch entry.Name() {
			case ".git":
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if strings.Contains(string(content), disallowed) {
			t.Fatalf("%s imports or references disallowed CMS module", path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("WalkDir() error = %v", err)
	}
}
