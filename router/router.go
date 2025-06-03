package router

import (
	"net/http"

	"github.com/arthurmarkus2013/simple-rest-server/routes"
	"github.com/gin-gonic/gin"
)

func Register_Routes(engine *gin.Engine) {
	engine.GET("/", func(ctx *gin.Context) {
		ctx.Writer.WriteString("Simple REST Server")
	})
	engine.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	protected := engine.Group("/movie")

	protected.POST("/create", routes.CreateMovie)
	protected.GET("/read", routes.ReadMovie)
	protected.POST("/update", routes.UpdateMovie)
	protected.POST("/delete", routes.DeleteMovie)

	engine.POST("/register", routes.Register)
	engine.POST("/login", routes.Login)
	engine.POST("/logout", routes.Logout)
}
