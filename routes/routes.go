package routes

import (
	"fmt"
	"hotelbooking/controllers"
	"hotelbooking/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	permissions "github.com/xyproto/permissions2"
)

type UserInput struct {
	Username string `json:"username"`
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var filename string

var perm, err = permissions.New2()
var userstate = perm.UserState()

func Routes(db *gorm.DB) *gin.Engine {
	r := gin.New()

	if err != nil {
		log.Fatal(err)
	}

	permissionHandler := func(c *gin.Context) {
		if perm.Rejected(c.Writer, c.Request) {
			c.AbortWithStatus(http.StatusForbidden)
			fmt.Fprint(c.Writer, "Permission denied!")
			return
		}
		c.Next()
	}

	r.Use(gin.Logger())
	r.Use(permissionHandler)
	r.Use(gin.Recovery())
	//userstate := perm.UserState()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	// /
	r.GET("/", func(c *gin.Context) {
		var input UserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		uI := models.User{
			Username: input.Username,
		}
		msg := ""
		msg += fmt.Sprintf("Has user: %v\n", userstate.HasUser(uI.Username))
		msg += fmt.Sprintf("Logged in on server: %v\n", userstate.IsLoggedIn(uI.Username))
		msg += fmt.Sprintf("Is confirmed: %v\n", userstate.IsConfirmed(uI.Username))
		msg += fmt.Sprintf("Username stored in cookies (or blank): %v\n", userstate.Username(c.Request))
		msg += fmt.Sprintf("Current user is logged in, has a valid cookie and *user rights*: %v\n", userstate.UserRights(c.Request))
		msg += fmt.Sprintf("Current user is logged in, has a valid cookie and *admin rights*: %v\n", userstate.AdminRights(c.Request))
		msg += fmt.Sprintln("\nTry: /register, /confirm, /remove, /login, /logout, /makeadmin, /clear, /data and /admin")
		c.String(http.StatusOK, msg)
		if userstate.IsLoggedIn(uI.Username) == true {

		}
	})

	// REGISTER

	r.GET("/register", func(c *gin.Context) {

		var input UserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		uI := models.User{
			Username: input.Username,
			FullName: input.FullName,
			Email:    input.Email,
			Password: input.Password,
		}

		userstate.AddUser(uI.Username, uI.FullName, uI.Email)
		userstate.SetPassword(uI.Username, uI.Password)
		c.JSON(200, gin.H{
			"username : ": uI.Username + " added with password : " + uI.Password,
		})

	})

	//CONFIRM
	r.GET("/confirm", func(c *gin.Context) {
		var input UserInput
		uI := models.User{
			Username: input.Username,
		}
		userstate.MarkConfirmed(uI.Username)
		c.String(http.StatusOK, fmt.Sprintf("User confirmed: %v\n", userstate.IsConfirmed(uI.Username)))
	})

	//REMOVE
	r.GET("/remove", func(c *gin.Context) {
		var input UserInput
		uI := models.User{
			Username: input.Username,
		}
		userstate.RemoveUser(uI.Username)
		c.String(http.StatusOK, fmt.Sprintf("User removed: %v\n", !userstate.HasUser(uI.Username)))
	})

	//LOGIN

	r.GET("/login", func(c *gin.Context) {
		var input UserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		uI := models.User{
			Username: input.Username,
			// FullName: input.FullName,
			// Email:    input.Email,
			Password: input.Password,
		}
		if userstate.CorrectPassword(uI.Username, uI.Password) {
			userstate.Login(c.Writer, uI.Username)
			c.JSON(200, gin.H{
				"message": uI.Username + " Login successfully",
			})
			filename = uI.Username
			filename = "./" + filename
			f, err := os.Create(filename)
			if err != nil {
				panic(err.Error())
			}
			defer f.Close()
			f.WriteString(uI.Username)
		} else {
			c.JSON(200, gin.H{
				"message": "username/password not matched",
			})
		}

	})

	//LOGOUT

	r.GET("/logout", func(c *gin.Context) {
		var input UserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		uI := models.User{
			Username: input.Username,
		}
		userstate.Logout(uI.Username)
		c.String(http.StatusOK, fmt.Sprintf("logged out successfully: %v\n", !userstate.IsLoggedIn(uI.Username)))
		e := os.Remove(filename)
		if e != nil {
			log.Fatal(e)
		}
	})

	//MAKE ADMIN
	r.GET("/makeadmin", func(c *gin.Context) {
		var input UserInput
		uI := models.User{
			Username: input.Username,
		}
		userstate.SetAdminStatus(uI.Username)
		c.String(http.StatusOK, fmt.Sprintf("is now administrator: %v\n", userstate.IsAdmin(uI.Username)))
	})

	//USER DATA
	r.GET("/data", func(c *gin.Context) {
		c.String(http.StatusOK, "user page that only logged in users must see!")
	})

	//ADMIN
	r.GET("/admin", func(c *gin.Context) {
		c.String(http.StatusOK, "super secret information that only logged in administrators must see!\n\n")
		if usernames, err := userstate.AllUsernames(); err == nil {
			c.String(http.StatusOK, "list of all users: "+strings.Join(usernames, ", "))
		}
	})

	// r.GET("/guests", controllers.GetAllGuest)
	r.GET("/guests", func(c *gin.Context) {

		dat, err := ioutil.ReadFile(filename)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"data": "CANT ACCESS",
			})
		}
		if userstate.IsLoggedIn(string(dat)) {
			db := c.MustGet("db").(*gorm.DB)
			var guest []models.Guest
			db.Find(&guest)
			c.JSON(http.StatusOK, gin.H{
				"data": guest,
			})
		}

	})

	r.GET("/guests/:id", controllers.GetGuest)
	r.POST("/guests", controllers.AddGuest)
	r.PATCH("/guests/:id", controllers.UpdateGuest)
	r.DELETE("/guests/:id", controllers.DeleteGuest)
	return r
}
