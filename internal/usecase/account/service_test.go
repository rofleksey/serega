package account

import (
	"context"
	"errors"
	"testing"
)

type accountStore struct {
	users map[string]User
	err   error
}

func (s *accountStore) FindUserByUsername(_ context.Context, username string) (User, error) {
	if s.err != nil {
		return User{}, s.err
	}

	user, ok := s.users[username]
	if !ok {
		return User{}, ErrUserNotFound
	}

	return user, nil
}
func (s *accountStore) FindUserByID(context.Context, string) (User, error) {
	return User{}, ErrUserNotFound
}
func (s *accountStore) CreateUser(_ context.Context, username, hash string) (User, error) {
	if _, ok := s.users[username]; ok {
		return User{}, ErrUsernameTaken
	}

	user := User{ID: username, Username: username, PasswordHash: hash}
	s.users[username] = user

	return user, nil
}

func TestProvisioningKeepsEveryAccountAndHashesPasswords(t *testing.T) {
	store := &accountStore{users: map[string]User{}}
	service := NewService(store)

	const password = "correct horse battery"

	first, err := service.CreateUser(context.Background(), " alice ", password)
	if err != nil {
		t.Fatal(err)
	}

	if first.Username != "alice" || first.PasswordHash == password || !CheckPassword(first.PasswordHash, password) {
		t.Fatalf("invalid stored credential for %s", first.Username)
	}

	if _, err := service.CreateUser(context.Background(), "bob", password); err != nil {
		t.Fatal(err)
	}

	if len(store.users) != 2 {
		t.Fatal("provisioning removed another account")
	}

	if _, err := service.Authenticate(context.Background(), "alice", password); err != nil {
		t.Fatal(err)
	}

	if _, err := service.CreateUser(context.Background(), "alice", password); !errors.Is(err, ErrUsernameTaken) {
		t.Fatalf("duplicate error = %v", err)
	}
}

func TestAuthenticationSeparatesInvalidCredentialsFromStoreFailure(t *testing.T) {
	store := &accountStore{users: map[string]User{}}

	service := NewService(store)
	if _, err := service.Authenticate(context.Background(), "missing", "password"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("missing account error = %v", err)
	}

	store.err = errors.New("database unavailable")
	if _, err := service.Authenticate(context.Background(), "alice", "password"); !errors.Is(err, store.err) {
		t.Fatalf("store failure = %v", err)
	}
}

func TestPasswordValidationAndMalformedHashes(t *testing.T) {
	if _, err := HashPassword("short"); err == nil {
		t.Fatal("short password accepted")
	}

	for _, hash := range []string{"", "$argon2id$v=19$m=999999999,t=3,p=2$invalid$invalid", "$argon2id$v=19$m=65536,t=3,p=2$invalid$invalid"} {
		if CheckPassword(hash, "password") {
			t.Fatal("malformed hash accepted")
		}
	}
}
