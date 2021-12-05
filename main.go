package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sitetester/sochain-api-parser/route"

	"github.com/sitetester/sochain-api-parser/logger"
	"io"
	"log"
	"os"
)

// return the value of the key
func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

// https://github.com/gin-gonic/gin#how-to-write-log-file
func setupFileLogger() {
	// not needed when writing the logs to file
	gin.DisableConsoleColor()
	f, err := os.Create("logs/gin.log")
	if err != nil {
		panic(err)
	}
	gin.DefaultWriter = io.MultiWriter(f)
}

func main() {
	var engine *gin.Engine

	var ginMode = gin.DebugMode
	envGinMode := goDotEnvVariable("EnvGinMode")
	if envGinMode != "" {
		ginMode = envGinMode
	}
	gin.SetMode(ginMode)

	if ginMode == gin.ReleaseMode {
		setupFileLogger()
	}

	engine = route.SetupRouter()

	addr := "8081"
	err := engine.Run(fmt.Sprintf(":%s", addr))
	if err != nil {
		panic(err)
	}

	logger.GetLogger().Debug("GetLogger() called ONCE") // just to check `singleton` functionality ;)
}
