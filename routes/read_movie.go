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

	id := ctx.Param("id")

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

		var movies []Movie

		for result.Next() {
			var movie Movie

			result.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseYear)

			movies = append(movies, movie)
		}

		ctx.JSON(http.StatusOK, gin.H{
			"movies": movies,
		})
	} else {
		stmt, err := db.Prepare("SELECT * FROM movies WHERE id = ?")

		if err != nil {
			slog.Error("something went wrong", "error", err.Error())

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}

		defer stmt.Close()

		result, err := stmt.Query(id)

		if err != nil {
			slog.Error("something went wrong", "error", err.Error())

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}

		defer result.Close()

		var movie Movie

		if result.Next() {
			err = result.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.ReleaseYear)

			if err != nil {
				slog.Error("something went wrong", "error", err.Error())

				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": "something went wrong",
				})
				return
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"movie": movie,
		})
	}
}
