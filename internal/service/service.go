package service

import (
	"github.com/sirupsen/logrus"
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
