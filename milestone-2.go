package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/xuri/excelize/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	FirstName string `gorm:"type:varchar(100);not null"`
	LastName  string `gorm:"type:varchar(100);not null"`
	Age       int    `gorm:"not null"`
}

type DbConfig struct {
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBHost     string `json:"db_host"`
	DBName     string `json:"db_name"`
}

func main() {
	users := readUserDataFromExcel("users.xlsx")
	dbConfig := getDbCreds("db_config.json")
	db := configDatabase(dbConfig.DBUser, dbConfig.DBPassword, dbConfig.DBHost, dbConfig.DBName)
	save_data_from_struct_to_db(db, users)
}

func getDbCreds(filepath string) DbConfig {
	configFile, error := os.Open(filepath)
	if error != nil {
		log.Fatal("Error opening config.json:", error)
	}

	var dbConfig DbConfig
	if err := json.NewDecoder(configFile).Decode(&dbConfig); err != nil {
		log.Fatal("Error decoding config.json:", err)
	}
	defer configFile.Close()
	return dbConfig
}

func readUserDataFromExcel(filepath string) (users []User) {
	f, error := excelize.OpenFile(filepath)
	if error != nil {
		log.Fatal(error)
		return nil
	}
	rows, error := f.GetRows(f.GetSheetName(0))
	if error != nil {
		log.Fatal(error)
		return nil
	}
	for i := 1; i < len(rows); i++ {
		var age int
		fmt.Sscanf(rows[i][2], "%d", &age)
		user := User{
			FirstName: rows[i][0],
			LastName:  rows[i][1],
			Age:       age,
		}
		users = append(users, user)
	}
	return users
}

func configDatabase(user_name string, password string, host string, database_name string) *gorm.DB {
	db_uri := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user_name, password, host, database_name)
	db, error := gorm.Open(mysql.Open(db_uri), &gorm.Config{})
	if error != nil {
		log.Printf("Connection to database failed: %v", error)
		return nil
	}
	return db
}

func save_data_from_struct_to_db(db *gorm.DB, users []User) {
	error := db.AutoMigrate(&User{})
	if error != nil {
		log.Fatal("Migration of User table failed: ", error)
	}
	for _, user := range users {
		error := db.Create(&user).Error
		if error != nil {
			log.Printf("Failed to save the user data for the user: %s %s", user.FirstName, user.LastName)
		} else {
			log.Printf("User data saved: %s %s", user.FirstName, user.LastName)
		}
	}
}
