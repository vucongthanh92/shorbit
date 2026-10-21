package application

import (
	"context"
	"github.com/google/uuid"
	"github.com/vucongthanh92/shortbit/internal/share/domain"
	"time"
)

type ViewWorker struct {
	repo domain.ShareLinkRepository
	jobs chan viewJob
}
type viewJob struct {
	id uuid.UUID
	at time.Time
}

func NewViewWorker(repo domain.ShareLinkRepository, buffer int) *ViewWorker {
	return &ViewWorker{repo: repo, jobs: make(chan viewJob, buffer)}
}
func (w *ViewWorker) Record(id uuid.UUID, at time.Time) {
	select {
	case w.jobs <- viewJob{id, at}:
	default:
	}
}
func (w *ViewWorker) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case job := <-w.jobs:
				_ = w.repo.RecordView(context.Background(), job.id, job.at)
			case <-ctx.Done():
				return
			}
		}
	}()
}
