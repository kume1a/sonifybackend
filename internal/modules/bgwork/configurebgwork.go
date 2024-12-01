package bgwork

import (
	"log"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/kume1a/sonifybackend/internal/config"
	"github.com/kume1a/sonifybackend/internal/shared"
	"github.com/robfig/cron"
)

func ConfigureBackgroundWork(resourceConfig *config.ResourceConfig) (*work.Enqueuer, *work.WorkerPool) {
	crontab := cron.New()
	crontab.AddFunc("0 6 */2 * *", CreateHandleDeleteUnusedAudios(resourceConfig))

	crontab.Start()

	redisPool := createRedisPool()

	pool := registerWorkerPool(resourceConfig, redisPool)
	enqueuer := work.NewEnqueuer(shared.BackgroundJobNamespace, redisPool)

	return enqueuer, pool
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
			return redis.Dial(
				"tcp",
				envVars.RedisAddress,
				redis.DialPassword(envVars.RedisPassword),
			)
		},
	}
}

func registerWorkerPool(resourceConfig *config.ResourceConfig, redisPool *redis.Pool) *work.WorkerPool {
	pool := work.NewWorkerPool(struct{}{}, 2, shared.BackgroundJobNamespace, redisPool)

	pool.Job(shared.BackgroundJobDownloadPlaylistAudios, CreateHandleDownloadPlaylistAudios(resourceConfig))

	pool.Start()

	return pool
}
