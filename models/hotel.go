package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Guest struct {
	ID         int    `json:"id" gorm:"primary_key" gorm:"autoIncrement"`
	FirstName  string `json:"firstName"`
	Lastname   string `json:"lastname"`
	EntryDate  string `json:"entryDate"`
	ExitDate   string `json:"exitDate"`
	RoomNumber int    `json:"roomNumber"`
}
