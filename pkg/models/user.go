package models

type User struct {
	FirstName  string         `json:"firstName"`
	LastName   string         `json:"lastName"`
	Email      string         `json:"email"`
	Phone      string         `json:"phone"`
	Username   string         `json:"username"`
	Password   string         `json:"password"`
	IsAdmin    bool           `json:"isAdmin"`
	License    string         `json:"license"`
	Resources  []UserResource `json:"resources"`
	LastReset  string         `json:"lastReset"`
	ResetTime  int64          `json:"resetTime"`
	SecretCode []string       `json:"secretCode"`
}

type UserResource struct {
	Name         string `json:"name"`
	DownloadLink string `json:"downloadLink"`
	Status       string `json:"status"`
	DayLeft      int64  `json:"dayLeft"`
}
