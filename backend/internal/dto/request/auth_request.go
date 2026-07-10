package request

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshRequest không cần nữa — refresh token được đọc từ HttpOnly Cookie

type GoogleLoginRequest struct {
	IDToken  string `json:"id_token"`
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"full_name"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
}
