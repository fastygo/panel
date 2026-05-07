# Contributing

Thank you for contributing to FastyGo Panel.

## Principles

- Keep `github.com/fastygo/panel` generic. Do not import CMS, CRM, operations, chat, or storage-specific packages into the core package.
- Treat panel as a control plane. Data-plane implementations belong in application modules.
- Prefer typed descriptors and small interfaces over concrete UI or database assumptions.
- Write comments and documentation in English.
- Add tests for every new public contract.

## Checks

Run before submitting changes:

```bash
go test ./...
go vet ./...
```

When adapting GoCMS to a local panel checkout, also run the GoCMS verification suite.
