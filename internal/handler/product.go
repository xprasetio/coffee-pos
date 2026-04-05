package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xprasetio/coffee-pos/internal/entity"
	"github.com/xprasetio/coffee-pos/internal/repository"
	"github.com/xprasetio/coffee-pos/internal/service"
	"github.com/xprasetio/coffee-pos/pkg/response"
	"github.com/xprasetio/coffee-pos/pkg/validator"
)

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Items interface{} `json:"items"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
}

// ProductHandler handles product HTTP requests
type ProductHandler struct {
	productService service.ProductService
	validator      *validator.Validator
}

// NewProductHandler creates a new ProductHandler instance
func NewProductHandler(productService service.ProductService, v *validator.Validator) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		validator:      v,
	}
}

// FindAll handles GET /products
func (h *ProductHandler) FindAll(c *gin.Context) {
	search := c.Query("search")
	categoryID := c.Query("category_id")
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	filter := repository.ProductFilter{
		Search:     search,
		CategoryID: categoryID,
		Page:       page,
		Limit:      limit,
	}

	if activeStr := c.Query("is_active"); activeStr != "" {
		if active, err := strconv.ParseBool(activeStr); err == nil {
			filter.IsActive = &active
		}
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	products, total, err := h.productService.FindAll(c.Request.Context(), filter)
	if err != nil {
		response.InternalError(c, "Gagal mengambil data produk")
		return
	}

	paginatedResponse := PaginatedResponse{
		Items: products,
		Total: total,
		Page:  filter.Page,
		Limit: filter.Limit,
	}

	response.OK(c, "Berhasil", paginatedResponse)
}

// FindByID handles GET /products/:id
func (h *ProductHandler) FindByID(c *gin.Context) {
	id := c.Param("id")

	product, err := h.productService.FindByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.OK(c, "Berhasil", product)
}

// Create handles POST /products
func (h *ProductHandler) Create(c *gin.Context) {
	var req entity.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Format request tidak valid")
		return
	}

	if errors := h.validator.Validate(&req); errors != nil {
		response.ValidationError(c, errors)
		return
	}

	product, err := h.productService.Create(c.Request.Context(), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, "Produk berhasil dibuat", product)
}

// Update handles PUT /products/:id
func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req entity.UpdateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Format request tidak valid")
		return
	}

	if errors := h.validator.Validate(&req); errors != nil {
		response.ValidationError(c, errors)
		return
	}

	product, err := h.productService.Update(c.Request.Context(), id, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, "Produk berhasil diupdate", product)
}

// Delete handles DELETE /products/:id
func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.productService.Delete(c.Request.Context(), id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, "Produk berhasil dihapus", nil)
}
