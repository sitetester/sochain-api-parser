package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/sitetester/sochain-api-parser/logger"
	"github.com/sitetester/sochain-api-parser/service"
	"net/http"
	"strconv"
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
	ErrNotFound           = "Not found." // todo: remove
	ErrSomethingWentWrong = "Something went wrong in calling external API, try again"
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

	blockResponse, err := c.apiService.ApiClient.GetBlock(network, blockHashOrNumber)
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
		ctx.JSON(blockResponse.StatusCode, c.statusCodeToMsg(blockResponse.StatusCode))
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
		logger.GetLogger().Errorf("Erro while performing API call: %s", err.Error())
		ctx.IndentedJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	switch txResponse.StatusCode {
	case http.StatusNotFound:
		ctx.IndentedJSON(http.StatusNotFound, ErrorResponse{Error: ErrNotFound})
		return
	case http.StatusOK:
		desiredTxResponseData := c.apiService.GetTransactionInDesiredFormat(txResponse.Data)
		// put in cache
		c.cache.Set(cacheKey, &desiredTxResponseData, cache.DefaultExpiration)
		ctx.JSON(http.StatusOK, desiredTxResponseData)
		return
	default:
		logger.GetLogger().
			WithFields(logrus.Fields{"network": network, "hash": hash}).
			Debug("Unexpected API response code: ", txResponse.StatusCode)
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: ErrSomethingWentWrong})
		return
	}
}

func (c ApiController) statusCodeToMsg(statusCode int) string {
	var msg string
	code := strconv.Itoa(statusCode)
	if code[0:1] == "4" {
		msg = "Bad Request."
	} else {
		msg = "Unexpected Response."
	}
	return msg
}
