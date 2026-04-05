package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xprasetio/coffee-pos/internal/entity"
	"github.com/xprasetio/coffee-pos/internal/service"
	"github.com/xprasetio/coffee-pos/pkg/response"
	"github.com/xprasetio/coffee-pos/pkg/validator"
)

// CategoryHandler handles category HTTP requests
type CategoryHandler struct {
	categoryService service.CategoryService
	validator       *validator.Validator
}

// NewCategoryHandler creates a new CategoryHandler instance
func NewCategoryHandler(categoryService service.CategoryService, v *validator.Validator) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
		validator:       v,
	}
}

// FindAll handles GET /categories
func (h *CategoryHandler) FindAll(c *gin.Context) {
	categories, err := h.categoryService.FindAll(c.Request.Context())
	if err != nil {
		response.InternalError(c, "Gagal mengambil data kategori")
		return
	}

	response.OK(c, "Berhasil", categories)
}

// FindByID handles GET /categories/:id
func (h *CategoryHandler) FindByID(c *gin.Context) {
	id := c.Param("id")

	category, err := h.categoryService.FindByID(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.OK(c, "Berhasil", category)
}

// Create handles POST /categories
func (h *CategoryHandler) Create(c *gin.Context) {
	var req entity.CreateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Format request tidak valid")
		return
	}

	if errors := h.validator.Validate(&req); errors != nil {
		response.ValidationError(c, errors)
		return
	}

	category, err := h.categoryService.Create(c.Request.Context(), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, "Kategori berhasil dibuat", category)
}

// Update handles PUT /categories/:id
func (h *CategoryHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req entity.UpdateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Format request tidak valid")
		return
	}

	if errors := h.validator.Validate(&req); errors != nil {
		response.ValidationError(c, errors)
		return
	}

	category, err := h.categoryService.Update(c.Request.Context(), id, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, "Kategori berhasil diupdate", category)
}

// Delete handles DELETE /categories/:id
func (h *CategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.categoryService.Delete(c.Request.Context(), id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, "Kategori berhasil dihapus", nil)
}
