package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xprasetio/coffee-pos/internal/entity"
	"github.com/xprasetio/coffee-pos/internal/service"
	"github.com/xprasetio/coffee-pos/pkg/response"
	"github.com/xprasetio/coffee-pos/pkg/validator"
)

// StockHandler handles stock HTTP requests
type StockHandler struct {
	stockService service.StockService
	validator    *validator.Validator
}

// NewStockHandler creates a new StockHandler instance
func NewStockHandler(stockService service.StockService, v *validator.Validator) *StockHandler {
	return &StockHandler{
		stockService: stockService,
		validator:    v,
	}
}

// GetStock handles GET /stock/:id
func (h *StockHandler) GetStock(c *gin.Context) {
	id := c.Param("id")

	product, err := h.stockService.GetStock(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.OK(c, "Berhasil", gin.H{
		"stock":      product.Stock,
		"product_id": product.ID,
		"name":       product.Name,
	})
}

// Adjust handles POST /stock/:id/adjust
func (h *StockHandler) Adjust(c *gin.Context) {
	productID := c.Param("id")
	userID := c.GetString("user_id")

	var req entity.StockAdjustmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Format request tidak valid")
		return
	}

	if errors := h.validator.Validate(&req); errors != nil {
		response.ValidationError(c, errors)
		return
	}

	if err := h.stockService.Adjust(c.Request.Context(), productID, userID, req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, "Stok berhasil diupdate", nil)
}
