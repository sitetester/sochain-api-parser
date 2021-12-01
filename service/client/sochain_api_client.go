package client

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sitetester/sochain-api-parser/logger"
	"io/ioutil"
	"net/http"
	"time"
)

type RemoteApiStatus string

type BlockResponse struct {
	StatusCode int               `json:"-"`
	Status     RemoteApiStatus   `json:"status"`
	Data       BlockResponseData `json:"data"`
}

type BlockResponseData struct {
	Network           string   `json:"network"`
	Blockid           string   `json:"blockid"` // this field is needed in case of "status": "fail" response
	Blockhash         string   `json:"blockhash"`
	BlockNo           int      `json:"block_no"`
	Time              int64    `json:"time"` // int64 is needed for conversion to string
	PreviousBlockhash string   `json:"previous_blockhash"`
	NextBlockhash     string   `json:"next_blockhash"`
	Size              int      `json:"size"`
	Txs               []string `json:"txs"`
}

type TxResponse struct {
	StatusCode      int             `json:"-"`
	Status          RemoteApiStatus `json:"status"`
	Data            TxResponseData  `json:"data"`
	CustomSortOrder int             `json:"-"` // this will be used for sorting transactions
}

type TxResponseData struct {
	Network   string `json:"network"`
	Time      int64
	Txid      string `json:"txid"` /**/
	Fee       string `json:"fee"`
	SentValue string `json:"sent_value"`
}

type SoChainApiClient struct {
	ApiUrl  string
	timeout time.Duration
}

// NewSoChainApiClient - Golang constructor ;)
func NewSoChainApiClient() *SoChainApiClient {
	return &SoChainApiClient{ApiUrl: "https://sochain.com/api/v2", timeout: 5 * time.Second}
}

// GetBlock https://sochain.com/api/#get-block
// https://stackoverflow.com/questions/50676817/does-the-http-request-automatically-retry
func (c *SoChainApiClient) GetBlock(network string, blockNumberOrHash string) (BlockResponse, error) {
	url := fmt.Sprintf("%s/get_block/%s/%s", c.ApiUrl, network, blockNumberOrHash)
	bytes, code, err := c.performRequest(url)
	if err != nil { // e.g. remote API is temporarily down
		logger.GetLogger().Errorf("Erro while performing API call: %s", err.Error())
		return BlockResponse{}, err
	}

	if code == http.StatusOK || code == http.StatusNotFound {
		var blockResponse BlockResponse
		if err := json.Unmarshal(bytes, &blockResponse); err != nil {
			return BlockResponse{}, err
		}
		blockResponse.StatusCode = code
		return blockResponse, nil
	}

	// any other interesting response ? log it anyway
	logger.GetLogger().
		WithFields(logrus.Fields{"url": url, "statusCode": code, "body": string(bytes)}).
		Debug("Unexpected API response.")

	return BlockResponse{StatusCode: code}, nil
}

// GetTransaction https://sochain.com/api/#get-tx
func (c *SoChainApiClient) GetTransaction(network string, hash string) (TxResponse, error) {
	url := fmt.Sprintf("%s/tx/%s/%s", c.ApiUrl, network, hash)
	bytes, code, err := c.performRequest(url)
	if err != nil { // e.g. remote API is temporarily down
		return TxResponse{}, err
	}

	switch code {
	case http.StatusOK:
		var txResponse TxResponse
		if err := json.Unmarshal(bytes, &txResponse); err != nil {
			return TxResponse{}, err
		}
		txResponse.StatusCode = code
		return txResponse, nil

	case http.StatusNotFound:

	default:

	}

	// some other 5xx code ?
	return TxResponse{StatusCode: code}, nil
}

func (c *SoChainApiClient) performRequest(url string) ([]byte, int, error) {
	client := &http.Client{Timeout: c.timeout}
	resp, err := client.Get(url)
	if err != nil { // e.g. remote API is temporarily down
		return []byte(""), resp.StatusCode, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}
