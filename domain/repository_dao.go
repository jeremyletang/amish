package domain

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type RepositoryDao struct {
	db *gorm.DB
}

func NewRepositoryDao(db *gorm.DB) *RepositoryDao {
	return &RepositoryDao{db: db}
}

func (rd *RepositoryDao) List() ([]*Repository, error) {
	models := []*Repository{}
	return models, rd.db.Find(&models).Error
}

func (rd *RepositoryDao) Create(model *Repository) error {
	model.Id = uuid.NewV4().String()
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()
	return rd.db.Create(model).Error
}

func (rd *RepositoryDao) GetOrCreate(model *Repository) (*Repository, error) {
	r, err := rd.GetByOwnerAndName(model.Owner, model.Name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return model, rd.Create(model)
		}
		return nil, err
	}

	return r, nil
}

func (rd *RepositoryDao) GetById(id string) (*Repository, error) {
	model := &Repository{Id: id}
	return model, rd.db.First(model).Error
}

func (rd *RepositoryDao) GetByOwnerAndName(owner, name string) (*Repository, error) {
	model := &Repository{}
	return model, rd.db.Where("repositories.owner = ?", owner).
		Where("repositories.name = ?", name).
		First(model).Error
}
