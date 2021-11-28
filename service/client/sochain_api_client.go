package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type BlockResponse struct {
	Status string            `json:"status"`
	Data   BlockResponseData `json:"data"`
}

type BlockResponseData struct {
	Network           string   `json:"network"`
	Blockhash         string   `json:"blockhash"`
	BlockNo           int      `json:"block_no"`
	Blockid           string   `json:"blockid"` // this field is needed in case of "status": "fail" response
	Time              int64    `json:"time"`    // int64 is needed for conversion to string
	PreviousBlockhash string   `json:"previous_blockhash"`
	NextBlockhash     string   `json:"next_blockhash"`
	Size              int      `json:"size"`
	Txs               []string `json:"txs"`
}

const (
	// StatusSuccess  Indicates API response with `status` = success
	StatusSuccess = "success"
	// StatusFail Indicates API response with `status` = fail
	StatusFail = "fail"
)

type TxResponse struct {
	Status          string         `json:"status"`
	Data            TxResponseData `json:"data"`
	CustomSortOrder int            `json:"-"` // this will be used for sorting transactions
}

type TxResponseData struct {
	Network   string `json:"network"`
	Time      int64
	Txid      string `json:"txid"` /**/
	Fee       string `json:"fee"`
	SentValue string `json:"sent_value"`
}

type SoChainApiClient struct {
	ApiUrl string
}

// NewSoChainApiClient - Golang constructor ;)
func NewSoChainApiClient() *SoChainApiClient {
	return &SoChainApiClient{ApiUrl: "https://sochain.com/api/v2"}
}

// GetBlock https://sochain.com/api/#get-block
// https://stackoverflow.com/questions/50676817/does-the-http-request-automatically-retry
func (c *SoChainApiClient) GetBlock(network string, blockNumberOrHash string) BlockResponse {
	url := fmt.Sprintf("%s/get_block/%s/%s", c.ApiUrl, network, blockNumberOrHash)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return BlockResponse{Status: StatusFail}
	}

	var blockResponse BlockResponse
	if err := json.Unmarshal(body, &blockResponse); err != nil {
		panic(err)
	}

	return blockResponse
}

// GetTransaction https://sochain.com/api/#get-tx
func (c *SoChainApiClient) GetTransaction(network string, hash string) TxResponse {
	url := fmt.Sprintf("%s/tx/%s/%s", c.ApiUrl, network, hash)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return TxResponse{Status: StatusFail}
	}

	var txResponse TxResponse
	if err := json.Unmarshal(body, &txResponse); err != nil {
		panic(err)
	}

	return txResponse
}
