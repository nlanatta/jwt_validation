package main

import "sync"

type userDB struct {
	users map[string]*User
}

var o1 sync.Once
var singleton *userDB

func UserDB() *userDB {
	o1.Do(func() {
		singleton = initializeUserDaoService()
	})
	return singleton
}

func initializeUserDaoService() *userDB {
	udb := &userDB{
		users: make(map[string]*User),
	}
	return udb
}

func (udb *userDB) AddUser(user *User) {
	udb.users[user.Email] = user
}

func (udb *userDB) GetUser(email string) *User {
	u, ok := udb.users[email]
	if ok {
		return u
	}
	return nil
}

func (udb *userDB) GetUserById(uid int64) *User {
	for _, v := range udb.users {
		if v.ID == uid {
			return v
		}
	}
	return nil
}
