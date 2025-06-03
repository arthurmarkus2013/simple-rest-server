package routes

import (
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

	token, err := utils.GenerateToken(credentials.Username, credentials.Password, utils.USER)

	if err != nil {
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
	result := utils.InvalidateToken(ctx.GetHeader("Authorization"))

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
