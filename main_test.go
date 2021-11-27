package main

import (
	"encoding/json"
	"github.com/sitetester/sochain-api-parser/api/controller"
	"github.com/sitetester/sochain-api-parser/api/service"
	"github.com/sitetester/sochain-api-parser/api/service/client"
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

func TestHandleBlockGetRouteWithInvalidNetwork(t *testing.T) {
	assertions := assert.New(t)

	w := launchRequest(t, "/block/BTC123/000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf")
	checkStatusCode(t, 400, w.Code)

	var errorResponse controller.ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&errorResponse)
	if err != nil {
		t.Fatalf(err.Error())
	}

	assertions.Equal("Unsupported network.", errorResponse.Error)
}

func TestHandleBlockGetRouteWithInvalidBlocNum(t *testing.T) {
	assertions := assert.New(t)

	w := launchRequest(t, "/block/BTC/abcd")
	checkStatusCode(t, 400, w.Code)

	var blockResponse client.BlockResponse
	err := json.NewDecoder(w.Body).Decode(&blockResponse)
	if err != nil {
		t.Fatalf(err.Error())
	}

	assertions.Equal("fail", blockResponse.Status)
}

// https://sochain.com/api/v2/get_block/BTC/000000000000034a7dedef4a161fa058a2d67a173a90155f3a2fe6fc132e0ebf
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

}
