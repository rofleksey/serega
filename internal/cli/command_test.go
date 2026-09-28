package cli

import (
	"strings"
	"testing"
)

func TestPasswordStdinIsBoundedAndPreservesSpaces(t *testing.T) {
	password, err := provisionPassword(strings.NewReader("  twelve words  \r\n"), true)
	if err != nil || password != "  twelve words  " {
		t.Fatalf("password input changed: %v", err)
	}

	if _, err := provisionPassword(strings.NewReader(strings.Repeat("x", 1025)), true); err == nil {
		t.Error("accepted oversized password")
	}
}

func TestCLIRequiresCommand(t *testing.T) {
	if err := Execute(t.Context(), nil); err == nil {
		t.Error("missing command accepted")
	}

	if err := Execute(t.Context(), []string{"user", "create"}); err == nil {
		t.Error("missing username accepted")
	}
}
