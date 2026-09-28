package observability

import (
	"time"

	"github.com/rofleksey/serega/internal/entity"
)

func durationMilliseconds(duration time.Duration) float64 {
	return float64(duration) / float64(time.Millisecond)
}

// Keep Serega's bounded HTTP method vocabulary at the metric boundary.
func boundedHTTPMethod(method string) string {
	switch method {
	case "CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "PATCH", "POST", "PUT", "TRACE":
		return method
	default:
		return entity.MetricMethodOther
	}
}
