package repositories

import (
	"bast-request/internal/models"
	"gorm.io/gorm"
)

type AuditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) FindAll(targetTable, recordID, performedBy, dateFrom, dateTo string) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	query := r.db.Model(&models.AuditLog{})

	if targetTable != "" {
		query = query.Where("table_name = ?", targetTable)
	}
	if recordID != "" {
		query = query.Where("record_id = ?", recordID)
	}
	if performedBy != "" {
		query = query.Where("performed_by = ?", performedBy)
	}
	if dateFrom != "" {
		query = query.Where("created_at >= ?", dateFrom)
	}
	if dateTo != "" {
		query = query.Where("created_at <= ?", dateTo)
	}

	err := query.Order("created_at DESC").Find(&logs).Error
	return logs, err
}

func (r *AuditLogRepository) FindByID(id string) (models.AuditLog, error) {
	var log models.AuditLog
	err := r.db.First(&log, "audit_log_id = ?", id).Error
	return log, err
}

func (r *AuditLogRepository) Create(log *models.AuditLog) error {
	return r.db.Create(log).Error
}
