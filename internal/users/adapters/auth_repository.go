package adapters

type AuthRepository interface {
	// FindISU returns ISU of user with given token
	// If user is not logged returns error.
	FindISU(token string) (ISU int, err error)

	// Create add record [token, value] to database.
	Create(token string, ISU int) error
}
