package cli

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/term"
)

// readPassword reads directly from the terminal with echo disabled. It neither
// accepts an environment value nor returns the password to any log path.
func readPassword(prompt string) (string, error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return "", errors.New("password input requires an interactive terminal")
	}

	fmt.Fprint(os.Stderr, prompt)

	value, err := term.ReadPassword(int(os.Stdin.Fd()))

	fmt.Fprintln(os.Stderr)

	if err != nil {
		return "", fmt.Errorf("read password: %w", err)
	}

	return string(value), nil
}
