package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

// LoggerToFile Log to file
// https://programmer.help/blogs/gin-framework-logging-with-logrus.html
// https://chowdera.com/2021/04/20210427171527655n.html
// https://golangbyexample.com/go-logger-rotation/
func LoggerToFile() gin.HandlerFunc {
	logFilePath := "logs"
	logFileName := "log.txt"

	// log file
	fileName := path.Join(logFilePath, logFileName)
	println("fileName: ", fileName)

	// write file
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}

	// instantiation
	logger := logrus.New()

	// set output
	logger.Out = src

	// set log level
	logger.SetLevel(logrus.DebugLevel)

	// format log
	logger.SetFormatter(&logrus.TextFormatter{})

	return func(c *gin.Context) {
		// start time
		startTime := time.Now()

		// processing request
		c.Next()

		// end time
		endTime := time.Now()

		// execution time
		latencyTime := endTime.Sub(startTime)

		// Request mode
		reqMethod := c.Request.Method

		// Request routing
		reqUri := c.Request.RequestURI

		// status code
		statusCode := c.Writer.Status()

		// Request IP
		clientIP := c.ClientIP()

		// log format
		logger.Infof("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}

// LoggerToMongo Log to MongoDB
func LoggerToMongo() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func LoggerToES() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

func LoggerToMQ() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
