package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var dbMutex sync.Mutex

type UserAccess int

const (
	UANone UserAccess = iota
	UARevoked
	UAMember
	UAAdmin
)

type User struct {
	ID        int        `json:"id"`
	Username  string     `json:"username"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Access    UserAccess `json:"access"`
}

func (u User) IsAdmin() bool {
	return u.Access == UAAdmin
}

func (u User) IsMember() bool {
	return u.Access == UAMember
}

func (u User) IsRevoked() bool {
	return u.Access == UARevoked
}

func (u User) DisplayName() string {
	if u.Username != "" {
		return u.Username
	} else {
		if u.LastName != "" {
			return u.FirstName + " " + u.LastName
		}
		return u.FirstName
	}
}

func (u User) Recipient() string {
	return strconv.Itoa(u.ID)
}

type UserDB struct {
	users    []User
	usersMap map[int]User
	dbPath   string
}

func NewUserDB(dbPath string) (db *UserDB, err error) {
	db = &UserDB{
		users:    []User{},
		usersMap: map[int]User{},
		dbPath:   dbPath,
	}

	err = db.Load()
	return
}

func (u *UserDB) Create(user User) error {
	if u.Exists(user.ID) {
		return fmt.Errorf("user with ID %d already exists", user.ID)
	}

	u.users = append(u.users, user)
	u.usersMap[user.ID] = user

	u.Save()
	return nil
}

func (u *UserDB) Update(user User) error {
	if !u.Exists(user.ID) {
		return fmt.Errorf("user with ID %d doesn't exist", user.ID)
	}

	for i := 0; i < len(u.users); i++ {
		if u.users[i].ID == user.ID {
			u.users[i] = user
			break
		}
	}

	u.usersMap[user.ID] = user
	u.Save()
	return nil
}

func (u *UserDB) Delete(user User) error {
	if !u.Exists(user.ID) {
		return fmt.Errorf("user with ID %d doesn't exist", user.ID)
	}

	for i, usr := range u.users {
		if user.ID == usr.ID {
			u.users = append(u.users[:i], u.users[i+1:]...)
			break
		}
	}

	delete(u.usersMap, user.ID)
	u.Save()
	return nil
}

func (u *UserDB) User(id int) (User, bool) {
	user, exists := u.usersMap[id]
	return user, exists
}

func (u *UserDB) Exists(id int) bool {
	_, ok := u.usersMap[id]
	return ok
}

func (u *UserDB) Users() []User {
	return u.users
}

func (u *UserDB) Admins() []User {
	var result []User
	for _, user := range u.users {
		if user.Access == UAAdmin {
			result = append(result, user)
		}
	}
	return result
}

func (u *UserDB) Members() []User {
	var result []User
	for _, user := range u.users {
		if user.Access == UAMember {
			result = append(result, user)
		}
	}
	return result
}

func (u *UserDB) Revoked() []User {
	var result []User
	for _, user := range u.users {
		if user.Access == UARevoked {
			result = append(result, user)
		}
	}
	return result
}

func (u *UserDB) IsAdmin(id int) bool {
	user, ok := u.usersMap[id]
	if !ok {
		return false
	}
	return user.Access == UAAdmin
}

func (u *UserDB) IsMember(id int) bool {
	user, ok := u.usersMap[id]
	if !ok {
		return false
	}
	return user.Access == UAMember
}

func (u *UserDB) IsRevoked(id int) bool {
	user, ok := u.usersMap[id]
	if !ok {
		return false
	}
	return user.Access == UARevoked
}

func (u *UserDB) Save() error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	// Open a temporary file to hold the new database
	var tempDb *os.File
	tempDb, err := ioutil.TempFile(filepath.Dir(u.dbPath), filepath.Base(u.dbPath))
	if err != nil {
		return err
	}

	var db = struct {
		Users []User `json:"users"`
	}{
		u.users,
	}

	// Write the data to the new file
	enc := json.NewEncoder(tempDb)
	err = enc.Encode(db)
	if err != nil {
		return err
	}

	// Close the file if we succeeded in opening one
	err = tempDb.Close()
	if err != nil {
		return err
	}

	// Rename the temporary database over the permanent database
	err = os.Rename(tempDb.Name(), u.dbPath)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserDB) Load() error {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	raw, err := ioutil.ReadFile(u.dbPath)
	if err != nil {
		return err
	}

	var db struct {
		Users []User
	}

	json.Unmarshal(raw, &db)
	u.users = db.Users
	u.usersMap = map[int]User{}
	for _, user := range db.Users {
		u.usersMap[user.ID] = user
	}

	return nil
}
