package handlers

import (
	"api-mux/connections"
	"api-mux/structs"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Wilkommen!")
}

func CreateUsers(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body)

	var dbUsers structs.Users
	json.Unmarshal(payloads, &dbUsers)

	res := structs.Result{Code: 500, Data: dbUsers, Message: "Unknown Error"}

	switch dbUsers.Role {
	case "0":
		dbUsers.Role = "user"
	case "1":
		dbUsers.Role = "admin"
	default:
		dbUsers.Role = "invalid"
		res.Code = 400
		res.Message = "Invalid User Role"
	}

	if dbUsers.Role != "invalid" {
		genPass, err := EncriptPassword(dbUsers.Password)
		ReturnCheckError(w, err)
		dbUsers.Password = genPass

		if err := connections.DB.Create(&dbUsers).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		res.Data = dbUsers
		res.Code = 200
		res.Message = "Add new user successfully"
	}
	result, err := json.Marshal(res)
	ReturnCheckError(w, err)
	ReturnResult(w, result)
}

func GetUsersLimit(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if limit < 1 {
		limit = 10
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit < 1 {
		limit = 0
	}

	dbUsers := []structs.Users{}

	if err := connections.DB.Limit(limit).Offset(offset).Find(&dbUsers).Error; err != nil {
		ReturnCheckError(w, err)
	}

	res := structs.Result{Code: 200, Data: dbUsers, Message: "User has successfully retrieve"}
	result, err := json.Marshal(res)

	ReturnCheckError(w, err)
	ReturnResult(w, result)
}

func GetUserId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	dbUsers := structs.Users{}
	connections.DB.First(&dbUsers, id)

	res := structs.Result{Code: 200, Data: dbUsers, Message: "Users Ditemukan"}

	result, err := json.Marshal(res)

	ReturnCheckError(w, err)
	ReturnResult(w, result)
}

func UpdateUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	payloads, _ := ioutil.ReadAll(r.Body)

	var dbUsers structs.Users
	res := structs.Result{Code: 200, Data: dbUsers, Message: "Unknown Error"}

	connections.DB.First(&dbUsers, id)

	json.Unmarshal(payloads, &dbUsers)

	switch dbUsers.Role {
	case "0":
		dbUsers.Role = "user"
	case "1":
		dbUsers.Role = "admin"
	default:
		dbUsers.Role = "invalid"
		res.Code = 400
		res.Message = "Invalid User Role"
	}

	if dbUsers.Role != "invalid" {
		if err := connections.DB.Model(&dbUsers).Update(&dbUsers).Error; err != nil {
			ReturnCheckError(w, err)
		}
		if !dbUsers.Status {
			connections.DB.Model(&dbUsers).Updates(map[string]interface{}{"status": false})
		}
		res.Code = 200
		res.Data = dbUsers
		res.Message = "Update user data successfully"
	}

	result, err := json.Marshal(res)
	ReturnCheckError(w, err)
	ReturnResult(w, result)

}

func DeleteUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var dbUsers structs.Users

	connections.DB.First(&dbUsers, id)
	connections.DB.Delete(&dbUsers)

	res := structs.Result{Code: 200, Data: dbUsers, Message: "Berhasil Menghapus Users"}

	result, err := json.Marshal(res)
	ReturnCheckError(w, err)
	ReturnResult(w, result)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body)

	var dbUser structs.Users
	var userLogin structs.UsersLogin
	res := structs.Result{Code: 200, Data: userLogin, Message: "Gagal Login"}
	json.Unmarshal(payloads, &userLogin)
	connections.DB.Where("username = ?", &userLogin.Username).Find(&dbUser)

	if CekPassword(userLogin.Password, dbUser.Password) {
		res = structs.Result{Code: 200, Data: dbUser, Message: "Berhasil Login"}
	}

	result, err := json.Marshal(res)
	ReturnCheckError(w, err)
	ReturnResult(w, result)
}
