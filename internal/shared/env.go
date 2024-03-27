package shared

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	env := os.Getenv("SONIFY_ENV")
	if env == "" {
		env = "development"
	}

	// godotenv.Load(".env." + env + ".local")
	// if "test" != env {
	//   godotenv.Load(".env.local")
	// }

	envPath := ".env." + env

	log.Println("Loading env file: " + envPath)

	godotenv.Load(envPath)
}

type EnvVariables struct {
	IsDevelopment     bool
	IsProduction      bool
	Port              string
	DbUrl             string
	GoogleClientKey   string
	AccessTokenSecret string
	AccessTokenExp    int64
	PublicDIr         string
}

func ParseEnv() (*EnvVariables, error) {
	environment, err := getEnv("ENVIRONMENT")
	if err != nil {
		return nil, err
	}

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

	publicDir, err := getEnv("PUBLIC_DIR")
	if err != nil {
		return nil, err
	}

	return &EnvVariables{
		IsDevelopment:     environment == "development",
		IsProduction:      environment == "production",
		Port:              port,
		DbUrl:             dbUrl,
		GoogleClientKey:   googleClientKey,
		AccessTokenSecret: accessTokenSecret,
		AccessTokenExp:    accessTokenExp,
		PublicDIr:         publicDir,
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
