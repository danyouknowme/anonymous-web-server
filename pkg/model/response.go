package model

type ErrorResponse struct {
	Message string `json:"message"`
}

type UserResponse struct {
	FirstName  string         `json:"firstName"`
	LastName   string         `json:"lastName"`
	Email      string         `json:"email"`
	Phone      string         `json:"phone"`
	Username   string         `json:"username"`
	IsAdmin    bool           `json:"isAdmin"`
	License    string         `json:"license"`
	Resources  []UserResource `json:"resources"`
	LastReset  string         `json:"lastReset"`
	ResetTime  int64          `json:"resetTime"`
	SecretCode []string       `json:"secretCode"`
}

type LoginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}
