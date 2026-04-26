package dto

import "github.com/alexanderbs3/user-orders-api/internal/model"

type CreateOrderRequest struct {
    UserID      uint                `json:"user_id" binding:"required"`
    Description string              `json:"description" binding:"required,min=3"`
    Amount      float64             `json:"amount" binding:"required,gt=0"`
    Status      model.OrderStatus   `json:"status" binding:"omitempty,oneof=pending paid canceled"`
}

type OrderResponse struct {
    ID          uint                `json:"id"`
    UserID      uint                `json:"user_id"`
    Description string              `json:"description"`
    Amount      float64             `json:"amount"`
    Status      model.OrderStatus   `json:"status"`
    CreatedAt   string              `json:"created_at"`
}

// PaginationParams é reutilizado em qualquer listagem paginada
type PaginationParams struct {
    Page  int `form:"page,default=1" binding:"min=1"`
    Limit int `form:"limit,default=10" binding:"min=1,max=100"`
}