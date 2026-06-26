package services

import (
	"bast-request/internal/models"
	"bast-request/internal/repositories"
	"fmt"

	"github.com/google/uuid"
)

type BastSequenceService struct {
	repo *repositories.BastSequenceRepository
}

func NewBastSequenceService(repo *repositories.BastSequenceRepository) *BastSequenceService {
	return &BastSequenceService{repo: repo}
}

func (s *BastSequenceService) GetSequence(formatID string, year int, month int) (models.BastSequence, error) {
	return s.repo.FindSequence(formatID, year, month)
}

func (s *BastSequenceService) ResetSequence(formatID string, year int, month int, lastNumber int) (models.BastSequence, error) {
	seq, err := s.repo.FindSequence(formatID, year, month)
	if err != nil {
		// If not found, create new
		uid, _ := uuid.Parse(formatID)
		seq = models.BastSequence{
			FormatID:   uid,
			Year:       year,
			Month:      month,
			LastNumber: lastNumber,
		}
		err = s.repo.Create(&seq)
		return seq, err
	}

	seq.LastNumber = lastNumber
	err = s.repo.Update(&seq)
	return seq, err
}

func (s *BastSequenceService) GenerateNextNumber(formatID string, year int, month int) (int, error) {
	nextNumber, err := s.repo.IncrementAndGet(formatID, year, month)
	if err != nil {
		if err.Error() == "record not found" {
			// Create first record
			uid, errParse := uuid.Parse(formatID)
			if errParse != nil {
				return 0, fmt.Errorf("invalid format_id: %v", errParse)
			}
			
			seq := models.BastSequence{
				FormatID:   uid,
				Year:       year,
				Month:      month,
				LastNumber: 1,
			}
			if errCreate := s.repo.Create(&seq); errCreate != nil {
				return 0, errCreate
			}
			return 1, nil
		}
		return 0, err
	}
	return nextNumber, nil
}
