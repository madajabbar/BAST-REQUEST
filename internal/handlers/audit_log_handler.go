package handlers

import (
	_ "bast-request/internal/models"
	"bast-request/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuditLogHandler struct {
	service *services.AuditLogService
}

func NewAuditLogHandler(service *services.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{service: service}
}

// GetAllAuditLogs godoc
// @Summary Get all audit logs
// @Description Retrieve a list of audit logs, optionally filtered by parameters
// @Tags audit-logs
// @Produce json
// @Param table_name query string false "Table Name filter"
// @Param record_id query string false "Record ID filter"
// @Param performed_by query string false "Performed By filter"
// @Param date_from query string false "Date From filter (YYYY-MM-DD)"
// @Param date_to query string false "Date To filter (YYYY-MM-DD)"
// @Success 200 {array} models.AuditLog
// @Failure 500 {object} map[string]interface{}
// @Router /audit-logs [get]
func (h *AuditLogHandler) GetAllAuditLogs(c *gin.Context) {
	tableName := c.Query("table_name")
	recordID := c.Query("record_id")
	performedBy := c.Query("performed_by")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	logs, err := h.service.GetAllAuditLogs(tableName, recordID, performedBy, dateFrom, dateTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}

// GetAuditLogByID godoc
// @Summary Get an audit log by ID
// @Description Retrieve a single audit log based on its ID
// @Tags audit-logs
// @Produce json
// @Param id path string true "Audit Log ID"
// @Success 200 {object} models.AuditLog
// @Failure 404 {object} map[string]interface{}
// @Router /audit-logs/{id} [get]
func (h *AuditLogHandler) GetAuditLogByID(c *gin.Context) {
	id := c.Param("id")
	log, err := h.service.GetAuditLogByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Audit log not found"})
		return
	}
	c.JSON(http.StatusOK, log)
}
