package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/latif-ecommerce-microservices/user-service/internal/dto"

	userpb "github.com/latif-ecommerce-microservices/user-service/pkg/pb/user"

	"github.com/latif-ecommerce-microservices/user-service/internal/service/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	userpb.UnimplementedUserServiceServer
	userService user.ServiceProvider
}

func NewUserHandler(userService user.ServiceProvider) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.UserResponse, error) {

	payload := dto.CreateUserRequest{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	res, err := h.userService.CreateUser(ctx, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &userpb.UserResponse{
		Id:    res.ID.String(),
		Email: res.Email,
		Name:  res.Name,
	}, nil
}

func (h *UserHandler) GetUserByID(ctx context.Context, req *userpb.GetUserByIDRequest) (*userpb.UserResponse, error) {

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user id format")
	}

	res, err := h.userService.GetByID(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &userpb.UserResponse{
		Id:    res.ID.String(),
		Email: res.Email,
		Name:  res.Name,
	}, nil
}

func (h *UserHandler) GetAllUsers(ctx context.Context, id uuid.UUID) (*userpb.ListUserResponse, error) {
	res, err := h.userService.GetAllUsers(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	resp := &userpb.ListUserResponse{
		Users: make([]*userpb.UserResponse, 0, len(res)),
	}

	for _, u := range res {
		resp.Users = append(resp.Users, &userpb.UserResponse{
			Id:    u.ID.String(),
			Email: u.Email,
			Name:  u.Name,
		})
	}

	return resp, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UserResponse, error) {

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user id format")
	}

	var email *string
	if req.Email != "" {
		email = &req.Email
	}

	var name *string
	if req.Name != "" {
		name = &req.Name
	}

	payload := dto.UpdateUserRequest{
		Email: email,
		Name:  name,
	}

	res, err := h.userService.UpdateUser(ctx, id, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &userpb.UserResponse{
		Id:    res.ID.String(),
		Email: res.Email,
		Name:  res.Name,
	}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*userpb.DeleteUserResponse, error) {

	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user id format")
	}

	err = h.userService.DeleteUser(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &userpb.DeleteUserResponse{
		Success: true,
	}, nil
}
