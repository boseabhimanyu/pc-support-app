package services

import "github.com/boseabhimanyu/pc-support-app/backend/pc-support-system/internal/repository"

type JobService struct {
	jobRepo    repository.JobRepository
	userRepo   repository.UserRepository
	deviceRepo repository.DeviceRepository
}

func NewJobService(
	jobRepo repository.JobRepository,
	userRepo repository.UserRepository,
	deviceRepo repository.DeviceRepository,
) *JobService {

	return &JobService{
		jobRepo:    jobRepo,
		userRepo:   userRepo,
		deviceRepo: deviceRepo,
	}
}
