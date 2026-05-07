# Filament Comparison

FastyGo Panel borrows the high-level control-plane vocabulary from Filament without copying Laravel or Livewire assumptions.

| Filament concept | FastyGo Panel concept |
| --- | --- |
| Panel | `Panel` |
| Resource | `Resource` |
| Resource pages | `ResourceRoute` plus app delivery adapters |
| Custom page | `Page` |
| Dashboard widget | `Widget` |
| Forms | `FormSchema`, `Field`, `SchemaSection` |
| Tables | `TableSchema`, `Column`, `Filter` |
| Actions | `Action` |
| Navigation | `NavigationGroup`, `MenuItem` |
| Authorization | `Principal`, `Policy`, capability fields |
| Assets | `Asset` |

## Differences

- Panel is framework-light and does not require Laravel, Livewire, or an ORM.
- Data access is application-owned.
- UI rendering is intentionally outside the first core package.
- Descriptors should be usable by server-rendered, API-driven, or generated UI adapters.
