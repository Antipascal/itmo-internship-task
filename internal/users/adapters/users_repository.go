package adapters

type UsersRepository interface {
	FindByISU(ISU int) (user UserDTO, err error)
	FindByPhoneNumber(phoneNumber string) (user UserDTO, err error)
	Insert(user UserDTO) error
	Update(user UserDTO) error
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
