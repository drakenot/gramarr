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
	// UANone const
	UANone UserAccess = iota
	// UARevoked const
	UARevoked
	// UAMember const
	UAMember
	// UAAdmin const
	UAAdmin
)

// User struct
type User struct {
	ID        int        `json:"id"`
	Username  string     `json:"username"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Access    UserAccess `json:"access"`
}

// IsAdmin func
func (u User) IsAdmin() bool {
	return u.Access == UAAdmin
}

// IsMember func
func (u User) IsMember() bool {
	return u.Access == UAMember
}

// IsRevoked func
func (u User) IsRevoked() bool {
	return u.Access == UARevoked
}

// DisplayName func
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

// Recipient func
func (u User) Recipient() string {
	return strconv.Itoa(u.ID)
}

// UserDB struct
type UserDB struct {
	users    []User
	usersMap map[int]User
	dbPath   string
}

// NewUserDB func
func NewUserDB(dbPath string) (db *UserDB, err error) {
	db = &UserDB{
		users:    []User{},
		usersMap: map[int]User{},
		dbPath:   dbPath,
	}

	loadErr := db.Load()
	if !os.IsNotExist(loadErr) {
		err = loadErr
	}

	return
}

// Create func
func (u *UserDB) Create(user User) error {
	if u.Exists(user.ID) {
		return fmt.Errorf("user with ID %d already exists", user.ID)
	}

	u.users = append(u.users, user)
	u.usersMap[user.ID] = user

	u.Save()
	return nil
}

// Update func
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

// Delete func
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

// User func
func (u *UserDB) User(id int) (User, bool) {
	user, exists := u.usersMap[id]
	return user, exists
}

// Exists func
func (u *UserDB) Exists(id int) bool {
	_, ok := u.usersMap[id]
	return ok
}

// Users func
func (u *UserDB) Users() []User {
	return u.users
}

// Admins func
func (u *UserDB) Admins() []User {
	var result []User
	for _, user := range u.users {
		if user.Access == UAAdmin {
			result = append(result, user)
		}
	}
	return result
}

// Members func
func (u *UserDB) Members() []User {
	var result []User
	for _, user := range u.users {
		if user.Access == UAMember {
			result = append(result, user)
		}
	}
	return result
}

// Revoked func
func (u *UserDB) Revoked() []User {
	var result []User
	for _, user := range u.users {
		if user.Access == UARevoked {
			result = append(result, user)
		}
	}
	return result
}

// IsAdmin func
func (u *UserDB) IsAdmin(id int) bool {
	user, ok := u.usersMap[id]
	if !ok {
		return false
	}
	return user.Access == UAAdmin
}

// IsMember func
func (u *UserDB) IsMember(id int) bool {
	user, ok := u.usersMap[id]
	if !ok {
		return false
	}
	return user.Access == UAMember
}

// IsRevoked func
func (u *UserDB) IsRevoked(id int) bool {
	user, ok := u.usersMap[id]
	if !ok {
		return false
	}
	return user.Access == UARevoked
}

// Save func
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

// Load func
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
