package repository

import (
	"sync"
	"time"

	"video-platform/model"
)

type ReportRepository struct {
	lock  sync.RWMutex
	next  int
	items []model.Report
}

func NewReportRepository() *ReportRepository {
	return &ReportRepository{
		next:  1,
		items: []model.Report{},
	}
}

func (r *ReportRepository) Add(targetType string, targetID string, userID string, reason string) model.Report {
	r.lock.Lock()
	defer r.lock.Unlock()
	rep := model.Report{
		ID:         r.next,
		TargetType: targetType,
		TargetID:   targetID,
		UserID:     userID,
		Reason:     reason,
		CreatedAt:  time.Now().Unix(),
	}
	r.next++
	r.items = append(r.items, rep)
	return rep
}

func (r *ReportRepository) List() []model.Report {
	r.lock.RLock()
	defer r.lock.RUnlock()
	out := make([]model.Report, len(r.items))
	copy(out, r.items)
	return out
}
