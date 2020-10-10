package controllers

import (
	"hotelbooking/models"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AddGuestInput struct {
	FirstName  string `json:"firstName"`
	Lastname   string `json:"lastname"`
	EntryDate  string `json:"entryDate"`
	ExitDate   string `json:"exitDate"`
	RoomNumber int    `json:"roomNumber"`
}

type UpdateGuestInput struct {
	FirstName  string `json:"firstName"`
	Lastname   string `json:"lastname"`
	EntryDate  string `json:"entryDate"`
	ExitDate   string `json:"exitDate"`
	RoomNumber int    `json:"roomNumber"`
}

// func GetAllGuest(c *gin.Context) {

// 	db := c.MustGet("db").(*gorm.DB)
// 	var guest []models.Guest
// 	db.Find(&guest)
// 	c.JSON(http.StatusOK, gin.H{
// 		"data": guest,
// 	})
// }

func AddGuest(c *gin.Context) {

	var input AddGuestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aG := models.Guest{
		FirstName:  input.FirstName,
		Lastname:   input.Lastname,
		EntryDate:  input.EntryDate,
		ExitDate:   input.ExitDate,
		RoomNumber: input.RoomNumber,
	}

	db := c.MustGet("db").(*gorm.DB)
	db.Create(&aG)

	c.JSON(http.StatusOK, gin.H{
		"data": aG,
	})
}

func GetGuest(c *gin.Context) {
	var guest models.Guest
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("id = ?", c.Param("id")).First(&guest).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Record not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": guest,
	})
}

func UpdateGuest(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var guest models.Guest
	if err := db.Where("id = ?", c.Param("id")).First(&guest).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "record not found",
		})
		return
	}
	var input UpdateGuestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var uG models.Guest
	uG.FirstName = input.FirstName
	uG.Lastname = input.Lastname
	uG.EntryDate = input.EntryDate
	uG.ExitDate = input.ExitDate
	uG.RoomNumber = input.RoomNumber

	db.Model(&guest).Updates(uG)
	c.JSON(http.StatusOK, gin.H{
		"data": guest,
	})
}

func DeleteGuest(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var guest models.Guest
	if err := db.Where("id = ?", c.Param("id")).First(&guest).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Record not found",
		})
		return
	}
	db.Delete(&guest)
	c.JSON(http.StatusOK, gin.H{
		"data": true,
	})
}
