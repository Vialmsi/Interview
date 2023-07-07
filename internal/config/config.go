package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	PSQLDatabaseDriver   = "PSQL_DATABASE_DRIVER"
	PSQLDatabaseUser     = "PSQL_DATABASE_USER"
	PSQLDatabasePassword = "PSQL_DATABASE_PASSWORD"
	PSQLDatabaseHost     = "PSQL_DATABASE_HOST"
	PSQLDatabasePort     = "PSQL_DATABASE_PORT"
	PSQLDatabaseName     = "PSQL_DATABASE_NAME"
	PSQLDatabaseTimeout  = "PSQL_DATABASE_TIMEOUT"

	DefaultQueryLimit = "DEFAULT_QUERY_LIMIT"
	DefaultQueryPage  = "DEFAULT_QUERY_PAGE"

	TokenServiceProtocol = "TOKEN_SERVICE_PROTOCOL"
	TokenServiceHost     = "TOKEN_SERVICE_HOST"
	TokenServicePort     = "TOKEN_SERVICE_PORT"

	ServerHost = "SERVER_HOST"
	ServerPort = "SERVER_PORT"

	JWTSigningKey = "JWT_SIGNING_KEY"
)

type Config struct {
	PSQLDatabase       PSQLDatabase
	TokenServiceConfig TokenServiceConfig
	Server             Server
	TokenCredentials   TokenCredentials
}

type PSQLDatabase struct {
	Driver       string `required:"true" split_word:"true"`
	User         string `required:"true" split_word:"true"`
	Password     string `required:"true" split_word:"true"`
	Host         string `required:"true" split_word:"true"`
	Port         string `required:"true" split_word:"true"`
	Name         string `required:"true" split_word:"true"`
	Timeout      int    `required:"true" split_word:"true"`
	DefaultLimit int    `required:"true" split_word:"true"`
	DefaultPage  int    `required:"true" split_word:"true"`
	Address      string `required:"false"`
}

type TokenServiceConfig struct {
	Protocol string `required:"true" split_word:"true"`
	Host     string `required:"true" split_word:"true"`
	Port     string `required:"true" split_word:"true"`
	Address  string `required:"false"`
}

type Server struct {
	Host string `required:"true" split_word:"true"`
	Port string `required:"true" split_word:"true"`
}

type TokenCredentials struct {
	SigningKey string `required:"true" split_word:"true"`
}

func Init() (Config, error) {
	// .env - for docker
	// local.env -  for local load
	err := godotenv.Load(".env")
	if err != nil {
		return Config{}, err
	}

	var cfg = Config{}

	psql, err := initPSQL()
	if err != nil {
		return Config{}, err
	}
	cfg.PSQLDatabase = psql

	tokenServiceConfig, err := initTokenService()
	if err != nil {
		return Config{}, err
	}
	cfg.TokenServiceConfig = tokenServiceConfig

	serverConfig, err := initServer()
	if err != nil {
		return Config{}, err
	}
	cfg.Server = serverConfig

	tokenCredentials, err := initTokenCredentials()
	if err != nil {
		return Config{}, err
	}
	cfg.TokenCredentials = tokenCredentials

	return cfg, nil
}

func initPSQL() (PSQLDatabase, error) {
	var params = map[string]string{
		PSQLDatabaseDriver:   "",
		PSQLDatabaseUser:     "",
		PSQLDatabasePassword: "",
		PSQLDatabaseHost:     "",
		PSQLDatabasePort:     "",
		PSQLDatabaseName:     "",
		PSQLDatabaseTimeout:  "",
		DefaultQueryLimit:    "",
		DefaultQueryPage:     "",
	}

	params, err := LookupEnvs(params)
	if err != nil {
		return PSQLDatabase{}, err
	}
	var db = PSQLDatabase{}

	db.Driver = params[PSQLDatabaseDriver]
	db.User = params[PSQLDatabaseUser]
	db.Password = params[PSQLDatabasePassword]
	db.Host = params[PSQLDatabaseHost]
	db.Port = params[PSQLDatabasePort]
	db.Name = params[PSQLDatabaseName]

	timeout, err := strconv.Atoi(params[PSQLDatabaseTimeout])
	if err != nil {
		return PSQLDatabase{}, err
	}
	db.Timeout = timeout

	limit, err := strconv.Atoi(params[DefaultQueryLimit])
	if err != nil {
		return PSQLDatabase{}, err
	}
	db.DefaultLimit = limit

	page, err := strconv.Atoi(params[DefaultQueryPage])
	if err != nil {
		return PSQLDatabase{}, err
	}
	db.DefaultPage = page

	db.Address = fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", db.Driver, db.User, db.Password, db.Host, db.Port, db.Name)

	return db, nil
}

func initTokenService() (TokenServiceConfig, error) {
	var params = map[string]string{
		TokenServiceProtocol: "",
		TokenServiceHost:     "",
		TokenServicePort:     "",
	}

	params, err := LookupEnvs(params)
	if err != nil {
		return TokenServiceConfig{}, err
	}

	var tokenServiceConfig = TokenServiceConfig{}

	tokenServiceConfig.Protocol = params[TokenServiceProtocol]
	tokenServiceConfig.Host = params[TokenServiceHost]
	tokenServiceConfig.Port = params[TokenServicePort]

	tokenServiceConfig.Address = fmt.Sprintf("%s://%s:%s", tokenServiceConfig.Protocol, tokenServiceConfig.Host, tokenServiceConfig.Port)

	return tokenServiceConfig, nil
}

func initServer() (Server, error) {
	var params = map[string]string{
		ServerHost: "",
		ServerPort: "",
	}

	params, err := LookupEnvs(params)
	if err != nil {
		return Server{}, nil
	}

	var serverConfig = Server{}

	serverConfig.Host = params[ServerHost]
	serverConfig.Port = params[ServerPort]

	return serverConfig, nil
}

func initTokenCredentials() (TokenCredentials, error) {
	var params = map[string]string{
		JWTSigningKey: "",
	}

	params, err := LookupEnvs(params)
	if err != nil {
		return TokenCredentials{}, nil
	}

	var tokenCredentials = TokenCredentials{}

	tokenCredentials.SigningKey = params[JWTSigningKey]

	return tokenCredentials, nil
}

func LookupEnvs(params map[string]string) (map[string]string, error) {
	var errorMsg string

	for i := range params {
		envVar, ok := os.LookupEnv(i)
		if !ok {
			errorMsg += fmt.Sprintf("\nCannot find %s", i)
		}
		params[i] = envVar
	}

	if len(errorMsg) > 0 {
		return nil, fmt.Errorf("%s", errorMsg)
	}
	return params, nil
}
