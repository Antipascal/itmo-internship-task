package adapters

type UsersRepository interface {
	// FindByISU returns user with given ISU
	// If user is not found returns error.
	FindByISU(ISU int) (user UserDTO, err error)

	// FindByPhoneNumber returns user with given phone number
	// If user is not found returns error.
	FindByPhoneNumber(phoneNumber string) (user UserDTO, err error)

	// Insert adds new user to database
	// If user already exists returns error.
	Insert(user UserDTO) error

	// Update updates user info
	// If user is not found returns error.
	Update(user UserDTO) error

	// GetUsers returns users with offset and limit.
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
