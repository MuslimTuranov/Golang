package store

import (
	"Assignment1/internal/models"
	"errors"
	"sync"
)

type Storage struct {
	tasks  map[int]models.Task
	nextID int
	mu     sync.Mutex
}

func (s *Storage) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.tasks[id]
	if !exists {
		return errors.New("task not found")
	}

	delete(s.tasks, id)
	return nil
}

func NewStorage() *Storage {
	return &Storage{
		tasks:  make(map[int]models.Task),
		nextID: 1,
	}
}

func (s *Storage) GetAll() []models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []models.Task
	for _, task := range s.tasks {
		result = append(result, task)
	}
	return result
}

func (s *Storage) GetByID(id int) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[id]
	if !exists {
		return models.Task{}, errors.New("task not found")
	}
	return task, nil
}

func (s *Storage) Create(title string) models.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := models.Task{
		ID:    s.nextID,
		Title: title,
		Done:  false,
	}
	s.tasks[s.nextID] = task
	s.nextID++
	return task
}

func (s *Storage) UpdateDone(id int, done bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[id]
	if !exists {
		return errors.New("task not found")
	}

	task.Done = done
	s.tasks[id] = task
	return nil
}
