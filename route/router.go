package route

import (
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/sitetester/sochain-api-parser/controller"
	ginSwagger "github.com/swaggo/gin-swagger"
	"time"

	_ "github.com/sitetester/sochain-api-parser/docs"

	swaggerFiles "github.com/swaggo/files"
)

const ApiVersion = "/api/v1"

// SetupRouter These annotations are taken from https://github.com/swaggo/swag#general-api-info
// @title Sochain API Explorer
// @version 1.0
// @description This is an example server using Sochain API at backend
// @host localhost:8081
// @BasePath /api/v1
// @accept json
// @produce json
func SetupRouter() *gin.Engine {
	engine := gin.Default()

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	engine.Use(gin.Recovery())

	// https://github.com/patrickmn/go-cache
	cache := cache.New(60*time.Minute, 10*time.Minute)

	apiController := controller.NewApiController(cache)
	v1 := engine.Group(ApiVersion)
	{
		v1.GET("/", func(ctx *gin.Context) { ctx.String(200, "It works!") })
		v1.GET("/block/:network/:blockNumberOrHash", apiController.HandleBlockGetRoute)
		v1.GET("/tx/:network/:hash", apiController.HandleTransactionGetRoute)
	}

	// e.g. /swagger/index.html
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return engine
}
