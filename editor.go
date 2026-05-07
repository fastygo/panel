package panel

type EditorKind string

const (
	EditorRichText EditorKind = "richtext"
	EditorMarkdown EditorKind = "markdown"
	EditorCode     EditorKind = "code"
	EditorJSON     EditorKind = "json"
	EditorPlain    EditorKind = "plain"
)

type EditorProvider struct {
	ID          string
	Label       string
	Kind        EditorKind
	Description string
	AssetIDs    []string
	Priority    int
	Formats     []string
}
