package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Vialmsi/Interview/internal/entity"
	customErrors "github.com/Vialmsi/Interview/internal/errors"
)

type UserStore interface {
	NewUser(ctx context.Context, user entity.User) (int, error)
	RetrieveUser(ctx context.Context, login, password string) (entity.User, error)
}

func (s *Service) LoginUser(ctx context.Context, user entity.User) (*entity.User, error) {
	s.logger.Info("[RetrieveUser] started")

	user, err := s.store.RetrieveUser(ctx, user.Login, user.Password)
	if err != nil {
		s.logger.Errorf("[RetrieveUser] error in store: %s", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customErrors.IncorrectLoginOrPasswordError
		}
		return nil, fmt.Errorf("error while process request\n%w", err)
	}

	s.logger.Info(user)
	s.logger.Info("[RetrieveUser] ended")

	return &user, nil
}

func (s *Service) RegisterUser(ctx context.Context, user entity.User) (int, error) {
	s.logger.Info("[RegisterUser] started")

	userID, err := s.store.NewUser(ctx, user)
	if err != nil {
		s.logger.Errorf("[RegisterUser] error in store: %s", err.Error())
		if errors.Is(err, customErrors.UserAlreadyExistError) {
			return 0, fmt.Errorf("user with this login already exists")
		}
		return 0, fmt.Errorf("error while process request, try later\n%w", err)
	}

	s.logger.Info(userID)
	s.logger.Info("[RegisterUser] ended")

	return userID, nil
}
