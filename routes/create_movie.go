package routes

import (
	"errors"
	"net/http"

	"github.com/arthurmarkus2013/simple-rest-server/database"
	"github.com/arthurmarkus2013/simple-rest-server/utils"
	"github.com/gin-gonic/gin"
)

func CreateMovie(ctx *gin.Context) {
	if ctx.Keys["role"] == utils.ADMIN {
		ctx.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	db := database.OpenDatabase()

	defer db.Close()

	var movie Movie
	err := ctx.BindJSON(&movie)

	if err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	stmt, err := db.Prepare("INSERT INTO movies (title, description, release_year) VALUES (?, ?, ?)")

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
		return
	}

	defer stmt.Close()

	result, err := stmt.Exec(movie.Title, movie.Description, movie.ReleaseYear)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
		return
	}

	affectedRows, err := result.RowsAffected()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
		return
	}

	if affectedRows != 1 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "movie created",
	})
}
