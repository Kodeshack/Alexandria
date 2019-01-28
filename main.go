package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/russross/blackfriday.v2"
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
