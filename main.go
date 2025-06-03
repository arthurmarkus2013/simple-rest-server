package main

import (
	"github.com/arthurmarkus2013/simple-rest-server/database"
	"github.com/arthurmarkus2013/simple-rest-server/router"
	"github.com/gin-gonic/gin"
)

func main() {
	database.InitializeDatabase()

	engine := gin.Default()
	router.Register_Routes(engine)
	engine.Run()
}
