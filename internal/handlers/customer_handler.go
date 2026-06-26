package handlers

import (
	"bast-request/internal/models"
	"bast-request/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	service *services.CustomerService
}

func NewCustomerHandler(service *services.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

// GetAllCustomers godoc
// @Summary Get all customers
// @Description Retrieve a list of customers, optionally filtered by status and name
// @Tags customers
// @Produce json
// @Param status query string false "Status filter"
// @Param nama query string false "Name filter"
// @Success 200 {array} models.Customer
// @Failure 500 {object} map[string]interface{}
// @Router /customers [get]
func (h *CustomerHandler) GetAllCustomers(c *gin.Context) {
	status := c.Query("status")
	name := c.Query("nama") // mapping from query parameter 'nama' to name

	customers, err := h.service.GetAllCustomers(status, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customers)
}

// GetCustomerByID godoc
// @Summary Get a customer by ID
// @Description Retrieve a single customer based on their ID
// @Tags customers
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} models.Customer
// @Failure 404 {object} map[string]interface{}
// @Router /customers/{id} [get]
func (h *CustomerHandler) GetCustomerByID(c *gin.Context) {
	id := c.Param("id")
	customer, err := h.service.GetCustomerByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}
	c.JSON(http.StatusOK, customer)
}

// CreateCustomer godoc
// @Summary Create a new customer
// @Description Add a new customer to the master data
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body models.Customer true "Customer Data"
// @Success 201 {object} models.Customer
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /customers [post]
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var input models.Customer
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateCustomer(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, input)
}

// UpdateCustomer godoc
// @Summary Update an existing customer
// @Description Update data for a specific customer
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param customer body models.Customer true "Customer Data"
// @Success 200 {object} models.Customer
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /customers/{id} [put]
func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id := c.Param("id")
	var input models.Customer
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, err := h.service.UpdateCustomer(id, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, customer)
}

// DeleteCustomer godoc
// @Summary Delete a customer
// @Description Nonaktifkan customer based on ID (soft delete / status inactive)
// @Tags customers
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /customers/{id} [delete]
func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteCustomer(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}
