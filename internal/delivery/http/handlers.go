package http

import (
	"net/http"
	"strconv"

	"crypto-price-tracker-app/internal/application/dto"
	"crypto-price-tracker-app/internal/application/services"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	currencyService *services.CurrencyService
	priceService    *services.PriceService
}

func NewHandlers(currencyService *services.CurrencyService, priceService *services.PriceService) *Handlers {
	return &Handlers{
		currencyService: currencyService,
		priceService:    priceService,
	}
}

func (h *Handlers) AddCurrency(c *gin.Context) {
	var req dto.AddCurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
			Code:    400,
		})
		return
	}

	currency, err := h.currencyService.AddCurrency(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "currency_error",
			Message: err.Error(),
			Code:    400,
		})
		return
	}

	c.JSON(http.StatusCreated, currency)
}

func (h *Handlers) RemoveCurrency(c *gin.Context) {
	var req dto.RemoveCurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
			Code:    400,
		})
		return
	}

	if err := h.currencyService.RemoveCurrency(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "currency_error",
			Message: err.Error(),
			Code:    400,
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Currency removed successfully",
	})
}

func (h *Handlers) GetPrice(c *gin.Context) {
	coin := c.Query("coin")
	if coin == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "coin parameter is required",
			Code:    400,
		})
		return
	}

	timestampStr := c.Query("timestamp")
	if timestampStr == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "timestamp parameter is required",
			Code:    400,
		})
		return
	}

	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "invalid timestamp format",
			Code:    400,
		})
		return
	}

	req := &dto.GetPriceRequest{
		Coin:      coin,
		Timestamp: timestamp,
	}

	price, err := h.currencyService.GetPrice(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error:   "price_error",
			Message: err.Error(),
			Code:    404,
		})
		return
	}

	c.JSON(http.StatusOK, price)
}

func (h *Handlers) GetAllCurrencies(c *gin.Context) {
	currencies, err := h.currencyService.GetAllActiveCurrencies(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
			Code:    500,
		})
		return
	}

	c.JSON(http.StatusOK, currencies)
}

func (h *Handlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "crypto-price-tracker",
	})
}
