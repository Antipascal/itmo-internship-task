package adapters

import "gorm.io/gorm"

type UsersPostgres struct {
	db *gorm.DB
}

type userInfo struct {
	ISU         int    `gorm:"primaryKey"`
	GivenName   string `gorm:"not null"`
	MiddleName  string `gorm:"not null"`
	FamilyName  string
	Email       string
	PhoneNumber string
}

func (userInfo) TableName() string {
	return "user_info"
}

func (u userInfo) ToDTO() UserDTO {
	return UserDTO{
		ISU:         u.ISU,
		GivenName:   u.GivenName,
		MiddleName:  u.MiddleName,
		FamilyName:  u.FamilyName,
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
	}
}

// NewUsersPostgres creates new UsersPostgres instance
// and migrates "user_info" table structure
func NewUsersPostgres(db *gorm.DB) (*UsersPostgres, error) {
	if err := db.AutoMigrate(&userInfo{}); err != nil {
		return nil, err
	}

	return &UsersPostgres{db: db}, nil
}

func (u *UsersPostgres) FindByISU(ISU int) (UserDTO, error) {
	var userInfo userInfo
	if err := u.db.First(&userInfo, "isu = ?", ISU).Error; err != nil {
		return UserDTO{}, err
	}
	return userInfo.ToDTO(), nil
}

func (u *UsersPostgres) FindByPhoneNumber(phoneNumber string) (UserDTO, error) {
	var userInfo userInfo
	if err := u.db.First(&userInfo, "phone_number = ?", phoneNumber).Error; err != nil {
		return UserDTO{}, err
	}
	return userInfo.ToDTO(), nil
}

func (u *UsersPostgres) Insert(user UserDTO) error {
	userInfo := userInfo{
		ISU:         user.ISU,
		GivenName:   user.GivenName,
		MiddleName:  user.MiddleName,
		FamilyName:  user.FamilyName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
	}
	return u.db.Create(&userInfo).Error
}

func (u *UsersPostgres) Update(user UserDTO) error {
	userInfo := userInfo{
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
	var users []userInfo
	if err := u.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}

	var usersDTO []UserDTO
	for _, user := range users {
		usersDTO = append(usersDTO, user.ToDTO())
	}
	return usersDTO, nil
}
