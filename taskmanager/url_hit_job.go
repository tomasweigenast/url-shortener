package taskmanager

import "github.com/gocraft/work"

func (c *context) urlhit(job *work.Job) error {
	return nil
}

func EnqueueUrlHit() error {
	_, err := enqueuer.Enqueue(report_url_hit_jobname, work.Q{})
	return err
}
