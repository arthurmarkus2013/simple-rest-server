package router

import (
	"net/http"

	"github.com/arthurmarkus2013/simple-rest-server/routes"
	"github.com/arthurmarkus2013/simple-rest-server/utils"
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

	protected.Use(authMiddleware)

	protected.POST("/create", routes.CreateMovie)
	protected.GET("/list", routes.ReadMovie)
	protected.GET("/list/:id", routes.ReadMovie)
	protected.POST("/update/:id", routes.UpdateMovie)
	protected.DELETE("/delete/:id", routes.DeleteMovie)

	engine.POST("/register", routes.Register)
	engine.POST("/login", routes.Login)
	engine.POST("/logout", routes.Logout)
}

func authMiddleware(ctx *gin.Context) {
	result, role, err := utils.ValidateToken(ctx.GetHeader("Authorization"))

	if err != nil || !result {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

	if ctx.Keys == nil {
		ctx.Keys = make(map[string]any)
	}

	ctx.Keys["role"] = role

	ctx.Next()
}
