package services

import (
	"bast-request/internal/models"
	"bast-request/internal/repositories"
)

type CustomerService struct {
	repo *repositories.CustomerRepository
}

func NewCustomerService(repo *repositories.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) GetAllCustomers(status string, name string) ([]models.Customer, error) {
	return s.repo.FindAll(status, name)
}

func (s *CustomerService) GetCustomerByID(id string) (models.Customer, error) {
	return s.repo.FindByID(id)
}

func (s *CustomerService) CreateCustomer(customer *models.Customer) error {
	return s.repo.Create(customer)
}

func (s *CustomerService) UpdateCustomer(id string, input *models.Customer) (models.Customer, error) {
	customer, err := s.repo.FindByID(id)
	if err != nil {
		return customer, err
	}

	customer.CustomerName = input.CustomerName
	customer.CustomerCode = input.CustomerCode
	customer.Status = input.Status

	err = s.repo.Update(&customer)
	return customer, err
}

func (s *CustomerService) DeleteCustomer(id string) error {
	return s.repo.Delete(id)
}
