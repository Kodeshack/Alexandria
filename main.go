package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	blackfriday "gopkg.in/russross/blackfriday.v2"

	"./models"
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

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/user", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("content-type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		data, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		u := models.CreateUser{}
		err = json.Unmarshal(data, &u)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := models.NewUser(u.Email, u.DisplayName, u.Password, u.Admin)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resJSON, err := user.JSON()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(resJSON)
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
