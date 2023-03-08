package auth

import (
	"encoding/json"
	"errors"
	"github.com/Antipascal/itmo-internship-task/users/adapters"
	"net/http"
	"net/url"
	"os"
)

type Manager struct {
	authRepository  adapters.AuthRepository
	usersRepository adapters.UsersRepository
}

type itmoIdToken struct {
	AccessToken string `json:"access_token"`
}

type itmoIdUserInfo struct {
	ISU        int    `json:"isu"`
	GivenName  string `json:"given_name"`
	MiddleName string `json:"middle_name"`
	FamilyName string `json:"family_name"`
	Email      string `json:"email"`
}

// ItmoAuthLink TODO move to config
var (
	ItmoAuthLink = "https://id.itmo.ru/auth/realms/itmo/protocol/openid-connect/auth?client_id=test-dev&client_secret=" +
		os.Getenv("ITMO_ID_SECRET") +
		"&response_type=code&scope=openid&redirect_uri=http://localhost:8080/&state=itmoid"
	ItmoTokenLink = "id.itmo.ru/auth/realms/itmo/protocol/openid-connect/token"
	ItmoUserLink  = "id.itmo.ru/auth/realms/itmo/protocol/openid-connect/userinfo"
)

func NewManager(authRepository adapters.AuthRepository, usersRepository adapters.UsersRepository) *Manager {
	return &Manager{authRepository, usersRepository}
}

func (a *Manager) AuthItmoCode(code string) (AccessToken string, err error) {
	t, err := getItmoIdToken(code)
	if err != nil || t.AccessToken == "" {
		return "", errors.New("can't get access token")
	}
	u, err := getItmoIdUserInfo(t.AccessToken)
	if err != nil || u.ISU == 0 {
		return "", errors.New("can't get user info")
	}

	if err = a.authRepository.Create(t.AccessToken, u.ISU); err != nil {
		return "", err
	}

	if _, err = a.usersRepository.FindByISU(u.ISU); err != nil {
		user := adapters.UserDTO{
			ISU:        u.ISU,
			GivenName:  u.GivenName,
			MiddleName: u.MiddleName,
			FamilyName: u.FamilyName,
			Email:      u.Email,
		}

		if a.usersRepository.Insert(user) != nil {
			return "", err
		}
	}

	return t.AccessToken, nil
}

func getItmoIdToken(code string) (token itmoIdToken, err error) {
	u := url.URL{Path: ItmoTokenLink, Scheme: "https"}
	q := u.Query()
	q.Add("client_id", "test-dev")
	q.Add("client_secret", os.Getenv("ITMO_ID_SECRET"))
	q.Add("grant_type", "authorization_code")
	q.Add("redirect_uri", "http://localhost:8080/")
	q.Add("code", code)
	u.RawQuery = q.Encode()
	r, err := http.NewRequest("POST", u.String(), nil)
	resp, err := http.DefaultClient.Do(r)
	j := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	err = j.Decode(&token)
	if err != nil {
		return itmoIdToken{}, err
	}
	return token, nil
}

func getItmoIdUserInfo(token string) (User itmoIdUserInfo, err error) {
	u := url.URL{Path: ItmoUserLink, Scheme: "https"}
	r, err := http.NewRequest("GET", u.String(), nil)
	r.Header.Add("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(r)

	j := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	if err = j.Decode(&User); err != nil {
		return itmoIdUserInfo{}, err
	}
	return User, nil
}
