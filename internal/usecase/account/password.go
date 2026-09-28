package account

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	passwordMinLength = 12
	argonMemory       = 64 * 1024
	argonIterations   = 3
	argonParallelism  = 2
	argonSaltLength   = 16
	argonKeyLength    = 32
)

func HashPassword(password string) (string, error) {
	if len(password) < passwordMinLength {
		return "", fmt.Errorf("password must be at least %d characters", passwordMinLength)
	}

	salt := make([]byte, argonSaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate password salt: %w", err)
	}

	hash := argon2.IDKey([]byte(password), salt, argonIterations, argonMemory, argonParallelism, argonKeyLength)

	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s", argonMemory, argonIterations, argonParallelism,
		base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(hash)), nil
}

func CheckPassword(encoded, password string) bool {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" || parts[2] != "v=19" {
		return false
	}

	params := map[string]uint32{}

	for _, raw := range strings.Split(parts[3], ",") {
		keyValue := strings.SplitN(raw, "=", 2)
		if len(keyValue) != 2 {
			return false
		}

		value, err := strconv.ParseUint(keyValue[1], 10, 32)
		if err != nil {
			return false
		}

		params[keyValue[0]] = uint32(value)
	}

	if params["m"] != argonMemory || params["t"] != argonIterations || params["p"] != argonParallelism {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil || len(salt) != argonSaltLength {
		return false
	}

	expected, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil || len(expected) != argonKeyLength {
		return false
	}

	actual := argon2.IDKey([]byte(password), salt, params["t"], params["m"], uint8(params["p"]), uint32(len(expected))) //nolint:gosec // Parsed Argon2 parameters are bounds-validated above.

	return subtle.ConstantTimeCompare(actual, expected) == 1
}
