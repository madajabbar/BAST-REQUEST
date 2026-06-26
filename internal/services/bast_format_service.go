package services

import (
	"bast-request/internal/models"
	"bast-request/internal/repositories"
)

type BastFormatService struct {
	repo *repositories.BastFormatRepository
}

func NewBastFormatService(repo *repositories.BastFormatRepository) *BastFormatService {
	return &BastFormatService{repo: repo}
}

func (s *BastFormatService) GetAllFormats() ([]models.BastFormat, error) {
	return s.repo.FindAll()
}

func (s *BastFormatService) GetFormatByID(id string) (models.BastFormat, error) {
	return s.repo.FindByID(id)
}

func (s *BastFormatService) CreateFormat(format *models.BastFormat) error {
	return s.repo.Create(format)
}

func (s *BastFormatService) UpdateFormat(id string, input *models.BastFormat) (models.BastFormat, error) {
	format, err := s.repo.FindByID(id)
	if err != nil {
		return format, err
	}

	format.FormatName = input.FormatName
	format.FormatType = input.FormatType
	format.FormatPattern = input.FormatPattern
	format.IsActive = input.IsActive

	err = s.repo.Update(&format)
	return format, err
}

func (s *BastFormatService) DeleteFormat(id string) error {
	return s.repo.Delete(id)
}
