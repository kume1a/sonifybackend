package main

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type EnvVariables struct {
	Port  string
	DbUrl string
}

func parseEnv() (*EnvVariables, error) {
	port, err := getEnvVariable("PORT")
	if err != nil {
		return nil, err
	}

	dbUrl, err := getEnvVariable("DB_URL")
	if err != nil {
		return nil, err
	}

	return &EnvVariables{
		Port:  port,
		DbUrl: dbUrl,
	}, nil
}

func getEnvVariable(key string) (string, error) {
	envVar := os.Getenv(key)
	if envVar == "" {
		msg := fmt.Sprintf("%v is not found in the env", key)

		log.Fatal(msg)
		return "", errors.New(msg)
	}

	return envVar, nil
}
