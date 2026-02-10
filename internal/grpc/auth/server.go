package auth

import (
	"context"

	ssov1 "github.com/defan6/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	Register(ctx context.Context, email string, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}
type serverAPI struct {
	ssov1.UnimplementedAuthServer
	authService AuthService
}

func Register(gRPC *grpc.Server, authService AuthService) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{authService: authService})
}

func (s *serverAPI) Login(ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}
	token, err := s.authService.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		// TODO ...
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}
	userID, err := s.authService.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		// TODO ...
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &ssov1.RegisterResponse{UserId: userID}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context,
	req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}
	isAdmin, err := s.authService.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		// TODO ...
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}

func validateLogin(req *ssov1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	if req.GetAppId() == 0 {
		return status.Error(codes.InvalidArgument, "appId is required")
	}
	return nil
}

func validateRegister(req *ssov1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}
	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}
	return nil
}

func validateIsAdmin(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() <= 0 {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}
	return nil
}
