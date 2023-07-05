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
)

type Config struct {
	PSQLDatabase PSQLDatabase
}

type PSQLDatabase struct {
	Driver   string `required:"true" split_word:"true"`
	User     string `required:"true" split_word:"true"`
	Password string `required:"true" split_word:"true"`
	Host     string `required:"true" split_word:"true"`
	Port     string `required:"true" split_word:"true"`
	Name     string `required:"true" split_word:"true"`
	Timeout  int    `required:"true" split_word:"true"`
	Address  string `required:"false"`
}

func Init() (Config, error) {
	var cfg = Config{}

	psql, err := initPSQL()
	if err != nil {
		return Config{}, err
	}
	cfg.PSQLDatabase = psql

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

	db.Address = fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", db.Driver, db.User, db.Password, db.Host, db.Port, db.Name)

	return db, nil
}

func LookupEnvs(params map[string]string) (map[string]string, error) {
	var errorMsg string

	// .env - for docker
	// local.env -  for local load
	err := godotenv.Load("local.env")
	if err != nil {
		return nil, err
	}

	for i := range params {
		envVar, ok := os.LookupEnv(i)
		if !ok {
			errorMsg += fmt.Sprintf("\nCannot find %s", i)
		}
		params[i] = envVar
	}

	if len(errorMsg) > 0 {
		return nil, err
	}
	return params, nil
}
