package routes

import (
	"github.com/arthurmarkus2013/simple-rest-server/database"
	"github.com/arthurmarkus2013/simple-rest-server/utils"
	"github.com/gin-gonic/gin"

	"errors"
	"log/slog"
	"net/http"
)

func ReadMovie(ctx *gin.Context) {
	if ctx.Keys["role"] == utils.NONE {
		ctx.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	db := database.OpenDatabase()

	defer db.Close()

	id := ctx.Request.PathValue("id")

	if id == "" {
		result, err := db.Query("SELECT * FROM movies")

		if err != nil {
			slog.Error("something went wrong", "error", err.Error())

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}

		defer result.Close()

		ctx.JSON(http.StatusOK, gin.H{
			"movies": result,
		})
	} else {
		result, err := db.Prepare("SELECT * FROM movies WHERE id = ?")

		if err != nil {
			slog.Error("something went wrong", "error", err.Error())

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}

		defer result.Close()

		ctx.JSON(http.StatusOK, gin.H{
			"movie": result,
		})
	}
}
