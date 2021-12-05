package service

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/sitetester/sochain-api-parser/service/client"
	"net/http"
	"time"
)

type DesiredBlockResponseData struct {
	Network           string                  `json:"network"`
	BlockNo           int                     `json:"block_no"`
	Time              string                  `json:"time"` // actual type is `int`, but we need to display this time in string format in API response
	PreviousBlockhash string                  `json:"previous_blockhash"`
	NextBlockhash     string                  `json:"next_blockhash"`
	Size              int                     `json:"size"`
	Txs               []DesiredTxResponseData `json:"txs"`
	StatusCode        int                     `json:"-"`
}

type DesiredTxResponseData struct {
	Txid       string `json:"txid"`
	Time       string `json:"time"`
	Fee        string `json:"fee"`
	SentValue  string `json:"sent_value"`
	StatusCode int    `json:"-"` // this will be used in controller to show custom error message
}

type BlocksService struct {
	apiClient  *client.SoChainApiClient
	maxTxs     int // maximum number of transactions to parse
	cache      *cache.Cache
	txsService *TxsService
}

func NewBlocksService(maxTxs int, cache *cache.Cache, txsService *TxsService) *BlocksService {
	return &BlocksService{
		apiClient:  client.NewSoChainApiClient(),
		maxTxs:     maxTxs,
		cache:      cache,
		txsService: txsService,
	}
}

func (s *BlocksService) GetBlock(network string, blockNumberOrHash string) (DesiredBlockResponseData, error) {
	// try to retrieve from cache (if not expired)
	cacheKey := fmt.Sprintf("%s_%s", network, blockNumberOrHash)
	if x, found := s.cache.Get(cacheKey); found {
		desiredBlockResponseData := x.(*DesiredBlockResponseData)
		return *desiredBlockResponseData, nil
	}

	blockResponse, err := s.apiClient.GetBlock(network, blockNumberOrHash)
	if err != nil {
		return DesiredBlockResponseData{}, err
	}

	if blockResponse.IsSuccess() {
		desiredBlockResponseData := s.getBlockInDesiredFormat(network, blockResponse.Data)
		desiredBlockResponseData.StatusCode = http.StatusOK
		// put in cache
		s.cache.Set(cacheKey, &desiredBlockResponseData, cache.DefaultExpiration)
		return desiredBlockResponseData, nil
	}

	return DesiredBlockResponseData{StatusCode: blockResponse.StatusCode}, nil
}

func (s *BlocksService) getBlockInDesiredFormat(network string, blockResponseData client.BlockResponseData) DesiredBlockResponseData {
	timeInt64 := blockResponseData.Time

	return DesiredBlockResponseData{
		Network:           blockResponseData.Network,
		BlockNo:           blockResponseData.BlockNo,
		Time:              time.Unix(timeInt64, 0).Format("01/02/2006 15:04"),
		PreviousBlockhash: blockResponseData.PreviousBlockhash,
		NextBlockhash:     blockResponseData.NextBlockhash,
		Size:              blockResponseData.Size,
		Txs:               s.txsService.parseTransactions(network, blockResponseData.Txs, s.maxTxs),
	}
}
