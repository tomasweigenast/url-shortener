package taskmanager

import (
	"log"
	"os"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

var (
	enqueuer *work.Enqueuer
	pool     *work.WorkerPool
)

const (
	namespace              = "url_shortener_task_manager"
	report_url_hit_jobname = "url_hit"
)

func Start() {
	redisUrl, ok := os.LookupEnv("REDIS_URL")
	if !ok {
		log.Fatalf("REDIS_URL environment variable not found")
	}

	_ = redisUrl

	var redisPool = &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(redisUrl)
		},
	}

	enqueuer = work.NewEnqueuer(namespace, redisPool)
	pool = work.NewWorkerPool(context{}, 10, namespace, redisPool)
	pool.JobWithOptions(report_url_hit_jobname, work.JobOptions{
		Priority: 1,
		SkipDead: false,
		MaxFails: 5,
	}, (*context).urlhit)
	pool.Start()
}

func Stop() {
	pool.Stop()
	enqueuer.Pool.Close()
}

type context struct {
	metadata map[string]any
}
