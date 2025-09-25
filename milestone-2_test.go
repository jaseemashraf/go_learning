package main

import (
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
