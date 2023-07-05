package service

import (
	"context"
	"fmt"

	"github.com/Vialmsi/Interview/internal/entity"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const (
	// PgUniqueEntry 23505 - unique_violation [from pg doc]
	PgUniqueEntry = "23505"
)

type ServiceActions interface {
	ProductStore
	UserStore
}

type Service struct {
	logger       *logrus.Logger
	productStore ProductStore
	userStore    UserStore
}

func NewService(logger *logrus.Logger, productStore ProductStore, userStore UserStore) *Service {
	return &Service{
		logger:       logger,
		productStore: productStore,
		userStore:    userStore,
	}
}

func (s *Service) RegisterUser(ctx context.Context, user entity.User) (int, error) {
	s.logger.Info("[RegisterUser] started")

	userID, err := s.userStore.NewUser(ctx, user)
	if err != nil {
		s.logger.Errorf("[RegisterUser] error in store: %s", err.Error())
		if isUniqueViolation(err) {
			return 0, fmt.Errorf("user with this login already exists")
		}
		return 0, fmt.Errorf("error while process request, try later\n%w", err)
	}

	err = s.productStore.CreateUserStorage(userID)
	if err != nil {
		s.logger.Errorf("[RegisterUser] error while creating user storage: %s", err.Error())
		return 0, fmt.Errorf("error while process request, try later\n%w", err)
	}

	s.logger.Info(userID)
	s.logger.Info("[RegisterUser] ended")

	return userID, nil
}

// isUniqueViolation func which check postgresql error for unique violation
func isUniqueViolation(err error) bool {
	pgErr, ok := err.(*pq.Error)
	return ok && pgErr.Code == PgUniqueEntry
}
