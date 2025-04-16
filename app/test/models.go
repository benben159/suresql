package main

// Example struct
type User struct {
	Name   string `db:"name"`
	Email  string `db:"email"`
	Age    int    `db:"age"`
	Status string `db:"status"`
}

func (u *User) TableName() string {
	return "users"
}
