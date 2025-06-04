package main

import (
	"log/slog"

	"github.com/arthurmarkus2013/simple-rest-server/database"
	"github.com/arthurmarkus2013/simple-rest-server/router"
	"github.com/gin-gonic/gin"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(gin.DefaultWriter, nil)))

	database.InitializeDatabase()

	engine := gin.Default()
	router.Register_Routes(engine)
	engine.Run()
}
