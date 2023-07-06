package service

import (
	"errors"

	"github.com/sirupsen/logrus"
)

var (
	UniqueViolationError = errors.New("user with this login already exists")
)

type ServiceActions interface {
	ProductStore
	UserStore
}

type Service struct {
	logger *logrus.Logger
	store  ServiceActions
}

func NewService(logger *logrus.Logger, store ServiceActions) *Service {
	return &Service{
		logger: logger,
		store:  store,
	}
}
