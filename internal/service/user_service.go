package service

import (
    "github.com/alexanderbs3/user-orders-api/internal/dto"
    "github.com/alexanderbs3/user-orders-api/internal/model"
    "github.com/alexanderbs3/user-orders-api/internal/repository"
    apperrors "github.com/alexanderbs3/user-orders-api/pkg/errors"
)

type UserService interface {
    Create(req dto.CreateUserRequest) (*model.User, error)
    FindAll(page, limit int) ([]model.User, int64, error)
    FindByID(id uint) (*model.User, error)
    Update(id uint, req dto.UpdateUserRequest) (*model.User, error)
    Delete(id uint) error
}

type userServiceImpl struct {
    repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
    return &userServiceImpl{repo: repo}
}

func (s *userServiceImpl) Create(req dto.CreateUserRequest) (*model.User, error) {
    // Regra de negócio: e-mail deve ser único
    existing, err := s.repo.FindByEmail(req.Email)
    if err != nil {
        return nil, err
    }
    if existing != nil {
        return nil, apperrors.Conflict("email already in use")
    }

    user := &model.User{
        Name:  req.Name,
        Email: req.Email,
    }

    if err := s.repo.Create(user); err != nil {
        return nil, apperrors.Internal("failed to create user")
    }

    return user, nil
}

func (s *userServiceImpl) FindAll(page, limit int) ([]model.User, int64, error) {
    return s.repo.FindAll(page, limit)
}

func (s *userServiceImpl) FindByID(id uint) (*model.User, error) {
    return s.repo.FindByID(id)
}

func (s *userServiceImpl) Update(id uint, req dto.UpdateUserRequest) (*model.User, error) {
    user, err := s.repo.FindByID(id)
    if err != nil {
        return nil, err
    }

    // Atualização parcial: só altera campos que foram enviados (ponteiros não nulos)
    if req.Name != nil {
        user.Name = *req.Name
    }
    if req.Email != nil {
        // Verifica conflito de e-mail ao atualizar
        existing, err := s.repo.FindByEmail(*req.Email)
        if err != nil {
            return nil, err
        }
        if existing != nil && existing.ID != id {
            return nil, apperrors.Conflict("email already in use")
        }
        user.Email = *req.Email
    }

    if err := s.repo.Update(user); err != nil {
        return nil, apperrors.Internal("failed to update user")
    }

    return user, nil
}

func (s *userServiceImpl) Delete(id uint) error {
    return s.repo.Delete(id)
}