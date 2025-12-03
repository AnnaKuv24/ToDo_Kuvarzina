package domain

import "time"

type Task struct {
    Id          uint64
    UserId      uint64
    Title       string
    Description *string
    Status      TaskStatus
    Priority    TaskPriority 
    Deadline    *time.Time
    CreatedDate time.Time
    UpdatedDate time.Time
    DeletedDate *time.Time
}

type TaskStatus string

const (
    NewTaskStatus        TaskStatus = "NEW"
    DoneTaskStatus       TaskStatus = "DONE"
    InProgressTaskStatus TaskStatus = "IN_PROGRESS"
)

type TaskPriority string

const (
    PriorityLow    TaskPriority = "LOW"
    PriorityMedium TaskPriority = "MEDIUM"
    PriorityHigh   TaskPriority = "HIGH"
)