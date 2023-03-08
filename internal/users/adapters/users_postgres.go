package adapters

import "gorm.io/gorm"

type UsersPostgres struct {
	db *gorm.DB
}

type UserInfo struct {
	ISU         int    `gorm:"primaryKey"`
	GivenName   string `gorm:"not null"`
	MiddleName  string `gorm:"not null"`
	FamilyName  string
	Email       string
	PhoneNumber string
}

func (UserInfo) TableName() string {
	return "user_info"
}

func (u UserInfo) ToDTO() UserDTO {
	return UserDTO{
		ISU:         u.ISU,
		GivenName:   u.GivenName,
		MiddleName:  u.MiddleName,
		FamilyName:  u.FamilyName,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
	}
}

func NewUsersPostgres(db *gorm.DB) (*UsersPostgres, error) {
	if err := db.AutoMigrate(&UserInfo{}); err != nil {
		return nil, err
	}

	return &UsersPostgres{db: db}, nil
}

func (u *UsersPostgres) FindUserByISU(ISU int) (UserDTO, error) {
	var userInfo UserInfo
	if err := u.db.First(&userInfo, "isu = ?", ISU).Error; err != nil {
		return UserDTO{}, err
	}
	return userInfo.ToDTO(), nil
}

func (u *UsersPostgres) FindUserByPhoneNumber(phoneNumber string) (UserDTO, error) {
	var userInfo UserInfo
	if err := u.db.First(&userInfo, "phone_number = ?", phoneNumber).Error; err != nil {
		return UserDTO{}, err
	}
	return userInfo.ToDTO(), nil
}

func (u *UsersPostgres) InsertUser(user UserDTO) error {
	userInfo := UserInfo{
		ISU:         user.ISU,
		GivenName:   user.GivenName,
		MiddleName:  user.MiddleName,
		FamilyName:  user.FamilyName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
	}
	return u.db.Create(&userInfo).Error
}

func (u *UsersPostgres) UpdateUser(user UserDTO) error {
	userInfo := UserInfo{
		ISU:         user.ISU,
		GivenName:   user.GivenName,
		MiddleName:  user.MiddleName,
		FamilyName:  user.FamilyName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
	}
	return u.db.Save(&userInfo).Error
}

func (u *UsersPostgres) GetUsers(offset, limit int) ([]UserDTO, error) {
	var users []UserInfo
	if err := u.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}

	var usersDTO []UserDTO
	for _, user := range users {
		usersDTO = append(usersDTO, user.ToDTO())
	}
	return usersDTO, nil
}
