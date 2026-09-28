package board

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/rofleksey/serega/internal/entity"
)

const testCardID = "1056f797-2f30-4d21-b5c5-e40441cdbdaf"
const testTitle = "Task"

type recordingStore struct {
	calls          int
	input          CreateInput
	update         UpdateInput
	userID, cardID string
	version        int64
	err            error
}

func (s *recordingStore) ListCards(context.Context) ([]Card, error) { return []Card{}, s.err }
func (s *recordingStore) CreateCard(_ context.Context, userID string, input CreateInput) (Card, error) {
	s.calls++
	s.input = input
	s.userID = userID

	return Card{}, s.err
}
func (s *recordingStore) UpdateCard(_ context.Context, userID, cardID string, input UpdateInput) (Card, error) {
	s.calls++
	s.update = input
	s.userID = userID
	s.cardID = cardID

	return Card{}, s.err
}
func (s *recordingStore) DeleteCard(_ context.Context, cardID string, version int64) error {
	s.calls++
	s.cardID = cardID
	s.version = version

	return s.err
}

func TestCreateValidatesAndNormalizesBeforePersistence(t *testing.T) {
	for _, test := range []struct{ name, title, description, field string }{
		{name: "empty title", title: " \n ", field: entity.FieldTitle},
		{name: "long title", title: strings.Repeat("ж", 201), field: entity.FieldTitle},
		{name: "long description", title: testTitle, description: strings.Repeat("ж", 5001), field: entity.FieldDescription},
	} {
		t.Run(test.name, func(t *testing.T) {
			store := &recordingStore{}
			_, err := NewService(store).Create(context.Background(), "user", CreateInput{Title: test.title, Description: test.description})

			var invalid *ValidationError
			if !errors.As(err, &invalid) || invalid.Field != test.field {
				t.Fatalf("error = %v", err)
			}

			if store.calls != 0 {
				t.Fatal("invalid card reached persistence")
			}
		})
	}

	store := &recordingStore{}

	_, err := NewService(store).Create(context.Background(), "second-user", CreateInput{Title: "  " + strings.Repeat("ж", 200) + "  ", Description: "keep\n whitespace "})
	if err != nil {
		t.Fatal(err)
	}

	if store.input.Title != strings.Repeat("ж", 200) || store.input.Description != "keep\n whitespace " || store.userID != "second-user" {
		t.Fatalf("persisted input = %#v", store)
	}
}

func TestUpdateRequiresValidStatusIdentityAndVersion(t *testing.T) {
	for _, test := range []struct {
		name, id, status, field string
		version                 int64
	}{
		{name: "invalid status", id: testCardID, status: "archived", version: 1, field: entity.FieldStatus},
		{name: "missing version", id: testCardID, status: entity.CardStatusTodo, field: entity.FieldVersion},
		{name: "invalid id", id: "invalid", status: entity.CardStatusTodo, version: 1, field: entity.FieldCardID},
		{name: "noncanonical id", id: strings.ToUpper(testCardID), status: entity.CardStatusTodo, version: 1, field: entity.FieldCardID},
	} {
		t.Run(test.name, func(t *testing.T) {
			store := &recordingStore{}
			_, err := NewService(store).Update(context.Background(), "user", test.id, UpdateInput{Title: testTitle, Status: test.status, Version: test.version})

			var invalid *ValidationError
			if !errors.As(err, &invalid) || invalid.Field != test.field {
				t.Fatalf("error = %v", err)
			}

			if store.calls != 0 {
				t.Fatal("invalid update reached persistence")
			}
		})
	}
}

func TestSharedBoardPreservesActorVersionAndConflict(t *testing.T) {
	store := &recordingStore{err: ErrConflict}

	_, err := NewService(store).Update(context.Background(), "another-user", testCardID, UpdateInput{Title: " Task ", Status: entity.CardStatusDone, Version: 3})
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("error = %v", err)
	}

	if store.userID != "another-user" || store.cardID != testCardID || store.update.Title != testTitle || store.update.Version != 3 {
		t.Fatalf("update = %#v", store)
	}

	if err := NewService(store).Delete(context.Background(), testCardID, 3); !errors.Is(err, ErrConflict) {
		t.Fatalf("delete error = %v", err)
	}

	if store.version != 3 {
		t.Fatal("delete discarded optimistic version")
	}
}
