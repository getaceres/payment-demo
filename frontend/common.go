package frontend

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/getaceres/payment-demo/persistence"
)

func ReadBody(reader io.Reader, result interface{}) error {
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(result); err != nil {
		return fmt.Errorf("Invalid payload: %s", err.Error())
	}
	return nil
}

func Respond(w http.ResponseWriter, code int, payload []byte, contentType string) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)
	w.Write(payload)
}

func RespondWithText(w http.ResponseWriter, code int, text string) {
	Respond(w, code, []byte(text), "application/text")
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	content, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	Respond(w, code, content, "application/json")
	return nil
}

func RespondWithError(w http.ResponseWriter, code int, err error) {
	RespondWithText(w, code, err.Error())
}

func GetPersistenceErrorCode(err error) int {
	code := http.StatusInternalServerError
	switch err.(type) {
	case persistence.NotFoundError:
		code = http.StatusNotFound
	case persistence.AlreadyExistsError:
		code = http.StatusConflict
	}
	return code
}
