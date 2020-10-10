package main

import (
	"hotelbooking/models"
	"hotelbooking/routes"
)

func main() {
	db := models.DBConn()
	db.AutoMigrate(&models.Guest{})
	r := routes.Routes(db)
	r.Run(":3030")
}
