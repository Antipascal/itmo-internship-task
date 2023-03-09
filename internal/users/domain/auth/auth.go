package auth

import (
	"context"
	"errors"
	"github.com/Antipascal/itmo-internship-task/users/adapters"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"os"
)

type Manager struct {
	authRepository  adapters.AuthRepository
	usersRepository adapters.UsersRepository
	authConfig      *oauth2.Config
	oicdProvider    *oidc.Provider
}

type itmoIdUserInfo struct {
	ISU int `json:"isu"`
}

const (
	itmoProviderURI = "https://id.itmo.ru/auth/realms/itmo"
)

func NewManager(authRepository adapters.AuthRepository, usersRepository adapters.UsersRepository) (manager *Manager, err error) {

	provider, err := oidc.NewProvider(context.Background(), itmoProviderURI)
	if err != nil {
		return
	}

	authConfig := oauth2.Config{
		ClientID:     os.Getenv("ITMO_CLIENT_ID"),
		ClientSecret: os.Getenv("ITMO_CLIENT_SECRET"),
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://localhost:8080/",
		Scopes:       []string{oidc.ScopeOpenID},
	}

	return &Manager{authRepository, usersRepository, &authConfig, provider}, nil
}

func (m *Manager) GetUserAccessToken(code, state string) (AccessToken string, err error) {
	u, err := m.getUserInfo(code)
	if err != nil || u.ISU == 0 {
		return "", errors.New("can't get user info")
	}

	// Add new token record to db
	AccessToken = state
	if err = m.authRepository.Create(AccessToken, u.ISU); err != nil {
		return "", err
	}

	// Add user to db if not exists
	if _, err = m.usersRepository.FindByISU(u.ISU); err != nil {
		user := adapters.UserDTO{ISU: u.ISU}
		if m.usersRepository.Insert(user) != nil {
			return "", err
		}
	}

	return
}

func (m *Manager) getUserInfo(code string) (user itmoIdUserInfo, err error) {
	oauth2Token, err := m.authConfig.Exchange(context.Background(), code)
	if err != nil {
		return itmoIdUserInfo{}, err
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		// handle missing token
	}

	var verifier = m.oicdProvider.Verifier(&oidc.Config{ClientID: m.authConfig.ClientID})

	// Parse and verify ID Token payload.
	idToken, err := verifier.Verify(context.Background(), rawIDToken)
	if err != nil {
		return itmoIdUserInfo{}, err
	}

	if err := idToken.Claims(&user); err != nil {
		return itmoIdUserInfo{}, err
	}

	return user, nil
}

// GetAuthURL returns a URL to OAuth 2.0 provider's consent page that asks for permissions for the required scopes explicitly.
// State is a token to protect the user from CSRF attacks.
func (m *Manager) GetAuthURL(state string) string {
	return m.authConfig.AuthCodeURL(state)
}

// GetISU returns ISU of user by token and true if user is authorized
// otherwise returns 0 and false.
func (m *Manager) GetISU(token string) (ISU int, isAuthorized bool) {
	isu, err := m.authRepository.FindISU(token)
	if err != nil {
		return 0, false
	}
	return isu, true
}
