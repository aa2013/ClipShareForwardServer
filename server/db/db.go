package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

var AppDb *gorm.DB

func ConnectDB(dbPath string) {
	CloseDB()
	dir := filepath.Dir(dbPath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	AppDb = db
	if err != nil {
		panic("failed to connect database" + err.Error())
	}
	err = db.AutoMigrate(&PlanType{})
	if err != nil {
		panic("failed to migrate database" + err.Error())
		return
	}
	err = db.AutoMigrate(&PlanKey{})
	if err != nil {
		panic("failed to migrate database" + err.Error())
		return
	}
}
func CloseDB() {
	if AppDb == nil {
		return
	}

	db, _ := AppDb.DB()
	_ = db.Close()
}
func checkDb() {
	if AppDb == nil {
		panic("Database not initialized")
	}
}
