package routes

import (
	"encoding/json"
	"log"
)

type ErrorJSON struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewErrorJSON(code int, message string) []byte {
	errorJSON := ErrorJSON{
		Code:    code,
		Message: message,
	}

	data, err := json.Marshal(&errorJSON)
	if err != nil {
		log.Fatal(err)
	}

	return data
}
