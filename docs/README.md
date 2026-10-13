# Domain policy documents

`docs/` owns the product's mechanism-level policy: the decisions a contributor
needs that do not belong in the repository-wide rules ([AGENTS.md](../AGENTS.md))
or the operator instructions ([README.md](../README.md)).

- One document per mechanism or domain area with a descriptive kebab-case name
  (`docs/board.md` in this template).
- Record decisions and measured constraints, not tutorials: bounds, statuses,
  preconditions, concurrency semantics, retry behavior, and telemetry vocabulary.
- Link the document from README or the owning guide and keep the index below
  current.
- Keep the guidance hierarchy intact: AGENTS.md owns repository rules, README
  owns setup/commands/configuration, scoped guides own subtree constraints. A
  policy document explains what the product decided; it does not add repository
  law.
- When the template becomes a new product, replace the example policy documents
  with the new product's policies instead of carrying unused vocabulary.

## Index

| Document | Owns |
| --- | --- |
| [board.md](board.md) | Example domain policy for the shared board. |

## Growth paths

The template stays small on purpose. When a product needs a mechanism the
template does not demonstrate, add it as a complete vertical slice following the
rules in AGENTS.md. The starting shapes are:

- Outbound adapters (external services, filesystems, process execution): declare
  the consumer-owned port in the use case and implement it in a
  technology-named adapter package. Handlers and commands never reach the
  adapter; `app` wires it and tests substitute a fake.
- Asynchronous work: model durable work as explicit jobs with claimed leases
  rather than goroutine-only background tasks. Persist first, publish
  conservatively, and give each resource one lifecycle owner that shutdown joins.
- Multi-deployable products: keep one handwritten protocol package with
  canonical byte fixtures shared between the deployables, and compare fixtures
  in CI. Deployables exchange wire data; one never imports the other.
- Scale coverage: keep load or soak scenarios opt-in behind explicit environment
  gates with synthetic fixtures, separate from deterministic integration tests.

These are patterns, not platforms: add a mechanism only when the product owns
the responsibility.
