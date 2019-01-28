package models

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"golang.org/x/crypto/argon2"
)

func randomInt(max *big.Int) (int, error) {
	rand, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}

	return int(rand.Int64()), nil
}

// GetRandomString generate random string by specify chars.
func getRandomString(n int) (string, error) {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	buffer := make([]byte, n)
	max := big.NewInt(int64(len(alphanum)))

	for i := 0; i < n; i++ {
		index, err := randomInt(max)
		if err != nil {
			return "", err
		}

		buffer[i] = alphanum[index]
	}

	return string(buffer), nil
}

type CreateUser struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
	Admin       bool   `json:"admin"`
}

type JSONUser struct {
	ID           int64  `json:"id"`
	Admin        bool   `json:"admin"`
	Email        string `json:"email"`
	DisplayName  string `json:"display_name"`
	CreationDate int64  `json:"creation_date"`
}

type User struct {
	ID            int64  `json:"id"`
	Admin         bool   `json:"admin"`
	Email         string `json:"email"`
	DisplayName   string `json:"display_name"`
	CreationDate  int64
	Password      string `json:"password"`
	Salt          string
	argon2KeyLen  uint32
	argon2Memory  uint32
	argon2Threads uint8
	argon2Time    uint32
	argon2Version int
}

// Following https://tools.ietf.org/html/draft-irtf-cfrg-argon2-03#section-4
var (
	argon2KeyLen  uint32 = 32
	argon2Memory  uint32 = 1024 * 1024 // 1GiB
	argon2Threads uint8  = 4
	argon2Time    uint32 = 1
	argon2Version        = 0x13
)

func NewUser(email, displayName, password string, admin bool) (*User, error) {
	salt, err := getRandomString(16)
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
		argon2KeyLen:  argon2KeyLen,
		argon2Memory:  argon2Memory,
		argon2Threads: argon2Threads,
		argon2Time:    argon2Time,
		argon2Version: argon2Version,
	}, nil
}

func (u *User) JSON() ([]byte, error) {
	return json.Marshal(JSONUser{u.ID, u.Admin, u.Email, u.DisplayName, u.CreationDate})
}
