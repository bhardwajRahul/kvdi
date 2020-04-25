package rethinkdb

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/tinyzimmer/kvdi/pkg/auth/types"
	"github.com/tinyzimmer/kvdi/pkg/util/errors"
)

type fakeClock struct{}

func (*fakeClock) Now() time.Time { return time.Unix(0, 0) }

func newFakeUUID() uuid.UUID { return uuid.UUID{} }

func TestGetUserSession(t *testing.T) {
	mock := NewMock()
	if _, err := mock.GetUserSession(newItem); err != nil {
		t.Error("Expected no error, got:", err)
	}
	if _, err := mock.GetUserSession(nonExist); err == nil {
		t.Error("Expected error, got nil")
	} else if !errors.IsUserSessionNotFoundError(err) {
		t.Error("Expected session not found error, got", err)
	}
	if _, err := mock.GetUserSession(errItem); err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestCreateUserSession(t *testing.T) {
	mock := NewMock()
	fk := &fakeClock{}
	nowFunc = fk.Now
	uuidFunc = newFakeUUID
	if _, err := mock.CreateUserSession(&types.User{Name: newItem}); err != nil {
		t.Error("Expected no error got,", err)
	}
	if _, err := mock.CreateUserSession(&types.User{Name: errItem}); err == nil {
		t.Error("Expected error got nil")
	}
}

func TestDeleteUserSession(t *testing.T) {
	mock := NewMock()
	if err := mock.DeleteUserSession(&types.UserSession{Token: testToken}); err != nil {
		t.Error("Expected no error got,", err)
	}
	if err := mock.DeleteUserSession(&types.UserSession{Token: errItem}); err == nil {
		t.Error("Expected error got nil")
	}
}