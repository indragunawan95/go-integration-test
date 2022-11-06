package user

import (
	"context"
	"integration-test/internal/entity"
	pb "integration-test/proto/integration-test/user"

	"google.golang.org/grpc/status"
)

type UserUsecaseItf interface {
	CreateUser(context.Context, entity.User) (entity.User, error)
}

type User struct {
	pb.UnimplementedUserServiceServer
	uc UserUsecaseItf
}

func New(uc UserUsecaseItf) *User {
	return &User{
		uc: uc,
	}
}

func (h *User) CreateUser(ctx context.Context, user *pb.User) (*pb.User, error) {
	entUser := entity.User{
		ID:   user.Id,
		Name: user.Name,
	}
	newUser, err := h.uc.CreateUser(ctx, entUser)
	if err != nil {
		return &pb.User{}, status.Error(status.Code(err), "failed to create user")
	}

	pbUser := &pb.User{
		Id:   newUser.ID,
		Name: newUser.Name,
	}

	return pbUser, nil
}
