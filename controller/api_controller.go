package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/sitetester/sochain-api-parser/service"
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
	ErrUnsupportedNetwork = "Unsupported network."
	ErrUnexpectedResponse = "Unexpected response."
)

// HandleBlockGetRoute https://github.com/patrickmn/go-cache#usage
func (c *ApiController) HandleBlockGetRoute(ctx *gin.Context) {
	network := ctx.Param("network")
	blockNumberOrHash := ctx.Param("blockNumberOrHash")

	if !c.apiService.SupportsNetwork(network) {
		ctx.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: ErrUnsupportedNetwork})
		return
	}

	// try to retrieve from cache (if not expired)
	cacheKey := fmt.Sprintf("%s_%s", network, blockNumberOrHash)
	if x, found := c.cache.Get(cacheKey); found {
		ctx.IndentedJSON(http.StatusOK, x.(*service.DesiredBlockResponseData))
		return
	}

	blockResponse, err := c.apiService.ApiClient.GetBlock(network, blockNumberOrHash)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	switch blockResponse.StatusCode {
	case http.StatusNotFound:
		ctx.IndentedJSON(http.StatusNotFound, ErrorResponse{Error: blockResponse.Data.Blockid})
		return
	case http.StatusOK:
		desiredBlockResponseData := c.apiService.GetBlockInDesiredFormat(network, blockResponse.Data)
		// put in cache
		c.cache.Set(cacheKey, &desiredBlockResponseData, cache.DefaultExpiration)
		ctx.JSON(http.StatusOK, desiredBlockResponseData)
		return
	default:
		ctx.JSON(blockResponse.StatusCode, c.apiService.StatusCodeToMsg(blockResponse.StatusCode))
		return
	}
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

	txResponse, err := c.apiService.ApiClient.GetTransaction(network, hash)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	switch txResponse.StatusCode {
	case http.StatusOK:
		desiredTxResponseData := c.apiService.GetTransactionInDesiredFormat(txResponse.Data)
		// put in cache
		c.cache.Set(cacheKey, &desiredTxResponseData, cache.DefaultExpiration)
		ctx.JSON(http.StatusOK, desiredTxResponseData)
		return
	default:
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: ErrUnexpectedResponse})
		return
	}
}
