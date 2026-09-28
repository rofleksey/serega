package api

import "testing"

func walkContract(value any, visitReference func(string)) {
	switch value := value.(type) {
	case map[string]any:
		if reference, ok := value["$ref"].(string); ok {
			visitReference(reference)
		}

		for _, child := range value {
			walkContract(child, visitReference)
		}
	case []any:
		for _, child := range value {
			walkContract(child, visitReference)
		}
	}
}

func contractObject(t *testing.T, value any, name string) map[string]any {
	t.Helper()

	object, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("%s is not an object", name)
	}

	return object
}
