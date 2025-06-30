package routes

import (
	"log/slog"
	"net/http"

	"github.com/arthurmarkus2013/simple-rest-server/utils"
	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
	var credentials Credentials
	err := ctx.BindJSON(&credentials)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := utils.GenerateToken(credentials.Username, credentials.Password)

	if err != nil {
		slog.Error("something went wrong", "error", err.Error())

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func Logout(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")

	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "token is required",
		})
		return
	}

	result := utils.InvalidateToken(token)

	if !result {

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to logout",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "successfully logged out",
	})
}
