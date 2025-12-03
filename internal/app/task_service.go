package app

import (
	"log"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
)

type TaskService interface {
	Save(t domain.Task) (domain.Task, error)
	FindList(tf database.TasksFilters) ([]domain.Task, error)
	Find(id uint64) (interface{}, error)
	Update(id uint64, t domain.Task) (domain.Task, error)
	Delete(id uint64) error
}

type taskService struct {
	taskRepo database.TaskRepository
}

func NewTaskService(tr database.TaskRepository) TaskService {
	return taskService{
		taskRepo: tr,
	}
}

func (s taskService) Save(t domain.Task) (domain.Task, error) {
	task, err := s.taskRepo.Save(t)
	if err != nil {
		log.Printf("taskService.Save(s.TaskRepo.Save): %s", err)
		return domain.Task{}, err
	}

	return task, nil
}

func (s taskService) FindList(tf database.TasksFilters) ([]domain.Task, error) {
	tasks, err := s.taskRepo.FindList(tf)
	if err != nil {
		log.Printf("taskService.FindList(s.taskRepo.FindList): %s", err)
		return nil, err
	}

	return tasks, nil
}

func (s taskService) Find(id uint64) (interface{}, error) {
	task, err := s.taskRepo.Find(id)
	if err != nil {
		log.Printf("taskService.FindList(s.taskRepo.Find): %s", err)
		return nil, err
	}

	return task, nil
}

func (s taskService) Update(id uint64, t domain.Task) (domain.Task, error) {
	task, err := s.taskRepo.Update(id, t)
	if err != nil {
		log.Printf("taskService.Update(s.taskRepo.Update): %s", err)
		return domain.Task{}, err
	}

	return task, nil
}

func (s taskService) Delete(id uint64) error {
	err := s.taskRepo.Delete(id)
	if err != nil {
		log.Printf("taskService.Delete(s.taskRepo.Delete): %s", err)
		return err
	}

	return nil
}
