package services

import (
	"bast-request/internal/models"
	"bast-request/internal/repositories"
)

type ProjectService struct {
	repo *repositories.ProjectRepository
}

func NewProjectService(repo *repositories.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) GetAllProjects(customerID string) ([]models.Project, error) {
	return s.repo.FindAll(customerID)
}

func (s *ProjectService) GetProjectByID(id string) (models.Project, error) {
	return s.repo.FindByID(id)
}

func (s *ProjectService) CreateProject(project *models.Project) error {
	return s.repo.Create(project)
}

func (s *ProjectService) UpdateProject(id string, input *models.Project) (models.Project, error) {
	project, err := s.repo.FindByID(id)
	if err != nil {
		return project, err
	}

	project.ProjectCode = input.ProjectCode
	project.ProjectName = input.ProjectName
	project.Status = input.Status
	// Don't update CustomerID

	err = s.repo.Update(&project)
	return project, err
}

func (s *ProjectService) DeleteProject(id string) error {
	return s.repo.Delete(id)
}
