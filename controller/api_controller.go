package controller

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/sitetester/sochain-api-parser/service"
	"net/http"
)

// InputBlock validation constraints are taken from https://github.com/asaskevich/govalidator
// https://sochain.com/api/#get-block
type InputBlock struct {
	Network      string `valid:"alpha,required"`
	NumberOrHash string `valid:"alphanum,required"`
}

// InputTransaction
// https://sochain.com/api/#get-tx
type InputTransaction struct {
	Network string `valid:"alpha,required"`
	Hash    string `valid:"alphanum,required"`
}

type ErrorResponse struct {
	Error string
}

type ApiController struct {
	service *service.ApiService
	cache   *cache.Cache
}

func NewApiController(cache *cache.Cache) *ApiController {
	return &ApiController{
		service: service.NewApiService(10),
		cache:   cache,
	}
}

const (
	ErrUnsupportedNetwork = "Unsupported network."
	ErrUnexpectedResponse = "Unexpected response."
)

// HandleBlockGetRoute
// https: //github.com/swaggo/swag#general-api-info
// @Summary      Show block
// @Description  Show block by network & number/hash
// @Tags         blocks
// @Param        network path string true "Network"
// @Param        blockNumberOrHash path string true "block number or hash"
// @Router       /block/{network}/{blockNumberOrHash} [get]
// @Success      200  {object}   service.DesiredBlockResponseData
func (c *ApiController) HandleBlockGetRoute(ctx *gin.Context) {
	network := ctx.Param("network")
	blockNumberOrHash := ctx.Param("blockNumberOrHash")

	_, err := govalidator.ValidateStruct(InputBlock{Network: network, NumberOrHash: blockNumberOrHash})
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if !c.service.SupportsNetwork(network) {
		ctx.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: ErrUnsupportedNetwork})
		return
	}

	// try to retrieve from cache (if not expired)
	cacheKey := fmt.Sprintf("%s_%s", network, blockNumberOrHash)
	if x, found := c.cache.Get(cacheKey); found {
		ctx.IndentedJSON(http.StatusOK, x.(*service.DesiredBlockResponseData))
		return
	}

	blockResponse, err := c.service.ApiClient.GetBlock(network, blockNumberOrHash)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	switch blockResponse.StatusCode {
	case http.StatusNotFound:
		// show error returned from remote server
		ctx.IndentedJSON(http.StatusNotFound, ErrorResponse{Error: blockResponse.Data.Blockid})
		return
	case http.StatusOK:
		desiredBlockResponseData := c.service.GetBlockInDesiredFormat(network, blockResponse.Data)
		// put in cache
		c.cache.Set(cacheKey, &desiredBlockResponseData, cache.DefaultExpiration)
		ctx.JSON(http.StatusOK, desiredBlockResponseData)
		return
	default:
		ctx.JSON(blockResponse.StatusCode, c.service.StatusCodeToMsg(blockResponse.StatusCode))
		return
	}
}

// HandleTransactionGetRoute
// https: //github.com/swaggo/swag#general-api-info
// @Summary      Show transaction
// @Description  Show transaction by network & hash
// @Tags         transactions
// @Param        network path string true "network"
// @Param        hash path string true "transaction hash"
// @Router       /tx/{network}/{hash} [get]
// @Success      200  {object}   service.DesiredTxResponseData
func (c *ApiController) HandleTransactionGetRoute(ctx *gin.Context) {
	network := ctx.Param("network")
	hash := ctx.Param("hash")

	_, err := govalidator.ValidateStruct(InputTransaction{Network: network, Hash: hash})
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if !c.service.SupportsNetwork(network) {
		ctx.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: ErrUnsupportedNetwork})
		return
	}

	// try to retrieve from cache (if not expired)
	cacheKey := fmt.Sprintf("%s_%s", network, hash)
	if x, found := c.cache.Get(cacheKey); found {
		ctx.IndentedJSON(http.StatusOK, x.(*service.DesiredTxResponseData))
		return
	}

	txResponse, err := c.service.ApiClient.GetTransaction(network, hash)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	switch txResponse.StatusCode {
	case http.StatusOK:
		desiredTxResponseData := c.service.GetTransactionInDesiredFormat(txResponse.Data)
		// put in cache
		c.cache.Set(cacheKey, &desiredTxResponseData, cache.DefaultExpiration)
		ctx.JSON(http.StatusOK, desiredTxResponseData)
		return
	default:
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{Error: ErrUnexpectedResponse})
		return
	}
}
