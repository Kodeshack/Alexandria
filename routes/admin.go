package routes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"alexandria.app/models"
	"alexandria.app/view"
)

func AdminRoutes(r *mux.Router, templateDir string, userStorage models.UserStorage, sessionStorage *models.SessionStorage) {
	r.HandleFunc("/api/user", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.Header.Get("content-type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(NewErrorJSON(http.StatusBadRequest, `Content-Type must be "application/json"`))
			return
		}

		data, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(NewErrorJSON(http.StatusBadRequest, "Error while reading request body"))
			return
		}

		u := models.CreateUser{}
		err = json.Unmarshal(data, &u)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(NewErrorJSON(http.StatusBadRequest, "Invalid JSON"))
			return
		}

		user, err := models.NewUser(u.Email, u.DisplayName, u.Password, u.Admin)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = userStorage.AddUser(user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(NewErrorJSON(http.StatusBadRequest, "User already exists"))
			return
		}

		err = userStorage.Save()
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resJSON, err := user.JSON()
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(resJSON)
	})

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			v := view.New("layout", "login", templateDir, nil)

			if err := v.Render(w); err != nil {
				log.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		if len(email) == 0 || len(password) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			v := view.New("layout", "login", templateDir, nil)

			if err := v.Render(w); err != nil {
				log.Print(err)
				return
			}

			return
		}

		user := userStorage.CheckUserLogin(email, password)
		if user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			v := view.New("layout", "login", templateDir, nil)

			if err := v.Render(w); err != nil {
				log.Print(err)
				return
			}
			return
		}

		session := models.NewSession(user)

		sessionStorage.AddSession(session)

		http.SetCookie(w, session.Cookie(false))

		http.Redirect(w, r, "/", http.StatusFound)
	})

	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		if session := models.GetRequestSession(r); session != nil {
			sessionStorage.RemoveSession(session)
			http.SetCookie(w, models.RemoveSessionCookie(false))
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}).Methods(http.MethodGet)
}
