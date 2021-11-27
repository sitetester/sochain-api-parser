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

func (bc *BlocksController) HandleBlockGetRoute(c *gin.Context) {
	network := c.Param("network")
	blockHashOrNumber := c.Param("blockHashOrNumber")

	if !bc.blockService.SupportsNetwork(network) {
		c.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: "Unsupported network."})
		return
	}

	blockResponse := bc.blockService.ApiClient.GetBlock(network, blockHashOrNumber)
	if blockResponse.Status != "success" {
		c.JSON(http.StatusOK, blockResponse)
		return
	}

	c.JSON(http.StatusOK, bc.blockService.GetBlockInDesiredFormat(network, blockResponse))
	return
}

func (bc *BlocksController) HandleTransactionRoute(c *gin.Context) {

}
