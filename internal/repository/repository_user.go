package repository

import (
	"context"
	"yafgo/yafgo-layout/internal/model"
	"yafgo/yafgo-layout/pkg/database"

	"github.com/pkg/errors"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int64) (*model.User, error)
}

type userRepository struct {
	*Repository
}

func NewUserRepository(r *Repository) UserRepository {
	return &userRepository{
		Repository: r,
	}
}

// Create implements UserRepository.
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	userDo := r.q.User.WithContext(ctx)
	err := userDo.Create(user)
	if database.IsErrDuplicateEntryCode(err) {
		return errors.Wrap(err, "用户名已存在")
	}
	return err
}

// GetByID implements UserRepository.
func (r *userRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	userDo := r.q.User.WithContext(ctx)
	user, err := userDo.GetByID(id)
	return user, err
}

// Update implements UserRepository.
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	panic("unimplemented")
}
