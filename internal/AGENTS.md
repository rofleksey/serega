# Go application ownership

Follow the existing boundaries; subtree guides provide concrete rules.

| Package | Responsibility |
| --- | --- |
| app | Concrete dependency construction, runtime, cleanup |
| cli | Cobra parsing, standard input/output, service invocation |
| config | Explicit environment parsing and validation |
| entity | Shared domain data and stable vocabulary |
| api | Handwritten OpenAPI generation/validation/response boundary |
| handler | HTTP-to-service conversion and response mapping |
| middleware | Reusable HTTP mechanics |
| usecase | Policy, orchestration, consumer-owned ports |
| store/postgres | SQLC calls, transactions, row conversion |
| database/postgres | Explicit Goose lifecycle |
| frontend | Immutable embedded SPA serving |
| observability | Library bridges and domain metrics |

Handlers and commands never reach directly into persistence. Use cases never
import generated HTTP/SQL models or concrete adapters. app injects adapters
into the consumer-defined ports.

Keep shared contracts in types.go and private helpers with their behavior.
Avoid generic CRUD engines, global service locators, circular convenience
imports, and new layers with no owned responsibility.
