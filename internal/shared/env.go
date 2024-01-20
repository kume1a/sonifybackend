package shared

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

type EnvVariables struct {
	Port              string
	DbUrl             string
	GoogleClientKey   string
	AccessTokenSecret string
	AccessTokenExp    int64
}

func ParseEnv() (*EnvVariables, error) {
	port, err := getEnv("PORT")
	if err != nil {
		return nil, err
	}

	dbUrl, err := getEnv("DB_URL")
	if err != nil {
		return nil, err
	}

	googleClientKey, err := getEnv("GOOGLE_CLIENT_KEY")
	if err != nil {
		return nil, err
	}

	accessTokenSecret, err := getEnv("ACCESS_TOKEN_SECRET")
	if err != nil {
		return nil, err
	}

	accessTokenExp, err := getEnvInt("ACCESS_TOKEN_EXP")
	if err != nil {
		return nil, err
	}

	return &EnvVariables{
		Port:              port,
		DbUrl:             dbUrl,
		GoogleClientKey:   googleClientKey,
		AccessTokenSecret: accessTokenSecret,
		AccessTokenExp:    accessTokenExp,
	}, nil
}

func getEnvInt(key string) (int64, error) {
	val, err := getEnv(key)
	if err != nil {
		return 0, err
	}

	valInt, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}

	return valInt, nil
}

func getEnv(key string) (string, error) {
	envVar := os.Getenv(key)
	if envVar == "" {
		msg := fmt.Sprintf("%v is not found in the env", key)

		log.Fatal(msg)
		return "", errors.New(msg)
	}

	return envVar, nil
}
