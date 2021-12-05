package service

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/sitetester/sochain-api-parser/logger"
	"github.com/sitetester/sochain-api-parser/service/client"
	"net/http"
	"sort"
	"time"
)

type TxsService struct {
	apiClient *client.SoChainApiClient
	cache     *cache.Cache
}

func NewTxsService(cache *cache.Cache) *TxsService {
	return &TxsService{
		apiClient: client.NewSoChainApiClient(),
		cache:     cache,
	}
}

func (s *TxsService) GetTransaction(network string, hash string) (DesiredTxResponseData, error) {
	cacheKey := fmt.Sprintf("%s_%s", network, hash)
	if x, found := s.cache.Get(cacheKey); found {
		desiredTxResponseData := x.(*DesiredTxResponseData)
		return *desiredTxResponseData, nil
	}

	txResponse, err := s.apiClient.GetTransaction(network, hash)
	if err != nil {
		return DesiredTxResponseData{}, err
	}

	if txResponse.IsSuccess() {
		desiredTxResponseData := s.getTransactionInDesiredFormat(txResponse.Data)
		desiredTxResponseData.StatusCode = http.StatusOK
		// put in cache
		s.cache.Set(cacheKey, &desiredTxResponseData, cache.DefaultExpiration)
		return desiredTxResponseData, nil
	}

	return DesiredTxResponseData{StatusCode: txResponse.StatusCode}, nil
}

// getTransactionInDesiredFormat https://yourbasic.org/golang/format-parse-string-time-date-example/
func (s *TxsService) getTransactionInDesiredFormat(txResponseData client.TxResponseData) DesiredTxResponseData {
	return DesiredTxResponseData{
		Txid:      txResponseData.Txid,
		Time:      time.Unix(txResponseData.Time, 0).Format("01/02/2006 15:04"),
		Fee:       txResponseData.Fee,
		SentValue: txResponseData.SentValue,
	}
}

func (s *TxsService) parseTransactions(network string, txHashes []string, maxTxs int) []DesiredTxResponseData {
	var txResponses []client.TxResponse
	var desiredTxResponseData []DesiredTxResponseData

	if len(txHashes) == 0 {
		return desiredTxResponseData
	}

	var hashes []string
	if len(txHashes) >= maxTxs {
		hashes = txHashes[:maxTxs] // [10] = 0-9
	} else {
		hashes = txHashes[:]
	}

	txResponseChan := make(chan client.TxResponse, len(hashes))

	// let's parse them concurrently
	parseTxs := func() {
		for index, hash := range hashes {
			go s.parseConcurrently(network, index, hash, txResponseChan)
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
				txResponses = append(txResponses, txResponse)
				if count == len(hashes) {
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
		desiredTxResponseData = append(desiredTxResponseData, s.getTransactionInDesiredFormat(txResponse.Data))
	}

	return desiredTxResponseData
}

func (s *TxsService) parseConcurrently(network string, index int, hash string, txResponseChan chan client.TxResponse) {
	tx, err := s.apiClient.GetTransaction(network, hash)
	if err != nil {
		logger.GetLogger().Errorf("Error while retreiving txsService: %s", err.Error())
	}
	// attach a custom sort order, since goroutines are launched in random order
	tx.CustomSortOrder = index
	txResponseChan <- tx
}
