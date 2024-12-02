package handler

import (
	"fravega-tech/internal/domain"
	"fravega-tech/internal/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProductHandler struct {
	usecase *usecase.ProductUsecase
}

func NewProductHandler(router *gin.Engine, usecase *usecase.ProductUsecase) {
	handler := &ProductHandler{usecase: usecase}

	router.Use(cors.Default())

	v1 := router.Group("/api/v1/products")
	v1.POST("", handler.CreateProduct)
	v1.GET("", handler.GetProducts)
	v1.PUT("/:id", handler.UpdateProduct)
	v1.DELETE("", handler.DeleteProduct)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.CreateProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	name := c.DefaultQuery("name", "")
	category := c.DefaultQuery("category", "")

	products, err := h.usecase.GetProducts(name, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.UpdateProduct(id, &product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	ids := c.QueryArray("ids")

	if len(ids) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No product ids provided"})
		return
	}

	if err := h.usecase.DeleteProducts(ids); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
