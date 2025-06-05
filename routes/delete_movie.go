package routes

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/arthurmarkus2013/simple-rest-server/database"
	"github.com/arthurmarkus2013/simple-rest-server/utils"
	"github.com/gin-gonic/gin"
)

func DeleteMovie(ctx *gin.Context) {
	if ctx.Keys["role"] != utils.ADMIN {
		ctx.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "id is required",
		})
		return
	}

	db := database.OpenDatabase()

	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM movies WHERE id = ?")

	if err != nil {
		slog.Error("something went wrong", "error", err.Error())

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
		return
	}

	defer stmt.Close()

	result, err := stmt.Exec(id)

	if err != nil {
		slog.Error("something went wrong", "error", err.Error())

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
		return
	}

	affectedRows, err := result.RowsAffected()

	if err != nil {
		slog.Error("something went wrong", "error", err.Error())

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
		return
	}

	if affectedRows != 1 {
		slog.Error("something went wrong", "error", affectedRows)

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "movie deleted",
	})
}
