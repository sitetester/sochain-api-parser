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

func (h *BlocksController) HandleBlockGetRoute(c *gin.Context) {
	network := c.Param("network")
	blockHashOrNumber := c.Param("blockHashOrNumber")

	if !h.blockService.SupportsNetwork(network) {
		c.IndentedJSON(http.StatusBadRequest, ErrorResponse{Error: "Unsupported network."})
		return
	}

	blockResponse := h.blockService.ApiClient.GetBlock(network, blockHashOrNumber)
	if blockResponse.Status != "success" {
		c.JSON(http.StatusOK, blockResponse)
		return
	}

	c.JSON(http.StatusOK, h.blockService.GetBlockInDesiredFormat(network, blockResponse))
	return
}

func (h *BlocksController) HandleTransactionRoute(c *gin.Context) {

}
