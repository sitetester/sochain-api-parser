package controller

import (
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
	apiService    service.ApiService
	blocksService *service.BlocksService
	txsService    *service.TxsService
}

func NewApiController(cache *cache.Cache) *ApiController {
	txService := service.NewTxsService(cache)
	return &ApiController{
		apiService:    service.ApiService{},
		blocksService: service.NewBlocksService(10, cache, txService),
		txsService:    txService,
	}
}

const (
	ErrNotFound           = "Not found."
	ErrUnsupportedNetwork = "Unsupported network."
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

	if !c.apiService.SupportsNetwork(network) {
		ctx.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: ErrUnsupportedNetwork})
		return
	}

	desiredBlockResponseData, err := c.blocksService.GetBlock(network, blockNumberOrHash)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	switch desiredBlockResponseData.StatusCode {
	case http.StatusOK:
		ctx.JSON(http.StatusOK, desiredBlockResponseData)
		return
	case http.StatusNotFound:
		ctx.IndentedJSON(http.StatusNotFound, ErrorResponse{Error: ErrNotFound})
		return
	default:
		ctx.JSON(desiredBlockResponseData.StatusCode, ErrorResponse{Error: c.apiService.StatusCodeToMsg(desiredBlockResponseData.StatusCode)})
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

	if !c.apiService.SupportsNetwork(network) {
		ctx.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: ErrUnsupportedNetwork})
		return
	}

	desiredTxResponseData, err := c.txsService.GetTransaction(network, hash)
	if err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	switch desiredTxResponseData.StatusCode {
	case http.StatusOK:
		ctx.JSON(http.StatusOK, desiredTxResponseData)
		return
	case http.StatusNotFound:
		ctx.JSON(http.StatusNotFound, ErrorResponse{Error: ErrNotFound})
		return
	default:
		ctx.JSON(desiredTxResponseData.StatusCode, ErrorResponse{Error: c.apiService.StatusCodeToMsg(desiredTxResponseData.StatusCode)})
	}
}
