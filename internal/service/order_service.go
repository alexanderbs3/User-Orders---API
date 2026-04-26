package service

import (
    "github.com/alexanderbs3/user-orders-api/internal/dto"
    "github.com/alexanderbs3/user-orders-api/internal/model"
    "github.com/alexanderbs3/user-orders-api/internal/repository"
    apperrors "github.com/alexanderbs3/user-orders-api/pkg/errors"
)

type OrderService interface {
    Create(req dto.CreateOrderRequest) (*model.Order, error)
    FindAll(page, limit int) ([]model.Order, int64, error)
    FindByID(id uint) (*model.Order, error)
    FindByUserID(userID uint, page, limit int) ([]model.Order, int64, error)
    Delete(id uint) error
}

type orderServiceImpl struct {
    orderRepo repository.OrderRepository
    userRepo  repository.UserRepository
}

func NewOrderService(orderRepo repository.OrderRepository, userRepo repository.UserRepository) OrderService {
    return &orderServiceImpl{orderRepo: orderRepo, userRepo: userRepo}
}

func (s *orderServiceImpl) Create(req dto.CreateOrderRequest) (*model.Order, error) {
    // Regra: não pode criar pedido para usuário inexistente
    _, err := s.userRepo.FindByID(req.UserID)
    if err != nil {
        return nil, err // propaga o NotFound do repository
    }

    status := req.Status
    if status == "" {
        status = model.StatusPending // valor padrão
    }

    order := &model.Order{
        UserID:      req.UserID,
        Description: req.Description,
        Amount:      req.Amount,
        Status:      status,
    }

    if err := s.orderRepo.Create(order); err != nil {
        return nil, apperrors.Internal("failed to create order")
    }

    return order, nil
}

func (s *orderServiceImpl) FindAll(page, limit int) ([]model.Order, int64, error) {
    return s.orderRepo.FindAll(page, limit)
}

func (s *orderServiceImpl) FindByID(id uint) (*model.Order, error) {
    return s.orderRepo.FindByID(id)
}

func (s *orderServiceImpl) FindByUserID(userID uint, page, limit int) ([]model.Order, int64, error) {
    // Verifica se o usuário existe antes de listar os pedidos
    _, err := s.userRepo.FindByID(userID)
    if err != nil {
        return nil, 0, err
    }
    return s.orderRepo.FindByUserID(userID, page, limit)
}

func (s *orderServiceImpl) Delete(id uint) error {
    return s.orderRepo.Delete(id)
}