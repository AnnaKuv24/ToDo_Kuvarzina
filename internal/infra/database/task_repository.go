package database

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const TasksTableName = "tasks"

type task struct {
	Id          uint64              `db:"id,omitempty"`
	UserId      uint64              `db:"user_id"`
	Title       string              `db:"title"`
	Description *string             `db:"description"`
	Status      domain.TaskStatus   `db:"status"`
	Priority    domain.TaskPriority `db:"priority"`
	Deadline    *time.Time          `db:"deadline"`
	CreatedDate time.Time           `db:"created_date"`
	UpdatedDate time.Time           `db:"updated_date"`
	DeletedDate *time.Time          `db:"deleted_date"`
}

type TaskRepository interface {
	Save(t domain.Task) (domain.Task, error)
	FindList(tf TasksFilters) ([]domain.Task, error)
	Find(id uint64) (domain.Task, error)
	Update(id uint64, t domain.Task) (domain.Task, error)
	Delete(id uint64) error
}

type taskRepository struct {
	sess db.Session
	coll db.Collection
}

func NewTaskRepository(dbSession db.Session) TaskRepository {
	return taskRepository{
		sess: dbSession,
		coll: dbSession.Collection(TasksTableName),
	}
}

func (r taskRepository) Save(t domain.Task) (domain.Task, error) {
	tsk := r.mapDomainToModel(t)
	tsk.CreatedDate = time.Now()
	tsk.UpdatedDate = time.Now()

	err := r.coll.InsertReturning(&tsk)
	if err != nil {
		return domain.Task{}, err
	}

	t = r.mapModelToDomain(tsk)
	return t, nil
}

func (r taskRepository) Find(id uint64) (domain.Task, error) {
	var t task

	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&t)
	if err != nil {
		return domain.Task{}, err
	}

	return r.mapModelToDomain(t), nil
}

type TasksFilters struct {
	UserId       uint64
	Status       string
	Search       string
	Priority     string
	DeadlineFrom *time.Time
	DeadlineTo   *time.Time
	FilterType   string // "today", "week", "overdue"
}

func (r taskRepository) FindList(tf TasksFilters) ([]domain.Task, error) {
	var tasks []task

	conditions := db.Cond{"user_id": tf.UserId, "deleted_date": nil}

	if tf.Status != "" {
		if tf.Status == "NOT_DONE" {
			conditions["status !="] = domain.DoneTaskStatus
		} else {
			conditions["status"] = tf.Status
		}
	}

	if tf.Priority != "" {
		conditions["priority"] = tf.Priority
	}

	if tf.DeadlineFrom != nil && tf.DeadlineTo != nil {
		conditions["deadline >="] = *tf.DeadlineFrom
		conditions["deadline <="] = *tf.DeadlineTo
	} else if tf.DeadlineFrom != nil {
		conditions["deadline >="] = *tf.DeadlineFrom
	} else if tf.DeadlineTo != nil {
		conditions["deadline <="] = *tf.DeadlineTo
	}

	if tf.Search != "" {
		conditions["title ILIKE"] = "%" + tf.Search + "%"
	}

	query := r.coll.Find(conditions)

	query = query.OrderBy("deadline ASC", "created_date DESC")

	err := query.All(&tasks)
	if err != nil {
		return nil, err
	}

	return r.mapModelToDomainCollection(tasks), nil
}

func (r taskRepository) Update(id uint64, t domain.Task) (domain.Task, error) {
	var existingTask task
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&existingTask)
	if err != nil {
		return domain.Task{}, err
	}

	if t.Title != "" {
		existingTask.Title = t.Title
	}
	if t.Description != nil {
		existingTask.Description = t.Description
	}
	if t.Status != "" {
		existingTask.Status = t.Status
	}
	if t.Priority != "" {
		existingTask.Priority = t.Priority
	}
	if t.Deadline != nil {
		existingTask.Deadline = t.Deadline
	}

	existingTask.UpdatedDate = time.Now()

	err = r.coll.UpdateReturning(&existingTask)
	if err != nil {
		return domain.Task{}, err
	}

	return r.mapModelToDomain(existingTask), nil
}

func (r taskRepository) Delete(id uint64) error {
	return r.coll.Find(db.Cond{"id": id}).Update(map[string]interface{}{
		"deleted_date": time.Now(),
	})
}

func (r taskRepository) mapDomainToModel(t domain.Task) task {
	return task{
		Id:          t.Id,
		UserId:      t.UserId,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Priority:    t.Priority,
		Deadline:    t.Deadline,
		CreatedDate: t.CreatedDate,
		UpdatedDate: t.UpdatedDate,
		DeletedDate: t.DeletedDate,
	}
}

func (r taskRepository) mapModelToDomain(t task) domain.Task {
	return domain.Task{
		Id:          t.Id,
		UserId:      t.UserId,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Priority:    t.Priority,
		Deadline:    t.Deadline,
		CreatedDate: t.CreatedDate,
		UpdatedDate: t.UpdatedDate,
		DeletedDate: t.DeletedDate,
	}
}

func (r taskRepository) mapModelToDomainCollection(ts []task) []domain.Task {
	tasks := make([]domain.Task, len(ts))
	for i, t := range ts {
		tasks[i] = r.mapModelToDomain(t)
	}
	return tasks
}
