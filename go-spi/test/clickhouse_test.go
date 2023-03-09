package test

import (
	"go-spi/db"
	"testing"
)

// func TestConnectClickhouse(t *testing.T) {
// 	dbPtr := db.GetDB()
// 	err := dbPtr.Ping()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	dbPtr.Close()
// }

func TestInsertData(t *testing.T) {
	db.InsertData()
}

// func TestGeoHash(t *testing.T) {
// 	hash := geohash.Encode(0.000000000, -0.000000000)
// 	fmt.Println("打印hash: ", hash)
// }
