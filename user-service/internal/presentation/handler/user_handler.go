package handler

import (
	"context"
	"time"

	"github.com/AnnaKhairetdinova/user-service/docs/proto/user"
	svc "github.com/AnnaKhairetdinova/user-service/internal/app/service"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCHandler struct {
	user.UnimplementedUserServiceServer
	us svc.UserService
}

func NewGRPCHandler(us svc.UserService) *GRPCHandler {
	return &GRPCHandler{us: us}
}

func (grpc *GRPCHandler) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {
	createdUser, err := grpc.us.CreateUser(ctx, req.Name, req.Email)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userMsg := &user.UserMessage{
		Uuid:      createdUser.UUID.String(),
		Name:      createdUser.Name,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt.Format(time.RFC3339),
	}

	return &user.CreateUserResponse{User: userMsg}, nil
}

func (grpc *GRPCHandler) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.GetUserResponse, error) {
	parseUUID, err := uuid.Parse(req.Uuid)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid uuid")
	}

	getUser, err := grpc.us.GetUserByUUID(ctx, parseUUID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	userMsg := &user.UserMessage{
		Uuid:      getUser.UUID.String(),
		Name:      getUser.Name,
		Email:     getUser.Email,
		CreatedAt: getUser.CreatedAt.Format(time.RFC3339),
	}

	return &user.GetUserResponse{User: userMsg}, nil
}

func (grpc *GRPCHandler) ListUsers(ctx context.Context, _ *user.ListUsersRequest) (*user.ListUsersResponse, error) {
	users, err := grpc.us.ListUsers(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	userMsg := make([]*user.UserMessage, 0, len(users))
	for _, userItem := range users {
		userMsg = append(userMsg, &user.UserMessage{
			Uuid:      userItem.UUID.String(),
			Name:      userItem.Name,
			Email:     userItem.Email,
			CreatedAt: userItem.CreatedAt.Format(time.RFC3339),
		})
	}

	return &user.ListUsersResponse{Users: userMsg}, nil
}
