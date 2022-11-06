package user

import (
	"context"
	"integration-test/internal/entity"

	"gorm.io/gorm"
)

type User struct {
	ID   string
	Name string
}

type UserRepo struct {
	db *gorm.DB
}

func New(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (ur UserRepo) CreateUser(ctx context.Context, input entity.User) (entity.User, error) {
	var row User

	row.ID = input.ID
	row.Name = input.Name

	err := ur.db.WithContext(ctx).Create(&row).Error
	// fmt.Println("err di repo", err)
	if err != nil {
		return entity.User{}, err
	}

	retUser := entity.User{
		ID:   row.ID,
		Name: row.Name,
	}
	// fmt.Println("err di repo", retUser)

	return retUser, nil
}
