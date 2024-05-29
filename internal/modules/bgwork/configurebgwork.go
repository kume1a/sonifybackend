package bgwork

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/shared"
)

func ConfigureBackgroundWork(resourceConfig *config.ResourceConfig) *work.Enqueuer {
	redisPool := createRedisPool()

	registerWorkerPool(resourceConfig, redisPool)

	return work.NewEnqueuer(shared.BackgroundJobNamespace, redisPool)
}

func createRedisPool() *redis.Pool {
	envVars, err := config.ParseEnv()
	if err != nil {
		log.Fatal(err)
	}

	return &redis.Pool{
		MaxActive: 2,
		MaxIdle:   100,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", ":6379", redis.DialPassword(envVars.RedisPassword))
		},
	}
}

func registerWorkerPool(resourceConfig *config.ResourceConfig, redisPool *redis.Pool) {
	pool := work.NewWorkerPool(struct{}{}, 2, shared.BackgroundJobNamespace, redisPool)

	pool.Job(shared.BackgroundJobDownloadPlaylistAudios, CreateHandleDownloadPlaylistAudios(resourceConfig))

	pool.Start()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	pool.Stop()
}
