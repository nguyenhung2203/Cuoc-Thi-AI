package response

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"-"`
}

// AuthResponse — trả về cho Login/Register theo API_SPEC (gồm user info + access_token)
type AuthResponse struct {
	User        AuthUser `json:"user"`
	AccessToken string   `json:"access_token"`
}

type AuthUser struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

// LogoutAllResponse — trả về số sessions đã bị revoke
type LogoutAllResponse struct {
	Message         string `json:"message"`
	RevokedSessions int    `json:"revoked_sessions"`
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
