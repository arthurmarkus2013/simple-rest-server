package routes

import (
	"net/http"

	"github.com/arthurmarkus2013/simple-rest-server/database"
	"github.com/gin-gonic/gin"
)

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

func Register(ctx *gin.Context) {
	var user User
	err := ctx.BindJSON(&user)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	db := database.OpenDatabase()

	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO users (username, password) VALUES (?, ?)")

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(user.Username, user.Password)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "user registered",
	})
}
