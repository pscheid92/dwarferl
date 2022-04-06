package users

import (
	"errors"
	"github.com/jackc/pgx/v4"
	"github.com/pscheid92/dwarferl/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	testUser     = "00000000-0000-0000-0000-000000000000"
	testGoogleID = "testGoogleID"
	testEmail    = "example@example.com"
)

func TestService_CreateWithGoogleID(t *testing.T) {
	repo, sut := setupService()

	user, err := sut.CreateWithGoogleID(testGoogleID, testEmail)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)
	assert.NotEmptyf(t, user.ID, "Expected user ID to be set, got %v", user.ID)
	assert.Equalf(t, testGoogleID, user.GoogleID, "Expected GoogleID to be %v, got %v", testGoogleID, user.GoogleID)

	repo.FailMode = true
	_, err = sut.CreateWithGoogleID(testGoogleID, testEmail)
	assert.Errorf(t, err, "Expected error, got nil")
}

func TestService_GetOrCreateByGoogle(t *testing.T) {
	repo, sut := setupService()

	user, err := sut.GetOrCreateByGoogle(testGoogleID, testEmail)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)
	assert.Equalf(t, testUser, user.ID, "Expected user ID to be %v, got %v", testUser, user.ID)

	user, err = sut.GetOrCreateByGoogle("nonexistent", testEmail)
	assert.NoErrorf(t, err, "Expected no error, got %v", err)
	assert.NotEmptyf(t, user.ID, "Expected user ID to be set, got %v", user.ID)
	assert.Equalf(t, testEmail, user.Email, "Expected user email to be %v, got %v", testEmail, user.Email)
	assert.Equalf(t, "nonexistent", user.GoogleID, "Expected GoogleID to be %v, got %v", "nonexistent", user.GoogleID)

	repo.FailMode = true
	_, err = sut.GetOrCreateByGoogle(testGoogleID, testEmail)
	assert.Errorf(t, err, "Expected error, got nil")
	assert.NotErrorIsf(t, err, pgx.ErrNoRows, "Expected error to not be %v, got %v", pgx.ErrNoRows, err)
}

func setupService() (*usersRepositoryFake, *Service) {
	repo := &usersRepositoryFake{}
	svc := NewService(repo)
	return repo, svc
}

type usersRepositoryFake struct {
	FailMode bool
}

func (u usersRepositoryFake) Save(_ internal.User) error {
	if u.FailMode {
		return errors.New("fake error")
	}
	return nil
}

func (u usersRepositoryFake) GetByGoogleID(googleID string) (internal.User, error) {
	if u.FailMode {
		return internal.User{}, errors.New("fake error")
	}

	if googleID != testGoogleID {
		return internal.User{}, pgx.ErrNoRows
	}

	return internal.User{
		ID:       testUser,
		Email:    testEmail,
		GoogleID: testGoogleID,
	}, nil
}
