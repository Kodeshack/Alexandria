package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/argon2"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

var ContentPrefix = "content/"

func resolveRealFilePath(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(ContentPrefix + path)
	if err != nil {
		return []byte{}, err
	}

	options := blackfriday.WithExtensions(blackfriday.CommonExtensions | blackfriday.HardLineBreak)

	output := blackfriday.Run(data, options)

	return output, nil
}

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

func NewUser(email, displayName, password string, admin bool) (User, error) {
	salt, err := getRandomString(16)
	if err != nil {
		return User{}, err
	}

	tempPasswd := argon2.IDKey([]byte(password), []byte(salt), argon2Time, argon2Memory, argon2Threads, argon2KeyLen)

	return User{
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

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/user", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		data, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		u := CreateUser{}
		err = json.Unmarshal(data, &u)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := NewUser(u.Email, u.DisplayName, u.Password, u.Admin)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Println(user)

		w.WriteHeader(http.StatusCreated)
	})

	r.HandleFunc("/{path:[\\w\\d_/-]+}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		path := vars["path"]
		if len(path) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		data, err := resolveRealFilePath(path)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
