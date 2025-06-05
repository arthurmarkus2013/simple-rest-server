package routes

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/arthurmarkus2013/simple-rest-server/database"
	"github.com/arthurmarkus2013/simple-rest-server/utils"
	"github.com/gin-gonic/gin"
)

func CreateMovie(ctx *gin.Context) {
	if ctx.Keys["role"] != utils.ADMIN {
		ctx.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	db := database.OpenDatabase()

	defer db.Close()

	var movie Movie
	err := ctx.BindJSON(&movie)

	if err != nil {
		slog.Error("something went wrong", "error", err.Error())

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	stmt, err := db.Prepare("INSERT INTO movies (title, description, release_year) VALUES (?, ?, ?)")

	if err != nil {
		slog.Error("something went wrong", "error", err.Error())

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
		return
	}

	defer stmt.Close()

	result, err := stmt.Exec(movie.Title, movie.Description, movie.ReleaseYear)

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

	movieId, err := result.LastInsertId()

	if err != nil {
		slog.Error("something went wrong", "error", err.Error())

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
		return
	}

	movie.ID = int(movieId)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "movie created",
		"movie":   movie,
	})
}
