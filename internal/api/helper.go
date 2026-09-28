package api

import "github.com/getkin/kin-openapi/openapi3"

func stripRequestParameters(document *openapi3.T) {
	if document.Paths == nil {
		return
	}

	for _, path := range document.Paths.Map() {
		path.Parameters = nil
		for _, operation := range path.Operations() {
			operation.Parameters = nil
		}
	}
}
