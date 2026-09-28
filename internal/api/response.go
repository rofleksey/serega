// Package api owns handwritten HTTP contract validation and response encoding.
package api

import (
	"encoding/json"
	"net/http"

	"github.com/rofleksey/serega/internal/api/generated"
)

// WriteError writes Serega's stable HTTP error envelope.
func WriteError(w http.ResponseWriter, status int, code, message, requestID string, fields ...generated.FieldError) {
	response := generated.ErrorResponse{}
	response.Error.Code = code
	response.Error.Message = message
	response.Error.RequestID = requestID

	if len(fields) > 0 {
		response.Error.Fields = &fields
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}
