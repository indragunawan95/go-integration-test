package user

import (
	"context"
	"integration-test/internal/entity"
)

type UserResourceItf interface {
	CreateUser(ctx context.Context, input entity.User) (entity.User, error)
}

type Usecase struct {
	user UserResourceItf
}

func New(user UserResourceItf) *Usecase {
	return &Usecase{
		user: user,
	}
}

func (uc Usecase) CreateUser(ctx context.Context, input entity.User) (entity.User, error) {
	newUser, err := uc.user.CreateUser(ctx, input)
	if err != nil {
		return entity.User{}, err
	}

	return newUser, nil
}
