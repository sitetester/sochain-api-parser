package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sitetester/sochain-api-parser/api/controller"
)

func pong(c *gin.Context) {
	c.String(200, "pong")
}

func setupRouter(inTestMode bool) *gin.Engine {

	if inTestMode {
		// switch to test mode (to avoid debug output)
		gin.SetMode(gin.TestMode)
	}

	apiController := controller.NewApiController()

	route := gin.Default()
	route.GET("/block/:network/:blockHashOrNumber", apiController.HandleBlockGetRoute)
	route.GET("/tx/:network/:hash", apiController.HandleTransactionGetRoute)

	return route
}

func main() {
	route := setupRouter(false)
	err := route.Run(":3000")
	if err != nil {
		panic(err)
	}
}
