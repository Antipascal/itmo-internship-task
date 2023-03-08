package adapters

import (
	"gorm.io/gorm"
	"log"
)

type AuthPostgres struct {
	db *gorm.DB
}

type UserAuth struct {
	AccessToken string `gorm:"primaryKey"`
	ISU         int
}

func (UserAuth) TableName() string {
	return "user_auth"
}

func NewAuthPostgres(db *gorm.DB) (*AuthPostgres, error) {
	if err := db.AutoMigrate(&UserAuth{}); err != nil {
		return nil, err
	}

	return &AuthPostgres{db: db}, nil
}

func (a *AuthPostgres) FindISUByToken(token string) (int, error) {
	var authModel UserAuth
	if err := a.db.First(&authModel, "access_token = ?", token).Error; err != nil {
		return 0, err
	}
	return authModel.ISU, nil
}

func (a *AuthPostgres) AddAuthRecord(token string, ISU int) error {
	authModel := UserAuth{
		AccessToken: token,
		ISU:         ISU,
	}
	log.Println(authModel)
	return a.db.Create(&authModel).Error
}
