package users

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pscheid92/dwarferl/internal"
)

type Service struct {
	repository internal.UsersRepository
}

func NewService(repository internal.UsersRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) CreateWithGoogleID(googleID string, email string) (internal.User, error) {
	user := internal.User{
		ID:       uuid.New().String(),
		Email:    email,
		GoogleID: googleID,
	}

	err := s.repository.Save(user)
	if err != nil {
		return internal.User{}, err
	}

	return user, nil
}

func (s *Service) GetOrCreateByGoogle(googleID string, email string) (internal.User, error) {
	user, err := s.repository.GetByGoogleID(googleID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return internal.User{}, err
	}

	// found it!
	if err == nil {
		return user, nil
	}

	return s.CreateWithGoogleID(googleID, email)
}
