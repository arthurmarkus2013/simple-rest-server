package main

import (
	"log/slog"
	"time"

	"github.com/arthurmarkus2013/simple-rest-server/database"
	"github.com/arthurmarkus2013/simple-rest-server/router"
	"github.com/arthurmarkus2013/simple-rest-server/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(gin.DefaultWriter, nil)))

	database.InitializeDatabase()

	go utils.PurgeExpiredTokens(time.Minute * 5)

	engine := gin.Default()
	router.Register_Routes(engine)
	engine.Run()
}
