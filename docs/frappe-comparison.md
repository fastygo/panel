# Frappe Comparison

Frappe shows how far an admin platform can grow: DocTypes describe data, views, permissions, workflows, reports, print formats, and APIs.

FastyGo Panel keeps that ambition but separates control plane from data plane.

| Frappe concept | FastyGo Panel direction |
| --- | --- |
| DocType | `RecordType` or application `Resource` |
| DocField | `Field` |
| List view | `TableSchema` |
| Form view | `FormSchema` |
| Report | `Report` |
| Workflow | `Workflow` |
| Timeline/comments/audit | `Timeline`, `TimelineEvent` |
| Print format | `PrintDescriptor` |
| Import/export | `ImportExportDescriptor` |
| Permissions | `Policy`, capability descriptors |

## Design Rule

Frappe-like contracts are optional. A simple admin panel should not be forced into a full low-code platform model. Applications can opt in to record metadata, workflows, reports, imports, exports, and print descriptors as they need them.
