package tokensvc

import (
	"io"
	"net/http"
	"strconv"

	"github.com/Vialmsi/Interview/internal/config"
	customErrors "github.com/Vialmsi/Interview/internal/errors"

	"github.com/sirupsen/logrus"
)

const (
	pingPath     = "/ping"
	generatePath = "/generate?login="
	validatePath = "/validate"
)

type TokenService struct {
	logger  *logrus.Logger
	address string
}

func NewTokenService(logger *logrus.Logger, cfg config.TokenServiceConfig) *TokenService {
	return &TokenService{
		logger:  logger,
		address: cfg.Address,
	}
}

func (t *TokenService) Ping() error {
	t.logger.Info("[Ping] started")

	resp, err := http.Get(t.address + pingPath)
	if err != nil {
		t.logger.Errorf("[Ping] Error while request token service: %s", err.Error())
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.logger.Errorf("[Ping] Error while parse request: %s", err.Error())
		return err
	}

	t.logger.Info(string(body))
	t.logger.Info("[Ping] ended")

	return nil
}

func (t *TokenService) Generate(userID int) (string, error) {
	t.logger.Info("[Generate] started")

	strID := strconv.Itoa(userID)

	resp, err := http.Get(t.address + generatePath + strID)
	if err != nil {
		t.logger.Errorf("[Generate] Error while request token service: %s", err.Error())
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.logger.Errorf("[Generate] Error while parse response: %s", err.Error())
		return "", err
	}

	t.logger.Info(string(body))
	t.logger.Info("[Generate] ended")

	return string(body), nil
}

func (t *TokenService) Validate(token string) error {
	t.logger.Info("[Validate] started")

	req, err := http.NewRequest("GET", t.address+validatePath, nil)
	if err != nil {
		t.logger.Errorf("[Validate] Failed to create request: %s", err.Error())
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.logger.Errorf("[Validate] Error while parse response: %s", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.logger.Errorf("[Validate] Error while parse response: %s", err.Error())
		return err
	}

	if string(body) != "" {
		t.logger.Errorf("[Validate] Error from token service: %s", body)
		return customErrors.TokenExpiredError
	}

	t.logger.Info("[Validate] ended")

	return nil
}
