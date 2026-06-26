package repositories

import (
	"bast-request/internal/models"
	"gorm.io/gorm"
)

type BastRequestRepository struct {
	db *gorm.DB
}

func NewBastRequestRepository(db *gorm.DB) *BastRequestRepository {
	return &BastRequestRepository{db: db}
}

func (r *BastRequestRepository) FindAll(customerID, projectID, status string) ([]models.BastRequest, error) {
	var requests []models.BastRequest
	query := r.db.Preload("Customer").Preload("Project").Preload("Format")

	if customerID != "" {
		query = query.Where("customer_id = ?", customerID)
	}
	if projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Find(&requests).Error
	return requests, err
}

func (r *BastRequestRepository) FindByID(id string) (models.BastRequest, error) {
	var request models.BastRequest
	err := r.db.Preload("Customer").Preload("Project").Preload("Format").First(&request, "bast_request_id = ?", id).Error
	return request, err
}

func (r *BastRequestRepository) Create(request *models.BastRequest) error {
	return r.db.Create(request).Error
}

func (r *BastRequestRepository) UpdateStatus(id string, status string) error {
	return r.db.Model(&models.BastRequest{}).Where("bast_request_id = ?", id).Update("status", status).Error
}
