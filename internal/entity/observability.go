package entity

const (
	EventHTTPRequest = "http.request.completed"
	EventMetricsInit = "observability.metrics.initialization"
)

const (
	OutcomeSucceeded = "succeeded"
	OutcomeFailed    = "failed"
	OutcomeRejected  = "rejected"
	OutcomeCanceled  = "canceled"
)

const (
	AuthMechanismSession      = "session"
	DecisionAllowed           = "allowed"
	DecisionDenied            = "denied"
	DecisionAccepted          = "accepted"
	DecisionRejected          = "rejected"
	DecisionRouteAvailable    = "route_available"
	DecisionRouteUnavailable  = "route_unavailable"
	DecisionTargetNotFound    = "target_not_found"
	DecisionActionRejected    = "action_rejected"
	DecisionValidationFailed  = "validation_failed"
	DecisionPersistenceFailed = "persistence_failed"
	APIErrorNotFound          = "not_found"
)

const (
	FieldEvent            = "event.name"
	FieldOutcome          = "outcome"
	FieldDurationMS       = "duration_ms"
	FieldErrorKind        = "error.kind"
	FieldRequestID        = "request_id"
	FieldProvider         = "provider"
	FieldOperation        = "operation"
	FieldAction           = "action"
	FieldHTTPMethod       = "http.method"
	FieldHTTPRoute        = "http.route"
	FieldHTTPStatusCode   = "http.status"
	FieldHTTPSize         = "http.response.body.size"
	FieldAPIErrorCode     = "api.error.code"
	FieldHandler          = "handler.name"
	FieldValidation       = "handler.validation"
	FieldAuthMechanism    = "auth.mechanism"
	FieldAuthDecision     = "auth.decision"
	FieldAccountID        = "account_id"
	FieldUseCase          = "usecase.name"
	FieldDecision         = "decision"
	FieldResultCount      = "result.count"
	FieldAdapter          = "adapter.name"
	FieldAdapterOperation = "adapter.operation"
	FieldRowsAffected     = "db.rows_affected"
)

const (
	MetricHTTPRequests     = "serega.http.server.requests"
	MetricHTTPDuration     = "serega.http.server.request.duration"
	MetricHTTPActive       = "serega.http.server.active_requests"
	MetricDBConnections    = "serega.db.client.connections"
	MetricDBMaxConnections = "serega.db.client.max_connections"
)

const (
	MetricAttrMethod     = "http.request.method"
	MetricAttrRoute      = "http.route"
	MetricAttrStatusCode = "http.response.status_code"
	MetricAttrOutcome    = "serega.outcome"
	MetricAttrProvider   = "serega.provider"
	MetricAttrOperation  = "serega.operation"
	MetricAttrState      = "state"
)

const (
	MetricStateAcquired = "acquired"
	MetricStateIdle     = "idle"
	MetricMethodOther   = "OTHER"
)

const (
	OperationAuthenticate   = "authenticate"
	OperationChangePassword = "change_password"
	AdapterPostgres         = "postgres"
)
