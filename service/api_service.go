package service

import (
	"github.com/sitetester/sochain-api-parser/service/client"
	"sort"
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
}

type DesiredTxResponseData struct {
	Txid      string `json:"txid"`
	Time      string `json:"time"`
	Fee       string `json:"fee"`
	SentValue string `json:"sent_value"`
}

type ApiService struct {
	ApiClient *client.SoChainApiClient
	maxTxs    int // maximum number of transactions to parse
}

func NewApiService(maxTxs int) *ApiService {
	return &ApiService{
		ApiClient: client.NewSoChainApiClient(),
		maxTxs:    maxTxs,
	}
}

func (s *ApiService) SupportsNetwork(network string) bool {
	supportedNetworks := []string{"BTC", "LTC", "DOGE"}

	for _, supportedNetwork := range supportedNetworks {
		if network == supportedNetwork {
			return true
		}
	}

	return false
}

func (s *ApiService) GetBlockInDesiredFormat(network string, blockResponseData client.BlockResponseData) DesiredBlockResponseData {
	timeInt64 := blockResponseData.Time

	return DesiredBlockResponseData{
		Network:           blockResponseData.Network,
		BlockNo:           blockResponseData.BlockNo,
		Time:              time.Unix(timeInt64, 0).Format("01/02/2006 15:04"),
		PreviousBlockhash: blockResponseData.PreviousBlockhash,
		NextBlockhash:     blockResponseData.NextBlockhash,
		Size:              blockResponseData.Size,
		Txs:               s.parseTransactions(network, blockResponseData.Txs),
	}
}

func (s *ApiService) parseTransactions(network string, blockHashes []string) []DesiredTxResponseData {
	var txResponses []client.TxResponse
	var desiredTxResponseData []DesiredTxResponseData

	hashes := blockHashes[:s.maxTxs] // [10] = 0-9
	if len(hashes) == 0 {
		return desiredTxResponseData
	}

	txResponseChan := make(chan client.TxResponse, s.maxTxs)

	// let's parse them concurrently
	parseTxs := func() {
		for index, hash := range hashes {
			go func(index int, hash string, txResponseChan chan client.TxResponse) {
				tx := s.ApiClient.GetTransaction(network, hash)
				// attach a custom sort order, since goroutines are launched in random order
				tx.CustomSortOrder = index

				txResponseChan <- tx

			}(index, hash, txResponseChan)
		}
	}

	// collect each parsed transaction from provided channel
	collectTxs := func() {
		count := 0
	forLoop:
		for {
			select {
			case txResponse := <-txResponseChan:
				count += 1
				// txResponses = append(txResponses, s.GetTransactionInDesiredFormat(txResponse.Data))
				txResponses = append(txResponses, txResponse)
				if count == s.maxTxs {
					break forLoop
				}
			}
		}
	}

	parseTxs()
	collectTxs()

	// sort transactions
	sort.Slice(txResponses, func(i, j int) bool {
		return txResponses[i].CustomSortOrder < txResponses[j].CustomSortOrder
	})

	// to desired display format
	for _, txResponse := range txResponses {
		desiredTxResponseData = append(desiredTxResponseData, s.GetTransactionInDesiredFormat(txResponse.Data))
	}

	return desiredTxResponseData
}

// GetTransactionInDesiredFormat https://yourbasic.org/golang/format-parse-string-time-date-example/
func (s *ApiService) GetTransactionInDesiredFormat(txResponseData client.TxResponseData) DesiredTxResponseData {
	return DesiredTxResponseData{
		Txid:      txResponseData.Txid,
		Time:      time.Unix(txResponseData.Time, 0).Format("01/02/2006 15:04"),
		Fee:       txResponseData.Fee,
		SentValue: txResponseData.SentValue,
	}
}
