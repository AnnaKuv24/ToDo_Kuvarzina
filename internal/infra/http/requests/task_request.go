package requests

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type TaskRequest struct {
	Title       string               `json:"title" validate:"required"`
	Description *string              `json:"description"`
	Status      *domain.TaskStatus   `json:"status"`
	Priority    *domain.TaskPriority `json:"priority"`
	Deadline    *int64               `json:"deadline"`
}

func (r TaskRequest) ToDomainModel() (interface{}, error) {
	var deadline time.Time
	if r.Deadline != nil {
		if *r.Deadline != 0 {
			deadline = time.Unix(*r.Deadline, 0)
		}
	}

	var dl *time.Time
	if !deadline.IsZero() {
		dl = &deadline
	}

	var status domain.TaskStatus
	if r.Status != nil {
		status = *r.Status
	} else {
		status = domain.NewTaskStatus
	}

	var priority domain.TaskPriority
	if r.Priority != nil {
		priority = *r.Priority
	} else {
		priority = domain.PriorityMedium
	}

	return domain.Task{
		Title:       r.Title,
		Description: r.Description,
		Status:      status,
		Priority:    priority,
		Deadline:    dl,
	}, nil
}
