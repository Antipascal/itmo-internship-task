package adapters

type AuthRepository interface {
	FindISU(token string) (ISU int, err error)
	Create(token string, ISU int) error
}
