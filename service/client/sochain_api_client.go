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

const ApiResponseSuccess = "success"

type BlockResponse struct {
	StatusCode int               `json:"-"`
	Status     string            `json:"status"`
	Data       BlockResponseData `json:"data"`
}

func (r BlockResponse) IsSuccess() bool {
	return r.Status == ApiResponseSuccess
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
	StatusCode      int            `json:"-"`
	Status          string         `json:"status"`
	Data            TxResponseData `json:"data"`
	CustomSortOrder int            `json:"-"` // this will be used for sorting transactions
}

func (r TxResponse) IsSuccess() bool {
	return r.Status == ApiResponseSuccess
}

type TxResponseData struct {
	Network   string `json:"network"`
	Time      int64
	Txid      string `json:"txid"`
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

	if code == http.StatusOK {
		var blockResponse BlockResponse
		if err := json.Unmarshal(bytes, &blockResponse); err != nil {
			return blockResponse, err
		}
		return blockResponse, nil
	}

	// any other interesting response ? log it anyway
	logger.GetLogger().
		WithFields(logrus.Fields{"url": url, "statusCode": code, "body": string(bytes)}).
		Debug("Unexpected API response.")

	return BlockResponse{StatusCode: code}, nil
}

// GetTransaction https://sochain.com/api/#get-tx
// For some reason, there is inconsistent http response code for invalid hash
// e.g. https://chain.so/api/v2/get_tx/DOGE/abc123 (/get_tx) returns 404,
// while https://sochain.com/api/v2/tx/DOGE/abc123 (/tx) returns 500
func (c *SoChainApiClient) GetTransaction(network string, hash string) (TxResponse, error) {
	url := fmt.Sprintf("%s/tx/%s/%s", c.ApiUrl, network, hash)
	bytes, code, err := c.performRequest(url)
	if err != nil { // e.g. remote API is temporarily down
		logger.GetLogger().Errorf("Error while performing API call: %s", err.Error())
		return TxResponse{}, err
	}

	if code == http.StatusOK {
		var txResponse TxResponse
		if err := json.Unmarshal(bytes, &txResponse); err != nil {
			return txResponse, err
		}
		return txResponse, nil
	}

	// log any other response
	logger.GetLogger().WithFields(logrus.Fields{"url": url}).Debug("Unexpected API response.")
	return TxResponse{StatusCode: code}, nil
}

// Currently only `Timeout` option is configured, to configure `Transport, take a look at
// http://tleyden.github.io/blog/2016/11/21/tuning-the-go-http-client-library-for-load-testing/
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
