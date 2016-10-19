package domain

import (
	"time"

	"github.com/jinzhu/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}

func (ud *UserDao) List() ([]*User, error) {
	models := []*User{}
	return models, ud.db.Find(&models).Error
}

func (ud *UserDao) Create(model *User) error {
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()
	return ud.db.Create(model).Error
}

func (ud *UserDao) Update(model *User) error {
	return ud.db.Save(model).Error
}

func (ud *UserDao) GetById(id string) (*User, error) {
	model := &User{Id: id}
	return model, ud.db.First(model).Error
}

func (ud *UserDao) GetOrCreate(model *User) (*User, error) {
	u, err := ud.GetById(model.Id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return model, ud.Create(model)
		}
		return nil, err
	}

	return u, nil
}

func (ud *UserDao) GetContentNotUpdated() ([]*User, error) {
	models := []*User{}
	return models, ud.db.Where("users.content_updated = 0").
		Find(&models).Error
}
