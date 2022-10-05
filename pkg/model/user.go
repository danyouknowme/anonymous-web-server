package model

type User struct {
	FirstName  string         `json:"firstname" bson:"firstname"`
	LastName   string         `json:"lastname" bson:"lastname"`
	Email      string         `json:"email" bson:"email"`
	Phone      string         `json:"phone" bson:"phone"`
	Username   string         `json:"username" bson:"username"`
	Password   string         `json:"password" bson:"password"`
	IsAdmin    bool           `json:"is_admin" bson:"is_admin"`
	License    string         `json:"license" bson:"license"`
	Resources  []UserResource `json:"resources" bson:"resources"`
	LastReset  string         `json:"last_reset" bson:"last_reset"`
	ResetTime  int64          `json:"reset_time" bson:"reset_time"`
	SecretCode []string       `json:"secret_code" bson:"secret_code"`
}

type UserResource struct {
	Name    string  `json:"name"`
	Status  *string `json:"status"`
	DayLeft int64   `json:"dayLeft"`
}
