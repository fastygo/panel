# Changelog

## Unreleased

### Changed (breaking)
- Depends on `github.com/fastygo/codex` v0.4.2 and `github.com/fastygo/formset`
  (main, pseudo-version). See [ADR-0001](docs/adr/0001-codex-records-project-into-resources.md).
- `RecordType`, `RecordTypeID`, `ValidationRule` are Codex aliases; the
  panel-owned `RecordType` with Form/Table/Detail/Permissions is removed
  (`RecordPermission` stays).

### Added
- `ResourceFromRecord`, `ResourceOptions`, `FieldsFromRecord`,
  `TableFromRecord`, `WidgetFor`, `RecordValidator`.
- `Field.DataType`, `Field.Localized`, widget `FieldObject`.
- CI: gofmt, `go mod tidy -diff`, build (incl. examples), vet, `test -race`,
  golangci-lint. Coverage 67.3% → 98.1%; examples have smoke tests.
- `examples/cms-adapter` projects a Codex record.

### Initial

- Start the standalone `github.com/fastygo/panel` module.
- Add the first control-plane kernel for panels, registries, resources, pages, widgets, actions, schemas, policies, assets, editor providers, and optional data-source contracts.
- Add documentation and examples for CMS, CRM, operations, and chat-oriented panels.
