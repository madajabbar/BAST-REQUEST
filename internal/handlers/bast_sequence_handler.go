package handlers

import (
	_ "bast-request/internal/models"
	"bast-request/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BastSequenceHandler struct {
	service *services.BastSequenceService
}

func NewBastSequenceHandler(service *services.BastSequenceService) *BastSequenceHandler {
	return &BastSequenceHandler{service: service}
}

// GetSequence godoc
// @Summary Get running sequence
// @Description Retrieve the current running sequence for a specific BAST format, year, and month
// @Tags bast-sequences
// @Produce json
// @Param format_id query string true "Format ID"
// @Param year query int true "Year"
// @Param month query int true "Month"
// @Success 200 {object} models.BastSequence
// @Failure 404 {object} map[string]interface{}
// @Router /bast-sequences [get]
func (h *BastSequenceHandler) GetSequence(c *gin.Context) {
	formatID := c.Query("format_id")
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	year, _ := strconv.Atoi(yearStr)
	month, _ := strconv.Atoi(monthStr)

	seq, err := h.service.GetSequence(formatID, year, month)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sequence not found"})
		return
	}
	c.JSON(http.StatusOK, seq)
}

type ResetSequenceRequest struct {
	FormatID   string `json:"format_id" binding:"required"`
	Year       int    `json:"year" binding:"required"`
	Month      int    `json:"month" binding:"required"`
	LastNumber int    `json:"last_number" binding:"required"`
}

// ResetSequence godoc
// @Summary Reset running sequence
// @Description Override or reset the running sequence for a specific BAST format
// @Tags bast-sequences
// @Accept json
// @Produce json
// @Param sequence body ResetSequenceRequest true "Reset Sequence Data"
// @Success 200 {object} models.BastSequence
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /bast-sequences/reset [post]
func (h *BastSequenceHandler) ResetSequence(c *gin.Context) {
	var input ResetSequenceRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	seq, err := h.service.ResetSequence(input.FormatID, input.Year, input.Month, input.LastNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, seq)
}
