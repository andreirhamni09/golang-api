package connections

import (
	"fmt"

	"api-mux/structs"

	"github.com/jinzhu/gorm"
)

var (
	DB  *gorm.DB
	Err error
)

func Connect() {
	DB, Err = gorm.Open("mysql", "root:@/db_nasabah?charset=utf8&parseTime=True&loc=Local")
	if Err != nil {
		fmt.Println("Gagal Koneksi", Err)
	} else {
		fmt.Println("Berhasil Koneksi")
	}
	DB.AutoMigrate(&structs.Users{})
	DB.AutoMigrate(&structs.Products{})
}
