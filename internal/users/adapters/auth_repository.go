package adapters

type AuthRepository interface {
	FindISUByToken(token string) (ISU int, err error)
	AddAuthRecord(token string, ISU int) error
}
