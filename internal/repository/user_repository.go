package repository

import (
    "errors"

    "gorm.io/gorm"

    "github.com/alexanderbs3/user-orders-api/internal/model"
    apperrors "github.com/alexanderbs3/user-orders-api/pkg/errors"
)

// UserRepository define o contrato — qualquer implementação deve satisfazer essa interface
// Isso permite trocar o banco de dados sem alterar o Service (princípio da inversão de dependência)
type UserRepository interface {
    Create(user *model.User) error
    FindAll(page, limit int) ([]model.User, int64, error)
    FindByID(id uint) (*model.User, error)
    FindByEmail(email string) (*model.User, error)
    Update(user *model.User) error
    Delete(id uint) error
}

// userRepositoryImpl é a implementação concreta usando GORM + PostgreSQL
type userRepositoryImpl struct {
    db *gorm.DB
}

// NewUserRepository é o construtor — padrão de injeção de dependência em Go
func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(user *model.User) error {
    result := r.db.Create(user)
    return result.Error
}

func (r *userRepositoryImpl) FindAll(page, limit int) ([]model.User, int64, error) {
    var users []model.User
    var total int64

    offset := (page - 1) * limit

    // Conta o total para paginação
    r.db.Model(&model.User{}).Count(&total)

    result := r.db.Offset(offset).Limit(limit).Find(&users)
    return users, total, result.Error
}

func (r *userRepositoryImpl) FindByID(id uint) (*model.User, error) {
    var user model.User
    result := r.db.First(&user, id)

    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        return nil, apperrors.NotFound("user")
    }
    return &user, result.Error
}

func (r *userRepositoryImpl) FindByEmail(email string) (*model.User, error) {
    var user model.User
    result := r.db.Where("email = ?", email).First(&user)

    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        return nil, nil // email não encontrado = não é erro aqui
    }
    return &user, result.Error
}

func (r *userRepositoryImpl) Update(user *model.User) error {
    return r.db.Save(user).Error
}

func (r *userRepositoryImpl) Delete(id uint) error {
    result := r.db.Delete(&model.User{}, id)
    if result.RowsAffected == 0 {
        return apperrors.NotFound("user")
    }
    return result.Error
}