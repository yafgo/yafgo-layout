package service

import (
	"context"
	"yafgo/yafgo-layout/internal/model"
	"yafgo/yafgo-layout/internal/repository"
	"yafgo/yafgo-layout/pkg/hash"

	"github.com/pkg/errors"
)

// ReqRegisterUsername 用户名注册
type ReqRegisterUsername struct {
	Username   string `json:"username,omitempty" binding:"required"`
	Password   string `json:"password,omitempty" binding:"required"`
	VerifyCode string `json:"verify_code,omitempty" binding:"required"`
}

// ReqLoginUsername 用户名登录
type ReqLoginUsername struct {
	Username string `json:"username,omitempty" binding:"required"`
	Password string `json:"password,omitempty" binding:"required"`
}

type UserService interface {
	RegisterByUsername(ctx context.Context, req *ReqRegisterUsername) (*model.User, error)
	LoginByUsername(ctx context.Context, req *ReqLoginUsername) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
}

type userService struct {
	userRepo repository.UserRepository
	*Service
}

func NewUserService(service *Service, userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
		Service:  service,
	}
}

// RegisterByUsername implements UserService.
func (s *userService) RegisterByUsername(ctx context.Context, req *ReqRegisterUsername) (*model.User, error) {
	user := &model.User{
		Username: req.Username,
		Password: req.Password,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.Wrap(err, "创建用户失败")
	}

	return user, nil
}

// LoginByUsername implements UserService.
func (s *userService) LoginByUsername(ctx context.Context, req *ReqLoginUsername) (*model.User, error) {
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.Wrap(err, "用户不存在")
	}

	if !hash.BcryptCheck(req.Password, user.Password) {
		return nil, errors.New("密码不正确")
	}

	return user, nil
}

// GetByID implements UserService.
func (s *userService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	return s.userRepo.GetByID(ctx, id)
}
