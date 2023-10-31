package task

import (
	"errors"
	"orb-api/models"
	"orb-api/repositories/task"
)

func SetupTask(repository *task.Repository) *Service {
	return &Service{
		TaskRepo: repository,
	}
}

func (service *Service) ConcludedTask(id uint)(*models.Task, error) {
	return nil, nil
}

