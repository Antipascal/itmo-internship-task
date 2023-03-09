package users

import (
	"github.com/Antipascal/itmo-internship-task/users/adapters"
)

// Manager is a service that manages users.
// It is responsible for creating/updating/searching user info and storing it in database.
type Manager struct {
	authRepository  adapters.AuthRepository
	usersRepository adapters.UsersRepository
}

// NewManager creates new users Manager.
func NewManager(authRepository adapters.AuthRepository, usersRepository adapters.UsersRepository) *Manager {
	return &Manager{authRepository, usersRepository}
}

// GetUser returns user info by ISU.
// If user is not found returns error.
func (m *Manager) GetUser(ISU int) (user adapters.UserDTO, err error) {
	return m.usersRepository.FindByISU(ISU)
}

func (m *Manager) UpdateUser(user adapters.UserDTO) (err error) {
	return m.usersRepository.Update(user)
}
