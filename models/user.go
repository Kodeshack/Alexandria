package models

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"

	"github.com/google/uuid"

	"alexandria.app/crypto"
)

// A User of the wiki. Can create/edit/delete all articles on the wiki.
type User struct {
	ID            uint32
	Admin         bool
	Email         string
	DisplayName   string
	CreationDate  int64
	Password      string
	Salt          string
	Argon2KeyLen  uint32
	Argon2Memory  uint32
	Argon2Threads uint8
	Argon2Time    uint32
	Argon2Version int
}

// Following https://tools.ietf.org/html/draft-irtf-cfrg-argon2-03#section-4
const (
	argon2KeyLen  uint32 = 32
	argon2Memory  uint32 = 1024 * 1024 // 1GiB
	argon2Threads uint8  = 4
	argon2Time    uint32 = 1
	argon2Version        = 0x13
)

// NewUser will create a new User struct and hash the password.
// Will always use the current hashing parameters (which may change in future).
func NewUser(email, displayName, password string, admin bool) (*User, error) {
	salt, err := crypto.GetRandomString(16)
	if err != nil {
		return nil, err
	}

	tempPasswd := argon2.IDKey([]byte(password), []byte(salt), argon2Time, argon2Memory, argon2Threads, argon2KeyLen)

	return &User{
		Admin:         admin,
		Email:         email,
		DisplayName:   displayName,
		Password:      fmt.Sprintf("%x", tempPasswd),
		Salt:          salt,
		CreationDate:  time.Now().Unix(),
		Argon2KeyLen:  argon2KeyLen,
		Argon2Memory:  argon2Memory,
		Argon2Threads: argon2Threads,
		Argon2Time:    argon2Time,
		Argon2Version: argon2Version,
	}, nil
}

// UserStorage is a presistent database of all users in the system.
type UserStorage interface {
	GetUsers() []*User
	AddUser(*User) error
	DeleteUser(id uint32)
	Save() error
	CheckUserLogin(email, password string) *User
	CheckUserPassword(user *User, password string) *User
	SetUserPassword(user *User, newPassword string) error
	IsEmpty() bool
}

type userStorage struct {
	Version int
	path    string
	Users   []*User
}

// GetUsers returns a slice of all users.
func (udb *userStorage) GetUsers() []*User {
	return udb.Users
}

// GetUser retrieves a user by their email.
// Performs a simple linear search.
func (udb *userStorage) GetUser(email string) *User {
	for _, u := range udb.Users {
		if strings.EqualFold(u.Email, email) {
			return u
		}
	}

	return nil
}

// GetUserByID retrieves a user by their id.
// Performs a simple linear search.
func (udb *userStorage) GetUserByID(id uint32) *User {
	for _, u := range udb.Users {
		if u.ID == id {
			return u
		}
	}

	return nil
}

// AddUser inserts a new user into the database.
// Note: This doesn't save the database to the file system!
func (udb *userStorage) AddUser(newUser *User) error {
	if udb.GetUser(newUser.Email) != nil {
		return errors.New("User Already Exists")
	}

	id := uuid.New().ID()
	if udb.GetUserByID(id) != nil {
		return errors.New("UUID Collision")
	}

	newUser.ID = id
	udb.Users = append(udb.Users, newUser)
	return nil
}

// AddUser deletes a user from the database.
// Note: This doesn't save the database to the file system!
func (udb *userStorage) DeleteUser(id uint32) {
	for i, u := range udb.Users {
		if u.ID == id {
			udb.Users = append(udb.Users[:i], udb.Users[i+1:]...)
			return
		}
	}
}

// CheckUserLogin checks if the provided email and password match a record in the user database.
// If the user can't be found or the passwords don't match this will return nil.
func (udb *userStorage) CheckUserLogin(email, password string) *User {
	user := udb.GetUser(email)
	return udb.CheckUserPassword(user, password)
}

// CheckUserPassword hashes provided password and checks if it matches the user's.
// If the user is nil or the password doesn't match this will return nil.
func (udb *userStorage) CheckUserPassword(user *User, password string) *User {
	if user == nil {
		return nil
	}

	tempPasswd := argon2.IDKey([]byte(password), []byte(user.Salt), user.Argon2Time, user.Argon2Memory, user.Argon2Threads, user.Argon2KeyLen)

	if user.Password != fmt.Sprintf("%x", tempPasswd) {
		return nil
	}

	return user
}

// SetUserPassword hashes the password using the current hashing parameters (which may change in future)
// and sets the fields on the user struct accordingly.
// Note: This doesn't save the database to the file system!
func (udb *userStorage) SetUserPassword(user *User, newPassword string) error {
	salt, err := crypto.GetRandomString(16)
	if err != nil {
		return err
	}

	tempPasswd := argon2.IDKey([]byte(newPassword), []byte(salt), argon2Time, argon2Memory, argon2Threads, argon2KeyLen)

	user.Password = fmt.Sprintf("%x", tempPasswd)
	user.Salt = salt
	user.Argon2KeyLen = argon2KeyLen
	user.Argon2Memory = argon2Memory
	user.Argon2Threads = argon2Threads
	user.Argon2Time = argon2Time
	user.Argon2Version = argon2Version

	return nil
}

// Save will encode the database and save it to the file system.
func (udb *userStorage) Save() error {
	gob.Register(User{})
	gob.Register(userStorage{})

	file, err := os.OpenFile(udb.path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := gob.NewEncoder(file)

	return enc.Encode(udb)
}

// IsEmpty checks if the user database is empty. This should only be the case when initially setting up the system.
func (udb *userStorage) IsEmpty() bool {
	return len(udb.Users) == 0
}

// LoadUserStorage loads the database from the file system or creates a new database if it can't find one at the provided path.
// Note: This doesn't save the database to the file system when creating it!
func LoadUserStorage(path string) (UserStorage, error) {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return &userStorage{
			Version: 1,
			path:    path,
			Users:   []*User{},
		}, nil
	} else if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	gob.Register(User{})
	gob.Register(userStorage{})

	dec := gob.NewDecoder(file)

	udb := userStorage{path: path}
	err = dec.Decode(&udb)
	if err != nil {
		return nil, err
	}

	return &udb, nil
}
