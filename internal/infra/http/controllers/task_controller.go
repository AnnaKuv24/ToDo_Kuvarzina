package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
)

type TaskController struct {
	taskService app.TaskService
}

func NewTaskController(ts app.TaskService) TaskController {
	return TaskController{
		taskService: ts,
	}
}

func (c TaskController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		task, err := requests.Bind(r, requests.TaskRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController.Save(requests.Bind): %s", err)
			BadRequest(w, err)
			return
		}

		task.UserId = user.Id
		task.Status = domain.NewTaskStatus
		task, err = c.taskService.Save(task)
		if err != nil {
			log.Printf("TaskController.Save(c.taskService.Save): %s", err)
			InternalServerError(w, err)
			return
		}

		taskDto := resources.TaskDto{}
		taskDto = taskDto.DomainToDto(task)
		Success(w, taskDto)
	}
}

func (c TaskController) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		task := r.Context().Value(TaskKey).(domain.Task)

		if user.Id != task.UserId {
			err := errors.New("access denied")
			Forbidden(w, err)
			return
		}

		taskDto := resources.TaskDto{}
		taskDto = taskDto.DomainToDto(task)
		Success(w, taskDto)
	}
}

func (c TaskController) FindList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)

		status := ""
		if r.URL.Query().Has("status") {
			status = r.URL.Query().Get("status")
		}

		search := ""
		if r.URL.Query().Has("search") {
			search = r.URL.Query().Get("search")
		}

		priority := ""
		if r.URL.Query().Has("priority") {
			priority = r.URL.Query().Get("priority")
		}

		var deadlineFrom, deadlineTo *time.Time

		if r.URL.Query().Has("deadline_from") {
			deadlineStr := r.URL.Query().Get("deadline_from")
			if unixTime, err := strconv.ParseInt(deadlineStr, 10, 64); err == nil && unixTime > 0 {
				t := time.Unix(unixTime, 0)
				deadlineFrom = &t
			}
		}

		if r.URL.Query().Has("deadline_to") {
			deadlineStr := r.URL.Query().Get("deadline_to")
			if unixTime, err := strconv.ParseInt(deadlineStr, 10, 64); err == nil && unixTime > 0 {
				t := time.Unix(unixTime, 0)
				deadlineTo = &t
			}
		}

		filterType := ""
		if r.URL.Query().Has("filter_type") {
			filterType = r.URL.Query().Get("filter_type")

			now := time.Now()
			switch filterType {
			case "today":
				startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
				endOfDay := startOfDay.Add(24 * time.Hour)
				deadlineFrom = &startOfDay
				deadlineTo = &endOfDay

			case "week":
				weekday := now.Weekday()
				daysToMonday := (weekday - time.Monday + 7) % 7 // Понеділок
				startOfWeek := now.AddDate(0, 0, -int(daysToMonday))
				startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())
				endOfWeek := startOfWeek.AddDate(0, 0, 7)
				deadlineFrom = &startOfWeek
				deadlineTo = &endOfWeek

			case "overdue":
				// Прострочені задачі (deadline в минулому, статус не DONE)
				deadlineTo = &now
				if status == "" {
					status = "NOT_DONE"
				}
			}
		}

		filters := database.TasksFilters{
			UserId:       user.Id,
			Status:       status,
			Search:       search,
			Priority:     priority,
			DeadlineFrom: deadlineFrom,
			DeadlineTo:   deadlineTo,
			FilterType:   filterType,
		}

		tasks, err := c.taskService.FindList(filters)
		if err != nil {
			log.Printf("TaskController.FindList(c.taskService.FindList): %s", err)
			InternalServerError(w, err)
			return
		}

		tasksDto := resources.TasksDto{}
		tasksDto = tasksDto.DomainToDto(tasks)
		Success(w, tasksDto)
	}
}

func (c TaskController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		task := r.Context().Value(TaskKey).(domain.Task)

		if user.Id != task.UserId {
			err := errors.New("access denied")
			Forbidden(w, err)
			return
		}

		updateData, err := requests.Bind(r, requests.TaskUpdateRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController.Update(requests.Bind): %s", err)
			BadRequest(w, err)
			return
		}

		if updateData.Status != "" {
			if !c.isValidStatus(updateData.Status) {
				err := errors.New("invalid status value")
				BadRequest(w, err)
				return
			}
		}

		if updateData.Priority != "" {
			if !c.isValidPriority(updateData.Priority) {
				err := errors.New("invalid priority value")
				BadRequest(w, err)
				return
			}
		}

		updatedTask, err := c.taskService.Update(task.Id, updateData)
		if err != nil {
			log.Printf("TaskController.Update(c.taskService.Update): %s", err)
			InternalServerError(w, err)
			return
		}

		taskDto := resources.TaskDto{}
		taskDto = taskDto.DomainToDto(updatedTask)
		Success(w, taskDto)
	}
}

func (c TaskController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		task := r.Context().Value(TaskKey).(domain.Task)

		if user.Id != task.UserId {
			err := errors.New("access denied")
			Forbidden(w, err)
			return
		}

		err := c.taskService.Delete(task.Id)
		if err != nil {
			log.Printf("TaskController.Delete(c.taskService.Delete): %s", err)
			InternalServerError(w, err)
			return
		}

		Success(w, nil)
	}
}

func (c TaskController) isValidStatus(status domain.TaskStatus) bool {
	switch status {
	case domain.NewTaskStatus, domain.InProgressTaskStatus, domain.DoneTaskStatus:
		return true
	default:
		return false
	}
}

func (c TaskController) isValidPriority(priority domain.TaskPriority) bool {
	switch priority {
	case domain.PriorityLow, domain.PriorityMedium, domain.PriorityHigh:
		return true
	default:
		return false
	}
}
