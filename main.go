package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sitetester/sochain-api-parser/api/controller"
)

func pong(c *gin.Context) {
	c.String(200, "pong")
}

func main() {
	blocksController := controller.NewBlocksController()
	route := gin.Default()

	route.GET("/ping", pong)
	route.GET("/block/:network/:blockHashOrNumber", blocksController.HandleBlockGetRoute)
	route.GET("/block/:network/:hash", blocksController.HandleTransactionGetRoute)

	err := route.Run(":3000")
	if err != nil {
		panic(err)
	}
}
