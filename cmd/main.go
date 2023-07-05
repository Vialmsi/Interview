package main

import (
	"errors"
	"os"

	"github.com/Vialmsi/Interview/internal/config"
	"github.com/Vialmsi/Interview/internal/store"

	"github.com/sirupsen/logrus"
)

const (
	FileStorage   = "files"
	DirPermission = 0777
)

func main() {
	logger := logrus.New()

	_, err := os.Open(FileStorage)
	if errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(FileStorage, DirPermission)
		if err != nil {
			logger.Fatalf("Error while creating file storage: %s", err.Error())
		}
	}

	cfg, err := config.Init()
	if err != nil {
		logger.Fatalf("Error while load config: %s", err.Error())
		return
	}

	userStore, err := store.NewUserStore(cfg.PSQLDatabase)
	if err != nil {
		logger.Fatalf("Error while init store: %s", err)
	}

	_ = userStore
}
