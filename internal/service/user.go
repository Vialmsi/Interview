package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Vialmsi/Interview/internal/entity"
)

type UserStore interface {
	NewUser(ctx context.Context, user entity.User) (int, error)
	RetrieveUser(ctx context.Context, login, password string) (entity.User, error)
}

func (s *Service) LoginUser(ctx context.Context, user entity.User) (*entity.User, error) {
	s.logger.Info("[RetrieveUser] started")

	user, err := s.userStore.RetrieveUser(ctx, user.Login, user.Password)
	if err != nil {
		s.logger.Errorf("[RetrieveUser] error in store: %s", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("incorrect login or password")
		}
		return nil, fmt.Errorf("error while process request\n%w", err)
	}

	s.logger.Info(user)
	s.logger.Info("[RetrieveUser] ended")

	return &user, nil
}
