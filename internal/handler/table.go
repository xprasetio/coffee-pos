package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xprasetio/coffee-pos/internal/entity"
	"github.com/xprasetio/coffee-pos/internal/service"
	"github.com/xprasetio/coffee-pos/pkg/response"
	"github.com/xprasetio/coffee-pos/pkg/validator"
)

// TableHandler handles table HTTP requests
type TableHandler struct {
	tableService *service.TableService
	validator    *validator.Validator
}

// NewTableHandler creates a new TableHandler instance
func NewTableHandler(tableService *service.TableService, v *validator.Validator) *TableHandler {
	return &TableHandler{
		tableService: tableService,
		validator:    v,
	}
}

// FindAll handles GET /tables
func (h *TableHandler) FindAll(c *gin.Context) {
	tables, err := h.tableService.FindAll(c.Request.Context())
	if err != nil {
		response.InternalError(c, "Gagal mengambil data meja")
		return
	}

	response.OK(c, "Berhasil", tables)
}

// Create handles POST /tables
func (h *TableHandler) Create(c *gin.Context) {
	var req entity.CreateTableRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Format request tidak valid")
		return
	}

	if errors := h.validator.Validate(&req); errors != nil {
		response.ValidationError(c, errors)
		return
	}

	table, err := h.tableService.Create(c.Request.Context(), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, "Meja berhasil dibuat", table)
}

// Update handles PUT /tables/:id
func (h *TableHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req entity.UpdateTableRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Format request tidak valid")
		return
	}

	if errors := h.validator.Validate(&req); errors != nil {
		response.ValidationError(c, errors)
		return
	}

	table, err := h.tableService.Update(c.Request.Context(), id, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, "Meja berhasil diupdate", table)
}

// Delete handles DELETE /tables/:id
func (h *TableHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.tableService.Delete(c.Request.Context(), id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, "Meja berhasil dihapus", nil)
}
