package config

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

const (
	BackgroundJobNamespace = "application_namespace"

	BackgroundJobDownloadPlaylistAudios = "download_playlist_audios"
)

func ConfigureBackgroundWork() *work.Enqueuer {
	redisPool := createRedisPool()

	registerWorkerPool(redisPool)

	return work.NewEnqueuer(BackgroundJobNamespace, redisPool)
}

func createRedisPool() *redis.Pool {
	envVars, err := ParseEnv()
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

func registerWorkerPool(redisPool *redis.Pool) {
	pool := work.NewWorkerPool(struct{}{}, 2, BackgroundJobNamespace, redisPool)

	pool.Job(BackgroundJobDownloadPlaylistAudios, handleDownloadPlaylistAudios)

	pool.Start()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	pool.Stop()
}

func handleDownloadPlaylistAudios(job *work.Job) error {
	playlistId := job.ArgString("playlistId")
	if err := job.ArgError(); err != nil {
		return err
	}

	fmt.Println("Downloading playlist: ", playlistId)

	return nil
}
