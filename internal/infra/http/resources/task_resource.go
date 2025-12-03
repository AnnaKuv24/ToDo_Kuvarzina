package resources

import (
    "time"
    "github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type TaskDto struct {
    Id          uint64              `json:"id"`
    UserId      uint64              `json:"userId"`
    Title       string              `json:"title"`
    Description *string             `json:"description,omitempty"`
    Status      domain.TaskStatus   `json:"status"`
    Priority    domain.TaskPriority `json:"priority"`
    Deadline    *time.Time          `json:"deadline,omitempty"`
    CreatedDate time.Time           `json:"createdDate"`
    UpdatedDate time.Time           `json:"updatedDate"`
}

type TasksDto struct {
    Tasks []TaskDto `json:"tasks"`
}

func (d TasksDto) DomainToDto(ts []domain.Task) TasksDto {
    tasks := make([]TaskDto, len(ts))
    for i, t := range ts {
        tasks[i] = TaskDto{}.DomainToDto(t)
    }

    return TasksDto{
        Tasks: tasks,
    }
}

func (d TaskDto) DomainToDto(t domain.Task) TaskDto {
    return TaskDto{
        Id:          t.Id,
        UserId:      t.UserId,
        Title:       t.Title,
        Description: t.Description,
        Status:      t.Status,
        Priority:    t.Priority,
        Deadline:    t.Deadline,
        CreatedDate: t.CreatedDate,
        UpdatedDate: t.UpdatedDate,
    }
}