package domain

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type StarDao struct {
	db *gorm.DB
}

func NewStarDao(db *gorm.DB) *StarDao {
	return &StarDao{db: db}
}

func (sd *StarDao) GetNewStarsForRepository(respositoryId string) ([]Star, error) {
	models := []Star{}
	yesterday := time.Now().Add(-(time.Hour * 24))
	return models, sd.db.Where("stars.repository_id = ?", respositoryId).
		Where("stars.starred_at > ?", yesterday).
		Find(&models).Error
}

func (sd *StarDao) GetNewUnStarsForRepository(respositoryId string) ([]Star, error) {
	models := []Star{}
	yesterday := time.Now().Add(-(time.Hour * 24))
	return models, sd.db.Where("stars.repository_id = ?", respositoryId).
		Where("stars.updated_at > ?", yesterday).
		Find(&models).Error
}

func (sd *StarDao) GetValidForRepository(repositoryId string) ([]Star, error) {
	models := []Star{}
	return models, sd.db.Where("stars.repository_id = ?", repositoryId).
		Where("stars.valid = 1").
		Find(&models).Error
}

func (sd *StarDao) GetInValidForRepository(repositoryId string) ([]Star, error) {
	models := []Star{}
	return models, sd.db.Where("stars.repository_id = ?", repositoryId).
		Where("stars.valid = 0").
		Find(&models).Error
}

func (sd *StarDao) GetForRepository(repositoryId string) ([]Star, error) {
	models := []Star{}
	return models, sd.db.Where("stars.repository_id = ?", repositoryId).
		Find(&models).Error
}

func (sd *StarDao) Create(model *Star) error {
	model.Id = uuid.NewV4().String()
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()
	return sd.db.Create(model).Error
}

func (sd *StarDao) CreateIfNotExists(model *Star) (bool, error) {
	old := Star{}
	err := sd.db.Where("stars.repository_id = ?", model.RepositoryId).
		Where("stars.user_id = ?", model.UserId).
		First(&old).Error
	if err != nil {
		return true, sd.Create(model)
	} else {
		// if old was unstar, re star id
		if old.Valid == 0 {
			old.Valid = 1
			sd.Update(&old)
			return true, nil
		}
	}
	return false, nil
}

func (sd *StarDao) Update(model *Star) error {
	model.UpdatedAt = time.Now()
	return sd.db.Save(model).Error
}

// get the list of stars which is not associated to a id in the list
func (sd *StarDao) NotOneOfUsers(repositoryId string, usersIds []string) ([]Star, error) {
	models := []Star{}
	return models, sd.db.Where("stars.repository_id = ?", repositoryId).
		Where("stars.user_id not in (?)", usersIds).
		Where("stars.valid = 1").
		Find(&models).Error
}

func (sd *StarDao) SetInvalids(stars []Star) error {
	for _, m := range stars {
		m.Valid = 0
		sd.Update(&m)
	}
	return nil
}
