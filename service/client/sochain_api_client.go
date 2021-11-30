package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type RemoteApiStatus string

const (
	StatusSuccess RemoteApiStatus = "success"
	StatusFail    RemoteApiStatus = "fail"
)

type BlockResponse struct {
	Status RemoteApiStatus   `json:"status"`
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

type TxResponse struct {
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
func (c *SoChainApiClient) GetBlock(network string, blockNumberOrHash string) BlockResponse {
	url := fmt.Sprintf("%s/get_block/%s/%s", c.ApiUrl, network, blockNumberOrHash)
	client := &http.Client{Timeout: c.timeout}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return BlockResponse{Status: StatusFail}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
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
	client := &http.Client{Timeout: c.timeout}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return TxResponse{Status: StatusFail}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var txResponse TxResponse
	if err := json.Unmarshal(body, &txResponse); err != nil {
		panic(err)
	}

	return txResponse
}
