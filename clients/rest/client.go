package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	restmodels "github.com/joaquinbejar/deribit-api/clients/rest/models"
	"github.com/joaquinbejar/deribit-api/pkg/deribit"
	"github.com/joaquinbejar/deribit-api/pkg/models"
	"github.com/sirupsen/logrus"
	"io"
	"net/url"
	"strings"

	"github.com/shopspring/decimal"
	"net/http"
)

const (
	StdDepth    = 10
	BtcTickSize = 0.5
)

type DeribitRestClient struct {
	Client      *http.Client
	ClientID    string
	ApiSecret   string
	BaseURL     string
	AccessToken *string
	Logger      *logrus.Logger
}

func NewDeribitRestClient(cfg *deribit.Configuration) *DeribitRestClient {
	return &DeribitRestClient{
		Client:      &http.Client{},
		ClientID:    cfg.ApiKey,
		ApiSecret:   cfg.SecretKey,
		BaseURL:     cfg.RestAddr,
		AccessToken: nil,
		Logger:      cfg.Logger,
	}
}

func (d *DeribitRestClient) GetAuthToken() (string, error) {
	baseURL := strings.Replace(d.BaseURL, "/ws", "", 1)
	d.Logger.Debugf("Original Base URL: %s", baseURL)

	baseURL = strings.TrimSuffix(baseURL, "/")
	d.Logger.Debugf("URL after trim: %s", baseURL)

	authURL := baseURL + "/public/auth"
	d.Logger.Debugf("Final auth URL: %s", authURL)

	params := url.Values{}
	params.Set("grant_type", "client_credentials")
	params.Set("client_id", d.ClientID)
	params.Set("client_secret", d.ApiSecret)

	fullURL := authURL + "?" + params.Encode()
	d.Logger.Debugf("Full URL with params: %s", fullURL)

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	d.Logger.Debugf("Response status: %s", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	d.Logger.Debugf("Response body: %s", string(body))

	var result struct {
		Result struct {
			AccessToken string `json:"access_token"`
		} `json:"result"`
		Error *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Error != nil {
		return "", fmt.Errorf("API error %d: %s", result.Error.Code, result.Error.Message)
	}

	d.AccessToken = &result.Result.AccessToken
	return *d.AccessToken, nil
}

func (d *DeribitRestClient) request(method string, params map[string]interface{}, private bool) (map[string]interface{}, error) {
	baseURL := strings.Replace(d.BaseURL, "/ws", "", 1)
	baseURL = strings.TrimSuffix(baseURL, "/")
	fullURL := baseURL + "/" + method

	d.Logger.Debugf("Request URL: %s", fullURL)
	d.Logger.Debugf("Request params: %+v", params)

	urlValues := url.Values{}
	for key, value := range params {
		urlValues.Set(key, fmt.Sprintf("%v", value))
	}

	fullURL = fullURL + "?" + urlValues.Encode()

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if private && d.AccessToken != nil {
		req.Header.Set("Authorization", "Bearer "+*d.AccessToken)
	}

	resp, err := d.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	d.Logger.Debugf("Response body: %s", string(body))

	var result struct {
		JSONRPC string                 `json:"jsonrpc"`
		ID      int                    `json:"id"`
		Result  map[string]interface{} `json:"result"`
		Error   *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Error != nil {
		return nil, fmt.Errorf("API error %d: %s", result.Error.Code, result.Error.Message)
	}

	return result.Result, nil
}

func (d *DeribitRestClient) GetOrderbook(instrument string, depth *int) (restmodels.OrderBook, error) {
	depthValue := StdDepth
	if depth != nil {
		depthValue = *depth
	}

	params := map[string]interface{}{
		"instrument_name": instrument,
		"depth":           depthValue,
	}

	d.Logger.Debugf("Getting orderbook for %s with depth %d", instrument, depthValue)

	result, err := d.request("public/get_order_book", params, false)
	if err != nil {
		return restmodels.OrderBook{}, fmt.Errorf("failed to get orderbook: %w", err)
	}

	var orderBook restmodels.OrderBook
	resultBytes, err := json.Marshal(result)
	if err != nil {
		return restmodels.OrderBook{}, fmt.Errorf("failed to marshal result: %w", err)
	}

	d.Logger.Debugf("Orderbook raw data: %s", string(resultBytes))

	if err := json.Unmarshal(resultBytes, &orderBook); err != nil {
		return restmodels.OrderBook{}, fmt.Errorf("failed to unmarshal order book: %w", err)
	}

	d.Logger.Debugf("Parsed orderbook: %+v", orderBook)
	return orderBook, nil
}

func (d *DeribitRestClient) GetTicker(instrument string) (decimal.Decimal, error) {
	params := map[string]interface{}{
		"instrument_name": instrument,
	}

	response, err := d.request("public/ticker", params, false)
	if err != nil {
		return decimal.Zero, err
	}

	lastPrice, exists := response["last_price"]
	if !exists {
		return decimal.Zero, errors.New("missing last_price field in response")
	}

	switch v := lastPrice.(type) {
	case float64:
		return decimal.NewFromFloat(v), nil
	case string:
		return decimal.NewFromString(v)
	default:
		return decimal.Zero, errors.New("last price is neither a number nor a string")
	}
}

func (d *DeribitRestClient) PlaceLimitOrder(instrument string,
	price decimal.Decimal,
	amount decimal.Decimal,
	direction restmodels.Direction) (restmodels.Order, error) {
	priceToTick := roundToTickSize(price, decimal.NewFromFloat(BtcTickSize))
	method := "private/sell"
	if direction == "buy" {
		method = "private/buy"
	}

	params := map[string]interface{}{
		"instrument_name": instrument,
		"price":           priceToTick.String(),
		"amount":          amount.String(),
		"type":            models.OrderTypeLimit,
		"post_only":       true,
	}

	orderResult, err := d.request(method, params, true)
	if err != nil {
		return restmodels.Order{}, err
	}

	order := restmodels.Order{}
	orderData, _ := json.Marshal(orderResult["order"])
	if err := json.Unmarshal(orderData, &order); err != nil {
		return restmodels.Order{}, err
	}

	return order, nil
}

//func (d *DeribitRestClient) PlaceMarketOrder(instrument string, amount decimal.Decimal, direction Direction) (Order, error) {
//	amountF64, exact := amount.Float64()
//	if !exact {
//		return Order{}, errors.New("Could not convert amount to f64")
//	}
//
//	method := "private/sell"
//	if direction == DirectionBuy {
//		method = "private/buy"
//	}
//
//	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
//	label := fmt.Sprintf("market%d", timestamp)
//
//	params := map[string]interface{}{
//		"instrument_name": instrument,
//		"amount":          amountF64,
//		"type":            OrderTypeMarket,
//		"label":           label,
//	}
//
//	logger, _ := zap.NewProduction()
//	defer logger.Sync()
//	sugar := logger.Sugar()
//
//	sugar.Debugf("Placing market order: method=%s, params=%v", method, params)
//	orderResult, err := d.request(method, params, true)
//	if err != nil {
//		return Order{}, err
//	}
//
//	order := Order{}
//	orderData, _ := json.Marshal(orderResult["order"])
//	if err := json.Unmarshal(orderData, &order); err != nil {
//		return Order{}, err
//	}
//
//	return order, nil
//}
//
//func (d *DeribitRestClient) CancelOrder(orderID string) (CancelOrderResponse, error) {
//	params := map[string]interface{}{
//		"order_id": orderID,
//	}
//
//	cancelOrderResult, err := d.request("private/cancel", params, true)
//	if err != nil {
//		return CancelOrderResponse{}, err
//	}
//
//	cancelOrderResponse := CancelOrderResponse{}
//	cancelOrderData, _ := json.Marshal(cancelOrderResult)
//	if err := json.Unmarshal(cancelOrderData, &cancelOrderResponse); err != nil {
//		return CancelOrderResponse{}, err
//	}
//
//	return cancelOrderResponse, nil
//}
//
//func (d *DeribitRestClient) GetOpenOrders(instrument string) (Orders, error) {
//	params := map[string]interface{}{
//		"instrument_name": instrument,
//	}
//
//	result, err := d.request("private/get_open_orders_by_instrument", params, true)
//	if err != nil {
//		return Orders{}, err
//	}
//
//	orders := Orders{}
//	ordersData, _ := json.Marshal(result)
//	if err := json.Unmarshal(ordersData, &orders); err != nil {
//		return Orders{}, err
//	}
//
//	return orders, nil
//}
//
//func (d *DeribitRestClient) GetAccountSummary() (AccountSummary, error) {
//	panic("Not implemented")
//}
//
//func (d *DeribitRestClient) GetPositions() ([]Position, error) {
//	panic("Not implemented")
//}
//
//func (d *DeribitRestClient) GetPosition(_instrument string) (*Position, error) {
//	panic("Not implemented")
//}
//
//func (d *DeribitRestClient) GetPortfolioMargins() (PortfolioMargins, error) {
//	panic("Not implemented")
//}
//
//func (d *DeribitRestClient) SimulatePortfolio(_position Position) (PortfolioMargins, error) {
//	panic("Not implemented")
//}
//
//func (d *DeribitRestClient) GetFundingRate(_instrument string, _startTimestamp int64, _endTimestamp int64) (decimal.Decimal, error) {
//	panic("Not implemented")
//}
