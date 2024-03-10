package taskmanager

import (
	"context"
	"log"

	"github.com/riverqueue/river"
	"tomasweigenast.com/url-shortener/database"
	"tomasweigenast.com/url-shortener/entities"
)

type PushHitArgs struct {
	Urlhit entities.UrlHit `json:"hit"`
}

func (PushHitArgs) Kind() string { return "push_url_hit" }

type PushHitWorker struct {
	river.WorkerDefaults[PushHitArgs]
}

func (w *PushHitWorker) Work(ctx context.Context, job *river.Job[PushHitArgs]) error {
	log.Printf("[PushHitWorker] [job=%d] pushing url [%d] hit to database\n", job.ID, job.Args.Urlhit.UrlId)
	return database.InsertUrlHit(ctx, &job.Args.Urlhit)
}

func EnqueueUrlHit(hit entities.UrlHit) {
	go func() {
		_, err := riverClient.Insert(context.Background(), PushHitArgs{hit}, &river.InsertOpts{
			Priority:    1,
			MaxAttempts: 10,
		})

		if err != nil {
			log.Println("unable to push job to postgres. Error:", err)
		}
	}()
}
