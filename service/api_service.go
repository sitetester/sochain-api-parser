package service

import "strconv"

type ApiService struct {
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

func (s *ApiService) StatusCodeToMsg(statusCode int) string {
	var msg string
	code := strconv.Itoa(statusCode)
	if code[0:1] == "4" {
		msg = "Bad Request."
	} else {
		msg = "Unexpected Response."
	}
	return msg
}
