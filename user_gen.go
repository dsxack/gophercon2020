package gophercon2020

import (
	"github.com/jinzhu/gorm"
)

type UserRepository struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{db: db}
}
func (r UserRepository) Get(id uint) (*User, error) {
	entity := new(User)
	err := r.db.Limit(1).Where("id = ?", id).Find(entity).Error
	return entity, err
}
func (r UserRepository) Create(entity *User) error {
	return r.db.Create(entity).Error
}
func (r UserRepository) Update(entity *User) error {
	return r.db.Model(entity).Update(entity).Error
}
func (r UserRepository) Delete(entity *User) error {
	return r.db.Delete(entity).Error
}
