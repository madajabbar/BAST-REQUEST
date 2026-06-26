package services

import (
	"bast-request/internal/models"
	"bast-request/internal/repositories"
)

type AuditLogService struct {
	repo *repositories.AuditLogRepository
}

func NewAuditLogService(repo *repositories.AuditLogRepository) *AuditLogService {
	return &AuditLogService{repo: repo}
}

func (s *AuditLogService) GetAllAuditLogs(tableName, recordID, performedBy, dateFrom, dateTo string) ([]models.AuditLog, error) {
	return s.repo.FindAll(tableName, recordID, performedBy, dateFrom, dateTo)
}

func (s *AuditLogService) GetAuditLogByID(id string) (models.AuditLog, error) {
	return s.repo.FindByID(id)
}

func (s *AuditLogService) LogAction(log *models.AuditLog) error {
	return s.repo.Create(log)
}
