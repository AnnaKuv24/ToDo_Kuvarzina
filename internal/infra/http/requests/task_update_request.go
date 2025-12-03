package requests

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type TaskUpdateRequest struct {
	Title       *string              `json:"title"`
	Description *string              `json:"description"`
	Status      *domain.TaskStatus   `json:"status"`
	Priority    *domain.TaskPriority `json:"priority"`
	Deadline    *int64               `json:"deadline"`
}

func (r TaskUpdateRequest) ToDomainModel() (interface{}, error) {
	var task domain.Task

	if r.Title != nil {
		task.Title = *r.Title
	}

	if r.Description != nil {
		task.Description = r.Description
	}

	if r.Status != nil {
		task.Status = *r.Status
	}

	if r.Priority != nil {
		task.Priority = *r.Priority
	}

	if r.Deadline != nil {
		if *r.Deadline != 0 {
			deadline := time.Unix(*r.Deadline, 0)
			task.Deadline = &deadline
		} else {
			task.Deadline = nil
		}
	}

	return task, nil
}
