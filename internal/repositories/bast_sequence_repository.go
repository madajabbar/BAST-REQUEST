package repositories

import (
	"bast-request/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BastSequenceRepository struct {
	db *gorm.DB
}

func NewBastSequenceRepository(db *gorm.DB) *BastSequenceRepository {
	return &BastSequenceRepository{db: db}
}

func (r *BastSequenceRepository) FindSequence(formatID string, year int, month int) (models.BastSequence, error) {
	var seq models.BastSequence
	err := r.db.Where("format_id = ? AND year = ? AND month = ?", formatID, year, month).First(&seq).Error
	return seq, err
}

// IncrementAndGet generates the next running number transactionally
func (r *BastSequenceRepository) IncrementAndGet(formatID string, year int, month int) (int, error) {
	var seq models.BastSequence
	
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Use clause.Locking to prevent race conditions during generation
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("format_id = ? AND year = ? AND month = ?", formatID, year, month).
			First(&seq).Error
			
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// Create new sequence for this month
				parsedUID, _ := uuid.Parse(formatID)
				seq = models.BastSequence{
					FormatID:   parsedUID,
					Year:       year,
					Month:      month,
					LastNumber: 1,
				}
				// We need to pass valid uuid. Wait, gorm handles string to uuid automatically if it's well-formed string. Let's fix this in service instead.
				return err // Return err to service to handle creation
			}
			return err
		}

		// Increment
		seq.LastNumber++
		if err := tx.Save(&seq).Error; err != nil {
			return err
		}
		
		return nil
	})

	if err != nil {
		return 0, err
	}
	
	return seq.LastNumber, nil
}

func (r *BastSequenceRepository) Create(seq *models.BastSequence) error {
	return r.db.Create(seq).Error
}

func (r *BastSequenceRepository) Update(seq *models.BastSequence) error {
	return r.db.Save(seq).Error
}
