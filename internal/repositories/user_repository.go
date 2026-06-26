package repositories

import (
	"bast-request/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository{
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(email string) (models.User, error){
	var user models.User
	err := r.db.First(&user, "email = ?", email).Error
	return user, err
}

func (r *UserRepository) Create(user *models.User) error{
	return r.db.Create(user).Error
}

func (r *UserRepository) Update(user *models.User) error{
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id string) error{
	return r.db.Model(&models.User{}).Where("user_id = ?", id).Update("is_active", false).Error
}