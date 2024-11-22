package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	restmodels "github.com/joaquinbejar/deribit-api/internal/rest/models"
	"github.com/joaquinbejar/deribit-api/pkg/deribit"
	"github.com/joaquinbejar/deribit-api/pkg/models"
	"io"

	"github.com/shopspring/decimal"

	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
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
}

func NewDeribitRestClient(cfg *deribit.Configuration) *DeribitRestClient {
	return &DeribitRestClient{
		Client:      &http.Client{},
		ClientID:    cfg.ApiKey,
		ApiSecret:   cfg.SecretKey,
		BaseURL:     cfg.RestBaseURL,
		AccessToken: nil,
	}
}

func (d *DeribitRestClient) GetAuthToken() (string, error) {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "public/auth",
		"params": map[string]interface{}{
			"grant_type":    "client_credentials",
			"client_id":     d.ClientID,
			"client_secret": d.ApiSecret,
			"timestamp":     timestamp,
		},
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/public/auth", d.BaseURL), strings.NewReader(string(reqBody)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	if errorField, exists := response["error"]; exists {
		errorMap := errorField.(map[string]interface{})
		return "", errors.New(fmt.Sprintf("API error: code %d, message %s", int(errorMap["code"].(float64)), errorMap["message"].(string)))
	}

	var auth AuthResponse
	authResult, _ := json.Marshal(response["result"])
	if err := json.Unmarshal(authResult, &auth); err != nil {
		return "", errors.New(fmt.Sprintf("Failed to parse auth response: %v", err))
	}

	return auth.AccessToken, nil
}

func (d *DeribitRestClient) request(method string, params map[string]interface{}, private bool) (map[string]interface{}, error) {
	methodURL := fmt.Sprintf("%s/%s", d.BaseURL, method)
	queryParams := url.Values{}
	for key, value := range params {
		queryParams.Set(key, fmt.Sprintf("%v", value))
	}

	methodURL += "?" + queryParams.Encode()
	req, err := http.NewRequest("GET", methodURL, nil)
	if err != nil {
		return nil, err
	}

	if private {
		token, err := d.GetAuthToken()
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	req.Header.Set("Content-Type", "application/json")

	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Println("Failed to sync logger")
		}
	}(logger)
	sugar := logger.Sugar()

	sugar.Debugf("Sending request to URL: %s", req.URL.String())
	resp, err := d.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Failed to close response body")
		}
	}(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("Request failed: %s", resp.Status))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	sugar.Debugf("Request response: %v", result)

	if errorField, exists := result["error"]; exists {
		errorMap := errorField.(map[string]interface{})
		return nil, errors.New(fmt.Sprintf("API error: code %d, message %s", int(errorMap["code"].(float64)), errorMap["message"].(string)))
	}

	var deribitResponse map[string]interface{}
	responseData, _ := json.Marshal(result["result"])
	if err := json.Unmarshal(responseData, &deribitResponse); err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to deserialize response: %v", err))
	}

	return deribitResponse, nil
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

	result, err := d.request("public/get_order_book", params, false)
	if err != nil {
		return restmodels.OrderBook{}, err
	}

	orderBook := restmodels.OrderBook{}
	orderBookData, _ := json.Marshal(result)
	if err := json.Unmarshal(orderBookData, &orderBook); err != nil {
		return restmodels.OrderBook{}, err
	}

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
