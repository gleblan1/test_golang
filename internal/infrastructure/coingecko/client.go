package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client представляет клиент для работы с CoinGecko API
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// PriceResponse представляет ответ от CoinGecko API
type PriceResponse struct {
	ID                           string  `json:"id"`
	Symbol                       string  `json:"symbol"`
	Name                         string  `json:"name"`
	Image                        string  `json:"image"`
	CurrentPrice                 float64 `json:"current_price"`
	MarketCap                    int64   `json:"market_cap"`
	MarketCapRank                int     `json:"market_cap_rank"`
	FullyDilutedValuation        int64   `json:"fully_diluted_valuation"`
	TotalVolume                  int64   `json:"total_volume"`
	High24h                      float64 `json:"high_24h"`
	Low24h                       float64 `json:"low_24h"`
	PriceChange24h               float64 `json:"price_change_24h"`
	PriceChangePercentage24h     float64 `json:"price_change_percentage_24h"`
	MarketCapChange24h           int64   `json:"market_cap_change_24h"`
	MarketCapChangePercentage24h float64 `json:"market_cap_change_percentage_24h"`
	CirculatingSupply            float64 `json:"circulating_supply"`
	TotalSupply                  float64 `json:"total_supply"`
	MaxSupply                    float64 `json:"max_supply"`
	Ath                          float64 `json:"ath"`
	AthChangePercentage          float64 `json:"ath_change_percentage"`
	AthDate                      string  `json:"ath_date"`
	Atl                          float64 `json:"atl"`
	AtlChangePercentage          float64 `json:"atl_change_percentage"`
	AtlDate                      string  `json:"atl_date"`
	Roi                          *ROI    `json:"roi"`
	LastUpdated                  string  `json:"last_updated"`
}

// ROI представляет Return on Investment
type ROI struct {
	Times      float64 `json:"times"`
	Currency   string  `json:"currency"`
	Percentage float64 `json:"percentage"`
}

// NewClient создает новый клиент для CoinGecko API
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetPrice возвращает текущую цену криптовалюты
func (c *Client) GetPrice(ctx context.Context, symbol string) (float64, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=usd", c.baseURL, symbol)

	// Создаем HTTP запрос
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Выполняем запрос
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %w", err)
	}

	// Парсим JSON ответ
	var priceData map[string]map[string]float64
	if err := json.Unmarshal(body, &priceData); err != nil {
		return 0, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Извлекаем цену
	if currencyData, exists := priceData[symbol]; exists {
		if price, exists := currencyData["usd"]; exists {
			return price, nil
		}
	}

	return 0, fmt.Errorf("price not found for symbol: %s", symbol)
}

// GetDetailedPrice возвращает детальную информацию о цене криптовалюты
func (c *Client) GetDetailedPrice(ctx context.Context, symbol string) (*PriceResponse, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("%s/coins/%s", c.baseURL, symbol)

	// Создаем HTTP запрос
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Выполняем запрос
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Парсим JSON ответ
	var priceResponse PriceResponse
	if err := json.Unmarshal(body, &priceResponse); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &priceResponse, nil
}
