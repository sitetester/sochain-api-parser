package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sitetester/sochain-api-parser/api/service"
	"net/http"
)

type ErrorResponse struct {
	Error string
}

type ApiController struct {
	apiService *service.ApiService
}

func NewApiController() *ApiController {
	return &ApiController{
		apiService: service.NewApiService(10),
	}
}

func (c *ApiController) HandleBlockGetRoute(ctx *gin.Context) {
	network := ctx.Param("network")
	blockHashOrNumber := ctx.Param("blockHashOrNumber")

	if !c.apiService.SupportsNetwork(network) {
		ctx.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: "Unsupported network."})
		return
	}

	blockResponse := c.apiService.ApiClient.GetBlock(network, blockHashOrNumber)
	// may be invalid block number/hash was provided ?
	if blockResponse.Status != "success" {
		// show response with relevant error message returned from remote server
		ctx.JSON(http.StatusBadRequest, blockResponse)
		return
	}

	ctx.JSON(http.StatusOK, c.apiService.GetBlockInDesiredFormat(network, blockResponse.Data))
	return
}

func (c *ApiController) HandleTransactionGetRoute(ctx *gin.Context) {
	network := ctx.Param("network")
	hash := ctx.Param("hash")

	if !c.apiService.SupportsNetwork(network) {
		ctx.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: "Unsupported network."})
		return
	}

	txResponse := c.apiService.ApiClient.GetTransaction(network, hash)
	// may be invalid hash was provided ?
	if txResponse.Status != "success" {
		// show response with relevant error message returned from remote server
		ctx.JSON(http.StatusBadRequest, txResponse)
		return
	}

	ctx.JSON(http.StatusOK, c.apiService.GetTransactionInDesiredFormat(txResponse.Data))
	return
}
