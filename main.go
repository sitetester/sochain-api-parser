package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sitetester/sochain-api-parser/route"

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
	envGinMode := goDotEnvVariable("EnvGinMode")
	if envGinMode != "" {
		gin.SetMode(envGinMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	if gin.Mode() == gin.ReleaseMode {
		setupFileLogger()
	}

	engine := route.SetupRouter()
	err := engine.Run(":8081")
	if err != nil {
		panic(err)
	}
}
