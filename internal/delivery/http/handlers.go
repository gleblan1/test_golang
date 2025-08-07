package http

import (
	"net/http"
	"strconv"

	"crypto-price-tracker-app/internal/application/dto"
	"crypto-price-tracker-app/internal/application/services"

	"github.com/gin-gonic/gin"
)

// CurrencyHandler обрабатывает HTTP запросы для работы с криптовалютами
type CurrencyHandler struct {
	currencyService *services.CurrencyService
}

// NewCurrencyHandler создает новый экземпляр CurrencyHandler
func NewCurrencyHandler(currencyService *services.CurrencyService) *CurrencyHandler {
	return &CurrencyHandler{
		currencyService: currencyService,
	}
}

// AddCurrency godoc
// @Summary Добавить криптовалюту в отслеживание
// @Description Добавляет новую криптовалюту для отслеживания цен
// @Tags currency
// @Accept json
// @Produce json
// @Param request body dto.AddCurrencyRequest true "Данные криптовалюты"
// @Success 201 {object} dto.CurrencyResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /currency/add [post]
func (h *CurrencyHandler) AddCurrency(c *gin.Context) {
	var req dto.AddCurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	currency, err := h.currencyService.AddCurrency(c.Request.Context(), &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "currency already exists" {
			status = http.StatusConflict
		}
		c.JSON(status, dto.ErrorResponse{
			Error:   "currency_error",
			Message: err.Error(),
			Code:    status,
		})
		return
	}

	c.JSON(http.StatusCreated, currency)
}

// RemoveCurrency godoc
// @Summary Удалить криптовалюту из отслеживания
// @Description Удаляет криптовалюту из списка отслеживаемых
// @Tags currency
// @Accept json
// @Produce json
// @Param request body dto.RemoveCurrencyRequest true "Данные для удаления"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /currency/remove [post]
func (h *CurrencyHandler) RemoveCurrency(c *gin.Context) {
	var req dto.RemoveCurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	err := h.currencyService.RemoveCurrency(c.Request.Context(), &req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "currency not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, dto.ErrorResponse{
			Error:   "currency_error",
			Message: err.Error(),
			Code:    status,
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Currency removed successfully",
	})
}

// GetPrice godoc
// @Summary Получить цену криптовалюты
// @Description Возвращает цену криптовалюты в указанное время
// @Tags currency
// @Accept json
// @Produce json
// @Param coin query string true "Символ криптовалюты"
// @Param timestamp query int true "Unix timestamp"
// @Success 200 {object} dto.PriceResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /currency/price [get]
func (h *CurrencyHandler) GetPrice(c *gin.Context) {
	coin := c.Query("coin")
	timestampStr := c.Query("timestamp")

	if coin == "" || timestampStr == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "coin and timestamp parameters are required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "invalid timestamp format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	req := &dto.GetPriceRequest{
		Coin:      coin,
		Timestamp: timestamp,
	}

	price, err := h.currencyService.GetPrice(c.Request.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "currency not found" || err.Error() == "price not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, dto.ErrorResponse{
			Error:   "price_error",
			Message: err.Error(),
			Code:    status,
		})
		return
	}

	c.JSON(http.StatusOK, price)
}

// GetAllCurrencies godoc
// @Summary Получить все активные криптовалюты
// @Description Возвращает список всех активных криптовалют
// @Tags currency
// @Accept json
// @Produce json
// @Success 200 {array} dto.CurrencyResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /currency/list [get]
func (h *CurrencyHandler) GetAllCurrencies(c *gin.Context) {
	currencies, err := h.currencyService.GetAllActiveCurrencies(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "currency_error",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, currencies)
}

// HealthCheck godoc
// @Summary Проверка здоровья сервиса
// @Description Возвращает статус сервиса
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} dto.SuccessResponse
// @Router /health [get]
func (h *CurrencyHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Service is healthy",
	})
}
