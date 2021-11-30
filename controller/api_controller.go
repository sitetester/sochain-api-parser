package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/sitetester/sochain-api-parser/logger"
	"github.com/sitetester/sochain-api-parser/service"
	"github.com/sitetester/sochain-api-parser/service/client"
	"net/http"
)

type ErrorResponse struct {
	Error string
}

type ApiController struct {
	apiService *service.ApiService
	cache      *cache.Cache
}

func NewApiController(cache *cache.Cache) *ApiController {
	return &ApiController{
		apiService: service.NewApiService(10),
		cache:      cache,
	}
}

const (
	ErrUnsupportedNetwork   = "Unsupported network."
	ErrInvalidInputProvided = "Invalid input provided"
)

// HandleBlockGetRoute https://github.com/patrickmn/go-cache#usage
func (c *ApiController) HandleBlockGetRoute(ctx *gin.Context) {
	network := ctx.Param("network")
	blockHashOrNumber := ctx.Param("blockHashOrNumber")

	if !c.apiService.SupportsNetwork(network) {
		ctx.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: ErrUnsupportedNetwork})
		return
	}

	// try to retrieve from cache (if not expired)
	cacheKey := fmt.Sprintf("%s_%s", network, blockHashOrNumber)
	if x, found := c.cache.Get(cacheKey); found {
		ctx.IndentedJSON(http.StatusOK, x.(*service.DesiredBlockResponseData))
		return
	}

	blockResponse := c.apiService.ApiClient.GetBlock(network, blockHashOrNumber)
	// may be invalid block number/hash was provided ?
	if blockResponse.Status != client.StatusSuccess {
		logger.GetLogger().
			WithFields(logrus.Fields{"network": network, "blockHashOrNumber": blockHashOrNumber}).
			Debug("bad request!")

		// show response with relevant error message returned from remote server
		ctx.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: ErrInvalidInputProvided})
		return
	}

	// put in cache
	desiredBlockResponseData := c.apiService.GetBlockInDesiredFormat(network, blockResponse.Data)
	c.cache.Set(cacheKey, &desiredBlockResponseData, cache.DefaultExpiration)

	ctx.JSON(http.StatusOK, desiredBlockResponseData)
	return
}

func (c *ApiController) HandleTransactionGetRoute(ctx *gin.Context) {
	network := ctx.Param("network")
	hash := ctx.Param("hash")

	if !c.apiService.SupportsNetwork(network) {
		ctx.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: ErrUnsupportedNetwork})
		return
	}

	// try to retrieve from cache (if not expired)
	cacheKey := fmt.Sprintf("%s_%s", network, hash)
	if x, found := c.cache.Get(cacheKey); found {
		ctx.IndentedJSON(http.StatusOK, x.(*service.DesiredTxResponseData))
		return
	}

	txResponse := c.apiService.ApiClient.GetTransaction(network, hash)
	// may be invalid hash was provided ?
	if txResponse.Status != client.StatusSuccess {
		logger.GetLogger().WithFields(logrus.Fields{"network": network, "hash": hash}).Debug("bad request!")
		// show response with relevant error message returned from remote server
		ctx.JSON(http.StatusBadRequest, ErrorResponse{Error: ErrInvalidInputProvided})
		return
	}

	// put in cache
	desiredTxResponseData := c.apiService.GetTransactionInDesiredFormat(txResponse.Data)
	c.cache.Set(cacheKey, &desiredTxResponseData, cache.DefaultExpiration)

	ctx.JSON(http.StatusOK, desiredTxResponseData)
	return
}
