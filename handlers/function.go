package handlers

import (
	"api-mux/structs"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func EncriptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	pwd := string(hashedPassword)
	return pwd, err
}
func CekPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ReturnCheckError(w http.ResponseWriter, err error) {
	if err != nil {
		res := structs.Result{Code: http.StatusInternalServerError, Data: nil, Message: err.Error()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		result, _ := json.Marshal(res)
		w.Write(result)
	}
}

func ReturnResult(w http.ResponseWriter, result []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
