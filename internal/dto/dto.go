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
