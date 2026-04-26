package repository

import (
    "errors"

    "gorm.io/gorm"

    "github.com/alexanderbs3/user-orders-api/internal/model"
    apperrors "github.com/alexanderbs3/user-orders-api/pkg/errors"
)

type OrderRepository interface {
    Create(order *model.Order) error
    FindAll(page, limit int) ([]model.Order, int64, error)
    FindByID(id uint) (*model.Order, error)
    FindByUserID(userID uint, page, limit int) ([]model.Order, int64, error)
    Delete(id uint) error
}

type orderRepositoryImpl struct {
    db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
    return &orderRepositoryImpl{db: db}
}

func (r *orderRepositoryImpl) Create(order *model.Order) error {
    return r.db.Create(order).Error
}

func (r *orderRepositoryImpl) FindAll(page, limit int) ([]model.Order, int64, error) {
    var orders []model.Order
    var total int64

    offset := (page - 1) * limit
    r.db.Model(&model.Order{}).Count(&total)

    // Preload("User") é o equivalente ao @ManyToOne com fetch EAGER — carrega o usuário junto
    result := r.db.Preload("User").Offset(offset).Limit(limit).Find(&orders)
    return orders, total, result.Error
}

func (r *orderRepositoryImpl) FindByID(id uint) (*model.Order, error) {
    var order model.Order
    result := r.db.Preload("User").First(&order, id)

    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        return nil, apperrors.NotFound("order")
    }
    return &order, result.Error
}

func (r *orderRepositoryImpl) FindByUserID(userID uint, page, limit int) ([]model.Order, int64, error) {
    var orders []model.Order
    var total int64

    offset := (page - 1) * limit
    r.db.Model(&model.Order{}).Where("user_id = ?", userID).Count(&total)

    result := r.db.Where("user_id = ?", userID).Offset(offset).Limit(limit).Find(&orders)
    return orders, total, result.Error
}

func (r *orderRepositoryImpl) Delete(id uint) error {
    result := r.db.Delete(&model.Order{}, id)
    if result.RowsAffected == 0 {
        return apperrors.NotFound("order")
    }
    return result.Error
}