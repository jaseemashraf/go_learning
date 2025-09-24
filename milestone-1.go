package main

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

type User struct {
	FirstName string
	LastName  string
	Age       int
}

func main() {
	users := readUserDataFromExcel()
	for i, p := range users {
		fmt.Printf("User %d: %+v\n", i+1, p)
	}
}

func readUserDataFromExcel() (users []User) {
	f, error := excelize.OpenFile("users.xlsx")
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
