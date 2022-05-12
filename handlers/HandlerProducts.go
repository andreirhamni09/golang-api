package handlers

import (
	"api-mux/connections"
	"api-mux/structs"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body)

	var dbProducts structs.Products

	json.Unmarshal(payloads, &dbProducts)

	if err := connections.DB.Create(&dbProducts).Error; err != nil {
		ReturnCheckError(w, err)
	}

	res := structs.Result{Code: 200, Data: dbProducts, Message: "Berhasil Menambahkan Product Baru"}

	result, err := json.Marshal(res)
	ReturnCheckError(w, err)
	ReturnResult(w, result)
}
func GetProductsLimit(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if limit < 1 {
		limit = 10
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	if limit < 1 {
		limit = 0
	}

	dbProducts := []structs.Products{}

	if err := connections.DB.Limit(limit).Offset(offset).Find(&dbProducts).Error; err != nil {
		ReturnCheckError(w, err)
	}

	res := structs.Result{Code: 200, Data: dbProducts, Message: "Products has successfully retrieve"}
	resuts, err := json.Marshal(res)

	ReturnCheckError(w, err)
	ReturnResult(w, resuts)
}

func GetProductId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	dbProducts := structs.Products{}

	if err := connections.DB.First(&dbProducts, id).Error; err != nil {
		ReturnCheckError(w, err)
	}

	res := structs.Result{Code: 200, Data: dbProducts, Message: "Product Ditemukan"}

	result, err := json.Marshal(res)

	ReturnCheckError(w, err)
	ReturnResult(w, result)
}

func UpdateProductById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	payloads, _ := ioutil.ReadAll(r.Body)

	var dbProducts structs.Products

	connections.DB.First(&dbProducts, id)

	json.Unmarshal(payloads, &dbProducts)

	if err := connections.DB.Model(&dbProducts).Update(dbProducts).Error; err != nil {
		ReturnCheckError(w, err)
	}

	res := structs.Result{Code: 200, Data: dbProducts, Message: "Berhasil Update Data Product"}

	result, err := json.Marshal(res)
	ReturnCheckError(w, err)
	ReturnResult(w, result)
}

func DeleteProductById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var dbProducts structs.Products

	if err := connections.DB.First(&dbProducts, id).Error; err != nil {
		ReturnCheckError(w, err)
	}
	if err := connections.DB.Delete(&dbProducts).Error; err != nil {
		ReturnCheckError(w, err)
	}

	res := structs.Result{Code: 200, Data: dbProducts, Message: "Berhasil Menghapus Product"}

	result, err := json.Marshal(res)

	ReturnCheckError(w, err)
	ReturnResult(w, result)
}
