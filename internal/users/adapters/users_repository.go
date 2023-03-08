package adapters

type UsersRepository interface {
	FindUserByISU(ISU int) (user UserDTO, err error)
	FindUserByPhoneNumber(phoneNumber string) (user UserDTO, err error)
	InsertUser(user UserDTO) error
	UpdateUser(user UserDTO) error
	GetUsers(offset, limit int) (users []UserDTO, err error)
}

type UserDTO struct {
	ISU         int
	GivenName   string
	MiddleName  string
	FamilyName  string
	Email       string
	PhoneNumber string
}
