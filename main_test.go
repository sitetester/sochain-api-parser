package main

import (
	"encoding/json"
	"github.com/sitetester/sochain-api-parser/controller"
	"github.com/sitetester/sochain-api-parser/service"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func checkStatusCode(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Fatalf("Expected %d, got %d", expected, actual)
	}
}

func launchRequest(t *testing.T, url string) *httptest.ResponseRecorder {
	r := setupRouter(true)

	// create a response recorder so you can inspect the response
	w := httptest.NewRecorder()

	// mock request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// perform the request
	r.ServeHTTP(w, req)

	return w
}

func parseErrorResponse(r *httptest.ResponseRecorder) controller.ErrorResponse {
	var er controller.ErrorResponse
	json.NewDecoder(r.Body).Decode(&er)
	return er
}

func TestHandleBlockGetRouteWithInvalidNetwork(t *testing.T) {
	assertions := assert.New(t)

	w := launchRequest(t, "/block/BTC123/000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf")
	checkStatusCode(t, 400, w.Code)

	errorResponse := parseErrorResponse(w)
	assertions.Equal(controller.ErrUnsupportedNetwork, errorResponse.Error)
}

func TestHandleBlockGetRouteWithInvalidBlocNum(t *testing.T) {
	assertions := assert.New(t)

	w := launchRequest(t, "/block/BTC/abcd")
	checkStatusCode(t, 400, w.Code)

	errorResponse := parseErrorResponse(w)
	assertions.Equal(controller.ErrInvalidInputProvided, errorResponse.Error)
}

// https://sochain.com/api/v2/get_block/BTC/000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf
// https://www.epochconverter.com
func TestHandleBlockGetRouteWithValidNetworkAndBlocNum(t *testing.T) {
	assertions := assert.New(t)

	w := launchRequest(t, "/block/BTC/000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf")
	checkStatusCode(t, 200, w.Code)

	var desiredBlockResponseData service.DesiredBlockResponseData
	err := json.NewDecoder(w.Body).Decode(&desiredBlockResponseData)
	if err != nil {
		t.Fatalf(err.Error())
	}

	assertions.Equal("BTC", desiredBlockResponseData.Network)
	assertions.Equal(200000, desiredBlockResponseData.BlockNo)
	assertions.Equal("09/22/2012 13:45", desiredBlockResponseData.Time)
	assertions.Equal("00000000000003a20def7a05a77361b9657ff954b2f2080e135ea6f5970da215", desiredBlockResponseData.PreviousBlockhash)
	assertions.Equal("00000000000002e3269b8a00caf315115297c626f954770e8398470d7f387e1c", desiredBlockResponseData.NextBlockhash)
	assertions.Len(desiredBlockResponseData.Txs, 10)
	assertions.Equal(247533, desiredBlockResponseData.Size)

	// let's check some transactions

	// first transaction
	// https://sochain.com/api/v2/tx/BTC/dbaf14e1c476e76ea05a8b71921a46d6b06f0a950f17c5f9f1a03b8fae467f10
	txId := "dbaf14e1c476e76ea05a8b71921a46d6b06f0a950f17c5f9f1a03b8fae467f10"
	checkTransaction(t, desiredBlockResponseData.Txs[0], txId, "09/22/2012 13:45", "0.0", "50.63517500")

	// last transaction
	// https://sochain.com/api/v2/tx/BTC/80efe43cf64a524d1417546a027786127ad87475f3af1c13b8f3719cd4268679
	txId = "80efe43cf64a524d1417546a027786127ad87475f3af1c13b8f3719cd4268679"
	checkTransaction(t, desiredBlockResponseData.Txs[9], txId, "09/22/2012 13:45", "0.0", "50.00000000")
}

func checkTransaction(t *testing.T, desiredTxResponseData service.DesiredTxResponseData, txId string, timeStr string, fee string, sentValue string) {
	assertions := assert.New(t)
	assertions.Equal(txId, desiredTxResponseData.Txid)
	assertions.Equal(timeStr, desiredTxResponseData.Time)
	assertions.Equal(fee, desiredTxResponseData.Fee)
	assertions.Equal(sentValue, desiredTxResponseData.SentValue)
}

// https://sochain.com/api/v2/tx/BTC/ee475443f1fbfff84ffba43ba092a70d291df233bd1428f3d09f7bd1a6054a1f
func TestHandleTransactionGetRouteWithInvalidValidNetwork(t *testing.T) {
	assertions := assert.New(t)

	w := launchRequest(t, "/tx/BTC123/ee475443f1fbfff84ffba43ba092a70d291df233bd1428f3d09f7bd1a6054a1f")
	checkStatusCode(t, 400, w.Code)

	errorResponse := parseErrorResponse(w)
	assertions.Equal(controller.ErrUnsupportedNetwork, errorResponse.Error)
}

func TestHandleTransactionGetRouteWithInvalidValidHash(t *testing.T) {
	assertions := assert.New(t)
	w := launchRequest(t, "/tx/BTC/xyz")
	checkStatusCode(t, 400, w.Code)

	errorResponse := parseErrorResponse(w)
	assertions.Equal(controller.ErrInvalidInputProvided, errorResponse.Error)
}
