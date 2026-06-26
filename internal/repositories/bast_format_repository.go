package repositories

import (
	"bast-request/internal/models"
	"gorm.io/gorm"
)

type BastFormatRepository struct {
	db *gorm.DB
}

func NewBastFormatRepository(db *gorm.DB) *BastFormatRepository {
	return &BastFormatRepository{db: db}
}

func (r *BastFormatRepository) FindAll() ([]models.BastFormat, error) {
	var formats []models.BastFormat
	err := r.db.Find(&formats).Error
	return formats, err
}

func (r *BastFormatRepository) FindByID(id string) (models.BastFormat, error) {
	var format models.BastFormat
	err := r.db.First(&format, "format_id = ?", id).Error
	return format, err
}

func (r *BastFormatRepository) Create(format *models.BastFormat) error {
	return r.db.Create(format).Error
}

func (r *BastFormatRepository) Update(format *models.BastFormat) error {
	return r.db.Save(format).Error
}

func (r *BastFormatRepository) Delete(id string) error {
	return r.db.Model(&models.BastFormat{}).Where("format_id = ?", id).Update("is_active", false).Error
}
