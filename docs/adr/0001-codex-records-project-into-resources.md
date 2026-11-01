# ADR-0001: Codex records project into panel resources

Status: accepted (2026-11-01). Supersedes the "zero-dependency panel" stance.

## Decision

- panel depends on `github.com/fastygo/codex` (record vocabulary) and
  `github.com/fastygo/formset` (binding/validation). Codex is the single
  source of truth for data structures.
- `panel.RecordType`, `RecordTypeID` and `ValidationRule` are aliases of
  `codex/schema`. The previous panel-owned `RecordType` (with its own
  Fields/Form/Table) is removed.
- `panel.FieldType` remains a **UI widget** vocabulary (repeater, file,
  hidden, status, ...). The data type lives in `Field.DataType`
  (`schema.FieldType`) together with Codex `Rules` and `Localized`. Widget
  names shared with Codex use the same strings.
- `ResourceFromRecord` is the default projection: form/detail fields from
  Codex fields, table columns from `Indexed`/`Searchable` (never
  `Sensitive`/`Encrypted`), filters from indexed select/boolean/datetime,
  relations from Codex relations (`Multiple` = cardinality != one-to-one),
  media via `UIHint: media` → file widget / image column.
- `RecordValidator` implements `Validator` with formset + Codex rules.

## Consequences

Backends stop hand-writing columns. Hand-written resources keep working;
projection output can be edited before registering.
