package config

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

	if env == "development" {
		envPath := ".env." + env

		log.Println("Loading env file: " + envPath)

		godotenv.Load(envPath)
	}
}

type EnvVariables struct {
	IsDevelopment        bool
	IsProduction         bool
	Port                 string
	DbUrl                string
	GoogleClientKey      string
	AccessTokenSecret    string
	AccessTokenExpMillis int64
	PublicDIr            string
	MaxUploadSizeBytes   int64
	SpotifyClientID      string
	SpotifyClientSecret  string
	SpotifyRedirectURI   string
	RedisPassword        string
	RedisAddress         string
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

	accessTokenExpMillis, err := getEnvInt("ACCESS_TOKEN_EXP_MILLIS")
	if err != nil {
		return nil, err
	}

	publicDir, err := getEnv("PUBLIC_DIR")
	if err != nil {
		return nil, err
	}

	maxUploadSizeBytes, err := getEnvInt("MAX_UPLOAD_SIZE_BYTES")
	if err != nil {
		return nil, err
	}

	spotifyClientID, err := getEnv("SPOTIFY_CLIENT_ID")
	if err != nil {
		return nil, err
	}

	spotifyClientSecret, err := getEnv("SPOTIFY_CLIENT_SECRET")
	if err != nil {
		return nil, err
	}

	spotifyRedirectURI, err := getEnv("SPOTIFY_REDIRECT_URI")
	if err != nil {
		return nil, err
	}

	redisPassword, err := getEnv("REDIS_PASSWORD")
	if err != nil {
		return nil, err
	}

	redisAddress, err := getEnv("REDIS_ADDRESS")
	if err != nil {
		return nil, err
	}

	return &EnvVariables{
		IsDevelopment:        environment == "development",
		IsProduction:         environment == "production",
		Port:                 port,
		DbUrl:                dbUrl,
		GoogleClientKey:      googleClientKey,
		AccessTokenSecret:    accessTokenSecret,
		AccessTokenExpMillis: accessTokenExpMillis,
		PublicDIr:            publicDir,
		MaxUploadSizeBytes:   maxUploadSizeBytes,
		SpotifyClientID:      spotifyClientID,
		SpotifyClientSecret:  spotifyClientSecret,
		SpotifyRedirectURI:   spotifyRedirectURI,
		RedisPassword:        redisPassword,
		RedisAddress:         redisAddress,
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
