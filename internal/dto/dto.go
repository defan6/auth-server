package dto

type TokenGenerateRequest struct {
	UserID int64
	Email  string
	AppID  int64
}

func NewTokenGenerateRequest(userID int64, email string, appID int64) *TokenGenerateRequest {
	return &TokenGenerateRequest{
		UserID: userID,
		Email:  email,
		AppID:  appID,
	}
}

type TokenGenerateResponse struct {
	Token string
}

func NewTokenGenerateResponse(token string) *TokenGenerateResponse {
	return &TokenGenerateResponse{
		Token: token,
	}
}

type ListUserRequest struct {
	Filters map[string]string
}

func NewListUserRequest(filters map[string]string) *ListUserRequest {
	return &ListUserRequest{
		Filters: filters,
	}
}

type ListUserResponse struct {
	Users []*UserResponse
}

func NewListUserResponse(users []*UserResponse) *ListUserResponse {
	return &ListUserResponse{
		Users: users,
	}
}

type UserResponse struct {
	ID    int64
	Email string
	Role  string
}

func NewUserResponse(id int64, email string, role string) *UserResponse {
	return &UserResponse{
		ID:    id,
		Email: email,
		Role:  role,
	}
}
