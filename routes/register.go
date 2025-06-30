package routes

import (
	"log/slog"
	"net/http"

	"github.com/arthurmarkus2013/simple-rest-server/database"
	"github.com/arthurmarkus2013/simple-rest-server/utils"
	"github.com/gin-gonic/gin"
)

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

	stmt, err := db.Prepare("SELECT * FROM users WHERE username = ?")

	if err != nil {
		slog.Error("something went wrong", "error", err.Error())

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
		return
	}

	defer stmt.Close()

	result, err := stmt.Query(user.Username)

	if err != nil {
		slog.Error("something went wrong", "error", err.Error())

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
		return
	}

	defer result.Close()

	if result.Next() {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "user already exists",
		})
		return
	}

	stmt, err = db.Prepare("INSERT INTO users (username, password, role) VALUES (?, ?, ?)")

	if err != nil {
		slog.Error("something went wrong", "error", err.Error())

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
		return
	}

	defer stmt.Close()

	passwordHash, err := utils.HashPassword(user.Password)

	if err != nil {
		slog.Error("something went wrong", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
		return
	}

	_, err = stmt.Exec(user.Username, passwordHash, user.Role)

	if err != nil {
		slog.Error("something went wrong", "error", err.Error())

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "user registered",
	})
}
