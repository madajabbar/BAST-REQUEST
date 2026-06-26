package handlers

import (
	"bast-request/internal/models"
	"bast-request/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BastFormatHandler struct {
	service *services.BastFormatService
}

func NewBastFormatHandler(service *services.BastFormatService) *BastFormatHandler {
	return &BastFormatHandler{service: service}
}

// GetAllFormats godoc
// @Summary Get all BAST formats
// @Description Retrieve a list of all BAST formats
// @Tags bast-formats
// @Produce json
// @Success 200 {array} models.BastFormat
// @Failure 500 {object} map[string]interface{}
// @Router /bast-formats [get]
func (h *BastFormatHandler) GetAllFormats(c *gin.Context) {
	formats, err := h.service.GetAllFormats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, formats)
}

// GetFormatByID godoc
// @Summary Get a BAST format by ID
// @Description Retrieve a single BAST format based on its ID
// @Tags bast-formats
// @Produce json
// @Param id path string true "Format ID"
// @Success 200 {object} models.BastFormat
// @Failure 404 {object} map[string]interface{}
// @Router /bast-formats/{id} [get]
func (h *BastFormatHandler) GetFormatByID(c *gin.Context) {
	id := c.Param("id")
	format, err := h.service.GetFormatByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Format not found"})
		return
	}
	c.JSON(http.StatusOK, format)
}

// CreateFormat godoc
// @Summary Create a new BAST format
// @Description Add a new BAST format to the master data
// @Tags bast-formats
// @Accept json
// @Produce json
// @Param format body models.BastFormat true "Format Data"
// @Success 201 {object} models.BastFormat
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /bast-formats [post]
func (h *BastFormatHandler) CreateFormat(c *gin.Context) {
	var input models.BastFormat
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateFormat(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, input)
}

// UpdateFormat godoc
// @Summary Update an existing BAST format
// @Description Update data for a specific BAST format
// @Tags bast-formats
// @Accept json
// @Produce json
// @Param id path string true "Format ID"
// @Param format body models.BastFormat true "Format Data"
// @Success 200 {object} models.BastFormat
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /bast-formats/{id} [put]
func (h *BastFormatHandler) UpdateFormat(c *gin.Context) {
	id := c.Param("id")
	var input models.BastFormat
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	format, err := h.service.UpdateFormat(id, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, format)
}

// DeleteFormat godoc
// @Summary Delete a BAST format
// @Description Nonaktifkan BAST format based on ID
// @Tags bast-formats
// @Produce json
// @Param id path string true "Format ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /bast-formats/{id} [delete]
func (h *BastFormatHandler) DeleteFormat(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteFormat(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Format deleted successfully"})
}
