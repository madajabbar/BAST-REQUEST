package repositories

import (
	"bast-request/internal/models"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) FindAll(customerID string) ([]models.Project, error) {
	var projects []models.Project
	query := r.db.Preload("Customer")

	if customerID != "" {
		query = query.Where("customer_id = ?", customerID)
	}

	err := query.Find(&projects).Error
	return projects, err
}

func (r *ProjectRepository) FindByID(id string) (models.Project, error) {
	var project models.Project
	err := r.db.Preload("Customer").First(&project, "project_id = ?", id).Error
	return project, err
}

func (r *ProjectRepository) Create(project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *ProjectRepository) Update(project *models.Project) error {
	return r.db.Save(project).Error
}

func (r *ProjectRepository) Delete(id string) error {
	return r.db.Model(&models.Project{}).Where("project_id = ?", id).Update("status", "inactive").Error
}
