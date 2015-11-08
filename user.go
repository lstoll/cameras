package main

import "github.com/martini-contrib/sessionauth"

type User struct {
	Passcode      string `form:"passcode"`
	authenticated bool   `form:"-"`
}

func GenerateAnonymousUser() sessionauth.User {
	return &User{}
}

func (u *User) Login() {
	u.authenticated = true
}

func (u *User) Logout() {
	u.authenticated = false
}

func (u *User) IsAuthenticated() bool {
	return u.authenticated
}

func (u *User) UniqueId() interface{} {
	return u.Passcode
}

func (u *User) GetById(id interface{}) error {
	u.Passcode = id.(string)
	return nil
}
