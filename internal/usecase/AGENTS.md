# Application services and ports

Keep small cohesive account/board packages following Akio's service pattern.

- Declare narrow consumer-owned ports and shared contracts in types.go.
- Inject dependencies through constructors; use entity/standard Go values.
- Do not import net/http, generated API/SQLC, pgx, or the concrete store.
- Own input validation/normalization even when OpenAPI/Zod validate transport/UI.
- Expose stable domain errors for decisions handlers must distinguish.
- Preserve diagnostic cause and request context for dependency failures.
- Add safe use-case/decision fields to the active completion event; do not log
  user input.
- Avoid generic use-case engines and layers with no responsibility, while
  preserving the transport-to-service-to-persistence boundary.
- Use small fake ports for meaningful validation/decision/failure tests,
  including proving rejected input never reaches persistence.
