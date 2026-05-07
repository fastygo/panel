# GoCMS Integration

GoCMS should consume `github.com/fastygo/panel` through a CMS-specific adapter package such as `cmspanel`.

## What belongs in panel

- Generic panel registry.
- Resource, page, widget, action, form, table, and policy descriptors.
- Editor provider registration.
- Optional data-source interfaces.

## What belongs in GoCMS

- Content services.
- Media services.
- Taxonomy, menu, theme, permalink, and settings services.
- CMS capabilities and roles.
- CMS plugin runtime.
- Admin delivery handlers and templates.

## Migration Shape

`cmspanel.PostsResource` and `cmspanel.PagesResource` should import `github.com/fastygo/panel` for descriptors, while continuing to use GoCMS application services for real reads and writes.
