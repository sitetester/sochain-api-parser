package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/patrickmn/go-cache"
	"github.com/sitetester/sochain-api-parser/controller"
	"github.com/sitetester/sochain-api-parser/logger"
	"io"
	"log"
	"os"
	"time"
)

func setupRouter(value string) *gin.Engine {
	gin.SetMode(value)
	engine := gin.Default()

	if value == gin.ReleaseMode {
		// disable Console Color, you don't need console color when writing the logs to file.
		gin.DisableConsoleColor()
		f, err := os.Create("logs/gin.log")
		if err != nil {
			panic(err)
		}
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
		engine.Use(gin.Recovery())
	}

	// https://github.com/patrickmn/go-cache
	cache := cache.New(60*time.Minute, 10*time.Minute)

	apiController := controller.NewApiController(cache)
	engine.GET("/", func(ctx *gin.Context) { ctx.String(200, "It works!") })
	engine.GET("/block/:network/:blockNumberOrHash", apiController.HandleBlockGetRoute)
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

	var ginMode = gin.DebugMode
	envGinMode := goDotEnvVariable("EnvGinMode")
	if envGinMode != "" {
		ginMode = envGinMode
	}
	engine = setupRouter(ginMode)

	addr := "8081"
	err := engine.Run(fmt.Sprintf(":%s", addr))
	if err != nil {
		panic(err)
	}

	logger.GetLogger().Debug("GetLogger() called ONCE") // just to check `singleton` functionality ;)
}
