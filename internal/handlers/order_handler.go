package handlers

import (
	"net/http"

	"github.com/ashil-poojary/gostruct/internal/repository"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderRepo *repository.OrderRepository
}

func NewOrderHandler(orderRepo *repository.OrderRepository) *OrderHandler {
	return &OrderHandler{orderRepo: orderRepo}
}

func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	orders, err := h.orderRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}
