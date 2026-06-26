package handlers

import (
	"bast-request/internal/models"
	"bast-request/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BastRequestHandler struct {
	service *services.BastRequestService
}

func NewBastRequestHandler(service *services.BastRequestService) *BastRequestHandler {
	return &BastRequestHandler{service: service}
}

// GetAllRequests godoc
// @Summary Get all BAST requests
// @Description Retrieve a list of BAST requests, optionally filtered by customer_id, project_id, and status
// @Tags bast-requests
// @Produce json
// @Param customer_id query string false "Customer ID filter"
// @Param project_id query string false "Project ID filter"
// @Param status query string false "Status filter"
// @Success 200 {array} models.BastRequest
// @Failure 500 {object} map[string]interface{}
// @Router /bast-requests [get]
func (h *BastRequestHandler) GetAllRequests(c *gin.Context) {
	customerID := c.Query("customer_id")
	projectID := c.Query("project_id")
	status := c.Query("status")

	requests, err := h.service.GetAllRequests(customerID, projectID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, requests)
}

// GetRequestByID godoc
// @Summary Get a BAST request by ID
// @Description Retrieve a single BAST request based on its ID
// @Tags bast-requests
// @Produce json
// @Param id path string true "Request ID"
// @Success 200 {object} models.BastRequest
// @Failure 404 {object} map[string]interface{}
// @Router /bast-requests/{id} [get]
func (h *BastRequestHandler) GetRequestByID(c *gin.Context) {
	id := c.Param("id")
	request, err := h.service.GetRequestByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}
	c.JSON(http.StatusOK, request)
}

// CreateRequest godoc
// @Summary Create a new BAST request
// @Description Submit a new BAST request and automatically generate running number or use PO number
// @Tags bast-requests
// @Accept json
// @Produce json
// @Param request body models.BastRequest true "Request Data"
// @Success 201 {object} models.BastRequest
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /bast-requests [post]
func (h *BastRequestHandler) CreateRequest(c *gin.Context) {
	var input models.BastRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateRequest(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, input)
}

type UpdateStatusInput struct {
	Status string `json:"status" binding:"required"`
}

// UpdateStatus godoc
// @Summary Update request status
// @Description Update the status of a specific BAST request (Active, Used, Void)
// @Tags bast-requests
// @Accept json
// @Produce json
// @Param id path string true "Request ID"
// @Param status body UpdateStatusInput true "Status Data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /bast-requests/{id}/status [patch]
func (h *BastRequestHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var input UpdateStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateStatus(id, input.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}

// GetRequestAudit godoc
// @Summary Get audit log for a request
// @Description Retrieve the history of changes for a specific request ID
// @Tags bast-requests
// @Produce json
// @Param id path string true "Request ID"
// @Success 200 {object} map[string]interface{}
// @Router /bast-requests/{id}/audit [get]
func (h *BastRequestHandler) GetRequestAudit(c *gin.Context) {
	id := c.Param("id")
	// TODO: Call audit service to get logs for this specific request ID
	c.JSON(http.StatusOK, gin.H{"message": "Audit logs for request " + id})
}
