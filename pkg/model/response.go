package model

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
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

type AllResourceResponse struct {
	Name      string `json:"name"`
	Label     string `json:"label"`
	Thumbnail string `json:"thumbnail"`
	Plan      []Plan `json:"plan"`
	IsPublish bool   `json:"is_publish"`
}

type GetResourceByNameResponse struct {
	IsPublish   bool                          `json:"is_publish" bson:"is_publish"`
	Name        string                        `json:"name" bson:"name"`
	Label       string                        `json:"label" bson:"label"`
	Description string                        `json:"description" bson:"description"`
	Document    string                        `json:"document" bson:"document"`
	Video       string                        `json:"video" bson:"video"`
	Thumbnail   string                        `json:"thumbnail" bson:"thumbnail"`
	Images      []string                      `json:"images" bson:"images"`
	Plan        []Plan                        `json:"plan" bson:"plan"`
	PatchNotes  []GetResourceByNamePatchNotes `json:"patch_notes" bson:"patch_notes"`
}

type GetResourceByNamePatchNotes struct {
	Version string   `json:"version" bson:"version"`
	Logs    []string `json:"logs" bson:"logs"`
}

type GetDownloadResourceResponse struct {
	Version  string `json:"version" bson:"version"`
	Download string `json:"download" bson:"download"`
}

type GetCounterStateResponse struct {
	Downloads int64 `json:"download" bson:"download"`
	Users     int64 `json:"users" bson:"users"`
	Orders    int64 `json:"orders" bson:"orders"`
}

type GetUserDataResponse struct {
	FirstName  string         `json:"firstname" bson:"firstname"`
	LastName   string         `json:"lastname" bson:"lastname"`
	Email      string         `json:"email" bson:"email"`
	Phone      string         `json:"phone" bson:"phone"`
	Username   string         `json:"username" bson:"username"`
	License    string         `json:"license" bson:"license"`
	Resources  []UserResource `json:"resources" bson:"resources"`
	SecretCode []string       `json:"secret_code" bson:"secret_code"`
}
