package main

import (
	"fmt"

	orm "github.com/medatechnology/simpleorm"
	"github.com/medatechnology/suresql"
)

func TestCreateTable(db *suresql.SureSQLDB) {
	err := (*db).ExecOneSQL("CREATE TABLE IF NOT EXISTS users (name TEXT, email TEXT, age INT, status TEXT)")
	if err.Error != nil {
		fmt.Println("Error executing create table users")
		panic(err)
	}
}

func TestInsertTable(db *suresql.SureSQLDB) {
	// Test Write
	// Convert structs to DBRecord
	users := []User{
		{"Alice", "alice@example.com", 30, "active"},
		{"Bob", "bob@example.com", 25, "active"},
		{"Penduser", "penduser@example.com", 15, "pending"},
		{"Newuser", "newuser@example.com", 45, "new"},
		{"Jonathan", "jonathan@example.com", 9, "active"},
	}

	var dbRecords []orm.DBRecord
	for _, user := range users {
		fmt.Println("User ==> ", user)
		record, err := orm.TableStructToDBRecord(&user)
		fmt.Println("Record ==> ", record)
		if err != nil {
			fmt.Println("Error creating DBRecord:", err)
			continue
		}
		dbRecords = append(dbRecords, record)
	}
	fmt.Println("Records = ", len(dbRecords))
	fmt.Println("Records:", dbRecords)
	// Write to database
	_, err := (*db).InsertManyDBRecords(dbRecords, true)
	if err != nil {
		fmt.Println("Error writing records to database:", err)
	} else {
		fmt.Println("Records inserted successfully!")
	}
}
