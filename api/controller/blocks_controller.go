package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sitetester/sochain-api-parser/api/service"
	"net/http"
)

type SuccessResponse struct {
	Message string
}

type ErrorResponse struct {
	Error string
}

type BlocksController struct {
	blockService *service.BlockService
}

func NewBlocksController() *BlocksController {
	return &BlocksController{
		blockService: service.NewBlockService(10),
	}
}

func (c *BlocksController) HandleBlockGetRoute(ctx *gin.Context) {
	network := ctx.Param("network")
	blockHashOrNumber := ctx.Param("blockHashOrNumber")

	if !c.blockService.SupportsNetwork(network) {
		ctx.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: "Unsupported network."})
		return
	}

	blockResponse := c.blockService.ApiClient.GetBlock(network, blockHashOrNumber)
	// may be invalid block number/hash was provided ?
	if blockResponse.Status != "success" {
		// show response with relevant error message returned from remote server
		ctx.JSON(http.StatusBadRequest, blockResponse)
		return
	}

	ctx.JSON(http.StatusOK, c.blockService.GetBlockInDesiredFormat(network, blockResponse.Data))
	return
}

func (c *BlocksController) HandleTransactionGetRoute(ctx *gin.Context) {
	network := ctx.Param("network")
	hash := ctx.Param("hash")

	if !c.blockService.SupportsNetwork(network) {
		ctx.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: "Unsupported network."})
		return
	}

	txResponse := c.blockService.ApiClient.GetTransaction(network, hash)
	// may be invalid hash was provided ?
	if txResponse.Status != "success" {
		// show response with relevant error message returned from remote server
		ctx.JSON(http.StatusBadRequest, txResponse)
		return
	}

	ctx.JSON(http.StatusOK, c.blockService.GetTransactionInDesiredFormat(txResponse.Data))
	return
}
