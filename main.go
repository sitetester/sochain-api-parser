package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sitetester/sochain-api-parser/controller"
)

func setupRouter(inTestMode bool) *gin.Engine {

	if inTestMode {
		// switch to test mode (to avoid debug output)
		gin.SetMode(gin.TestMode)
	}

	apiController := controller.NewApiController()

	route := gin.Default()
	route.GET("/", func(ctx *gin.Context) { ctx.String(200, "It works!") })
	route.GET("/block/:network/:blockHashOrNumber", apiController.HandleBlockGetRoute)
	route.GET("/tx/:network/:hash", apiController.HandleTransactionGetRoute)

	return route
}

func main() {
	route := setupRouter(false)
	err := route.Run(":8081")
	if err != nil {
		panic(err)
	}
}
