package repositories

import (
	"bast-request/internal/models"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) FindAll(status string, name string) ([]models.Customer, error) {
	var customers []models.Customer
	query := r.db.Model(&models.Customer{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if name != "" {
		query = query.Where("customer_name ILIKE ?", "%"+name+"%")
	}

	err := query.Find(&customers).Error
	return customers, err
}

func (r *CustomerRepository) FindByID(id string) (models.Customer, error) {
	var customer models.Customer
	err := r.db.First(&customer, "customer_id = ?", id).Error
	return customer, err
}

func (r *CustomerRepository) Create(customer *models.Customer) error {
	return r.db.Create(customer).Error
}

func (r *CustomerRepository) Update(customer *models.Customer) error {
	return r.db.Save(customer).Error
}

func (r *CustomerRepository) Delete(id string) error {
	// Soft delete, or update status to inactive
	return r.db.Model(&models.Customer{}).Where("customer_id = ?", id).Update("status", "inactive").Error
}
