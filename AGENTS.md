# Agent notes

Panel owns UI-neutral control-plane descriptors: resources, pages, widgets, actions, forms, tables, navigation, and policy.

It does not own content fields, persistence, HTTP delivery, or editor binding. Codex owns content fields. FormSet binds editor forms. GoBackend projects a Codex manifest into these descriptors.

## Rules

1. Do not add a database, REST API, renderer, or product schema.
2. Do not import Codex, FormSet, or Framework.
3. Descriptors stay generic over the application's principal and capability types.
4. Comments and documentation are written in English.

Run before completion:

```text
go test ./...
go vet ./...
gofmt -w .
```
