# Domain presentation and session identity

Entities may depend on entities/shared, not features/widgets/pages/app.

- Own reusable domain presentation and metadata, not route-level workflows.
- Keep domain display labels separate from transport values.
- Read transport shapes from shared/api instead of handwritten copies.
- Session identity/context belongs here so lower-level consumers do not import
  application composition.
- Keep local styles consistent with MUI/theme tokens and accessible semantics.
