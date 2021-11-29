package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sitetester/sochain-api-parser/controller"
	"github.com/sitetester/sochain-api-parser/logger"
	"io"
	"log"
	"os"
)

func setupRouter(value string) *gin.Engine {
	gin.SetMode(value) // switch to test mode (to avoid debug output)
	engine := gin.Default()

	if value == gin.ReleaseMode {
		// disable Console Color, you don't need console color when writing the logs to file.
		gin.DisableConsoleColor()
		f, _ := os.Create("logs/gin.log")
		gin.DefaultWriter = io.MultiWriter(f)
		engine.Use(gin.Recovery())
	}

	apiController := controller.NewApiController()
	engine.GET("/", func(ctx *gin.Context) { ctx.String(200, "It works!") })
	engine.GET("/block/:network/:blockHashOrNumber", apiController.HandleBlockGetRoute)
	engine.GET("/tx/:network/:hash", apiController.HandleTransactionGetRoute)

	return engine
}

// return the value of the key
func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	var engine *gin.Engine

	envGinMode := goDotEnvVariable("EnvGinMode")
	if envGinMode != "" {
		engine = setupRouter(envGinMode)
	} else {
		engine = setupRouter(gin.DebugMode)
	}

	addr := "8081"
	addrEnv := goDotEnvVariable("HTTP_PORT") // create `.env` file with example value (HTTP_PORT=8182)
	if addrEnv != "" {
		addr = addrEnv
	}

	err := engine.Run(fmt.Sprintf(":%s", addr))
	if err != nil {
		panic(err)
	}

	logger.GetLogger().Debug("GetLogger() called ONCE") // just to check `singleton` functionality ;)
}
