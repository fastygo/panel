# Concepts

## Panel

A `Panel` is a mounted control surface with its own path, title, navigation, assets, resources, pages, widgets, and policies.

## Resource

A `Resource` describes CRUD-oriented control of an application record type. It declares table, form, detail, relation, action, route, and capability metadata. The resource does not persist data itself.

## Page

A `Page` describes an arbitrary screen: dashboard, settings, report, setup wizard, runtime status, moderation queue, or custom workspace.

## Widget

A `Widget` describes reusable page blocks such as stat cards, charts, health checks, activity feeds, and table summaries.

## Action

An `Action` describes an operator-triggered UI operation. Actions can appear in headers, table rows, bulk selections, forms, or modals.

## Schema

Schemas describe forms, tables, detail views, fields, filters, columns, repeaters, relation selectors, validation metadata, and layout hints.

## Policy

Policies and capabilities decide whether a principal can see or execute a panel element.

## DataSource

Data sources are optional contracts for service-backed records and commands. Panel defines the shape; applications implement the behavior.
