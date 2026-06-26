package response

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserMeResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	Role      string `json:"role"`
	AvatarURL string `json:"avatar_url"`
	Companies []CompanyRole `json:"companies"`
}

type CompanyRole struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}
