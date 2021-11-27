package service

import (
	"github.com/sitetester/sochain-api-parser/api/service/client"
	"time"
)

type DesiredBlockResponseData struct {
	Network           string                  `json:"network"`
	BlockNo           int                     `json:"block_no"`
	Time              string                  `json:"time"` // actual type is `int`, but we need to display this time in string format in API response
	PreviousBlockhash string                  `json:"previous_blockhash"`
	NextBlockhash     string                  `json:"next_blockhash"`
	Size              int64                   `json:"size"`
	Txs               []DesiredTxResponseData `json:"txs"`
}

type DesiredTxResponseData struct {
	Txid      string `json:"txid"`
	Time      string `json:"time"`
	Fee       string `json:"fee"`
	SentValue string `json:"sent_value"`
}

type BlockService struct {
	ApiClient *client.SoChainApiClient
	maxTxs    int // maximum number of transactions to parse
}

func NewBlockService(maxTxs int) *BlockService {
	return &BlockService{
		ApiClient: client.NewSoChainApiClient(),
		maxTxs:    maxTxs,
	}
}

func (s *BlockService) SupportsNetwork(network string) bool {
	supportedNetworks := []string{"BTC", "LTC", "DOGE"}

	for _, supportedNetwork := range supportedNetworks {
		if network == supportedNetwork {
			return true
		}
	}

	return false
}

func (s *BlockService) GetBlockInDesiredFormat(network string, blockResponseData client.BlockResponseData) DesiredBlockResponseData {
	return DesiredBlockResponseData{
		Network:           blockResponseData.Network,
		BlockNo:           blockResponseData.BlockNo,
		Time:              time.Unix(blockResponseData.Time, 0).Format("01/02/2006 15:04"),
		PreviousBlockhash: blockResponseData.PreviousBlockhash,
		NextBlockhash:     blockResponseData.NextBlockhash,
		Size:              blockResponseData.Size,
		Txs:               s.parseTransactions(network, blockResponseData.Txs),
	}
}

func (s *BlockService) parseTransactions(network string, blockHashes []string) []DesiredTxResponseData {
	var desiredTxResponseData []DesiredTxResponseData

	hashes := blockHashes[:s.maxTxs] // [10] = 0-9
	if len(hashes) == 0 {
		return desiredTxResponseData
	}

	txResponseChan := make(chan client.TxResponse, s.maxTxs)
	for _, hash := range hashes {
		go func(hash string, txResponseChan chan client.TxResponse) {
			txResponseChan <- s.ApiClient.GetTransaction(network, hash)

		}(hash, txResponseChan)
	}

	// collect each parsed transaction from provided channel
	count := 0
forLoop:
	for {
		select {
		case txResponse := <-txResponseChan:
			count += 1

			data := txResponse.Data
			desiredTxResponseData = append(desiredTxResponseData, DesiredTxResponseData{
				Txid:      data.Txid,
				Time:      s.timeInt64ToString(data.Time),
				Fee:       data.Fee,
				SentValue: data.SentValue,
			})
			if count == s.maxTxs {
				break forLoop
			}
		}
	}

	return desiredTxResponseData
}

// https://yourbasic.org/golang/format-parse-string-time-date-example/
func (s BlockService) timeInt64ToString(blockTime int64) string {
	return time.Unix(blockTime, 0).Format("01/02/2006 15:04")
}
