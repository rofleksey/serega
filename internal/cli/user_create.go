package cli

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/rofleksey/serega/internal/app"
	"github.com/spf13/cobra"
)

func newUserCreateCommand() *cobra.Command {
	var (
		username      string
		passwordStdin bool
	)

	command := &cobra.Command{
		Use: "create", Short: "Create a user without changing existing accounts", Args: cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			databaseURL := os.Getenv("DATABASE_URL")
			if databaseURL == "" {
				return errors.New("DATABASE_URL is required")
			}

			password, err := provisionPassword(command.InOrStdin(), passwordStdin)
			if err != nil {
				return err
			}

			return app.CreateUser(command.Context(), databaseURL, username, password)
		},
	}
	command.Flags().StringVar(&username, "username", "", "unique username")
	command.Flags().BoolVar(&passwordStdin, "password-stdin", false, "read password from standard input instead of a hidden terminal prompt")
	_ = command.MarkFlagRequired("username")

	return command
}

func provisionPassword(input io.Reader, passwordStdin bool) (string, error) {
	if passwordStdin {
		data, err := io.ReadAll(io.LimitReader(input, 1025))
		if err != nil {
			return "", err
		}

		if len(data) > 1024 {
			return "", errors.New("password input exceeds 1024 bytes")
		}

		return strings.TrimSuffix(strings.TrimSuffix(string(data), "\n"), "\r"), nil
	}

	password, err := readPassword("New password: ")
	if err != nil {
		return "", err
	}

	confirmation, err := readPassword("Confirm password: ")
	if err != nil {
		return "", err
	}

	if password != confirmation {
		return "", errors.New("passwords do not match")
	}

	return password, nil
}
