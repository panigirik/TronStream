package background_jobs

import (
	"context"
	"sync"
)

type Sheduler struct {
	jobs []Job
}

type Job interface {
	Start(ctx context.Context)
}

func NewScheduler(jobs ...Job) *Sheduler {
	return &Sheduler{
		jobs: jobs,
	}
}

func (s *Sheduler) Start(ctx context.Context) {
	var wg sync.WaitGroup
	for _, job := range s.jobs {
		wg.Add(1)

		go func(j Job) {
			defer wg.Done()
			j.Start(ctx)
		}(job)
	}

	<-ctx.Done()
	wg.Wait()
}
