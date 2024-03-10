package taskmanager

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"tomasweigenast.com/url-shortener/database"
)

const (
	namespace              = "url_shortener_task_manager"
	report_url_hit_jobname = "url_hit"
)

var riverClient *river.Client[pgx.Tx]

func Start() {
	workers := river.NewWorkers()
	river.AddWorker(workers, &PushHitWorker{})

	var err error
	riverClient, err = river.NewClient(riverpgxv5.New(database.Pool()), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {
				MaxWorkers: 100,
			},
		},
		Workers: workers,
	})

	if err != nil {
		log.Fatalf("unable to start river client: %s", err)
	}

	if err := riverClient.Start(context.Background()); err != nil {
		log.Printf("unable to execute job: %s", err)
	}
}

func Stop() {
	riverClient.Stop(context.TODO())
	riverClient = nil
}
