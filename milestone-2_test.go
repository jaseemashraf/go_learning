package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestReadUserDataFromExcel_ValidExcelFile(t *testing.T) {
	f := excelize.NewFile()
	f.SetCellValue("Sheet1", "A1", "FirstName")
	f.SetCellValue("Sheet1", "B1", "LastName")
	f.SetCellValue("Sheet1", "C1", "Age")
	f.SetCellValue("Sheet1", "A2", "Alice")
	f.SetCellValue("Sheet1", "B2", "Smith")
	f.SetCellValue("Sheet1", "C2", "25")
	f.SetCellValue("Sheet1", "A3", "Bob")
	f.SetCellValue("Sheet1", "B3", "Jones")
	f.SetCellValue("Sheet1", "C3", "30")
	filePath := filepath.Join(t.TempDir(), "test_users.xlsx")
	if error := f.SaveAs(filePath); error != nil {
		t.Fatalf("Failed to save test Excel file: %v", error)
	}
	users := readUserDataFromExcel(filePath)
	want := []User{
		{FirstName: "Alice", LastName: "Smith", Age: 25},
		{FirstName: "Bob", LastName: "Jones", Age: 30},
	}
	if !reflect.DeepEqual(users, want) {
		t.Errorf("readUserDataFromExcelWithPath() = %v, want %v", users, want)
	}
}

func TestGetDbCreds(t *testing.T) {
	configContent := `{
		"db_user": "test_user",
		"db_password": "test_learner",
		"db_host": "127.0.0.1:3306",
		"db_name": "test_learning"
	}`
	filePath := filepath.Join(t.TempDir(), "test_db_config.json")
	if err := os.WriteFile(filePath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test JSON file: %v", err)
	}
	dbConfig := getDbCreds(filePath)
	want := DbConfig{
		DBUser:     "test_user",
		DBPassword: "test_learner",
		DBHost:     "127.0.0.1:3306",
		DBName:     "test_learning",
	}
	if !reflect.DeepEqual(dbConfig, want) {
		t.Errorf("getDbCreds() = %v, want %v", dbConfig, want)
	}
}

func TestConfigDatabase(t *testing.T) {
	testDbConfig := DbConfig{
		DBUser:     "test_user",
		DBPassword: "test_learner",
		DBHost:     "127.0.0.1:3306",
		DBName:     "test_learning",
	}
	testDb := configDatabase(testDbConfig.DBUser, testDbConfig.DBPassword, testDbConfig.DBHost, testDbConfig.DBName)
	if testDb == nil {
		t.Error("configDatabase returned nil, expected a valid *gorm.DB")
	}
}
