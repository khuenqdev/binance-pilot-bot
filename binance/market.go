package binance

import "fmt"

type ServerTimeResource struct {
    ServerTime int `json:"serverTime,omitempty"`
}

type ExchangeInfoResource struct {

}

type MarketServiceInterface interface {
    Ping() error
    CheckServerTime() (*ServerTimeResource, error)
    GetExchangeInfo() (*ExchangeInfoResource, error)
}

type MarketService struct {
    client *Client
}

func (s *MarketService) Ping() error {
    path := fmt.Sprintf("%s/%s", ApiBasePathV1, "ping")
    err := s.client.Get(path, nil, nil)
    return err
}

func (s *MarketService) CheckServerTime() (*ServerTimeResource, error) {
    path := fmt.Sprintf("%s/%s", ApiBasePathV1, "time")
    resource := new(ServerTimeResource)
    err := s.client.Get(path, resource, nil)
    return resource, err
}

func (s *MarketService) GetExchangeInfo() (*ExchangeInfoResource, error) {
    path := fmt.Sprintf("%s/%s", ApiBasePathV1, "exchangeInfo")
    resource := new(ExchangeInfoResource)
    err := s.client.Get(path, resource, nil)
    return resource, err
}
