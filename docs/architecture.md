# Architecture

FastyGo Panel is a control-plane package. It describes how operators interact with application data, but it does not own that data.

## Layers

```text
github.com/fastygo/framework
  HTTP, middleware, rendering, sessions, application composition

github.com/fastygo/panel
  panels, resources, pages, widgets, actions, schemas, navigation, policies

application packages
  CMS content, CRM leads, operations metrics, chat conversations, storage, APIs
```

## Control Plane

The control plane declares:

- where a panel is mounted;
- which resources, pages, and widgets are visible;
- which actions can be triggered;
- which schemas describe forms, tables, detail views, filters, and relations;
- which capabilities or policies gate access.

## Data Plane

The data plane implements:

- queries and commands;
- validation and business rules;
- persistence and external APIs;
- workflows, audits, timelines, metrics, and side effects.

Panel can expose optional `DataSource` contracts, but implementations must live in application modules.

## Boundary Rule

The core package must not import `github.com/fastygo/cms` or any application domain package. GoCMS, CRM, operations, and chat packages import panel, not the other way around.
