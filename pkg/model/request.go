package model

type CreateUserRequest struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	Username  string `json:"username" binding:"required,alphanum"`
	Password  string `json:"password" binding:"required,min=6"`
}

type LoginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type AddResourceToUserRequest struct {
	Username string         `json:"username" binding:"required,alphanum"`
	Resource ResourceToUser `json:"resource"`
}

type ResourceToUser struct {
	Name    string `json:"name"`
	DayLeft int64  `json:"day_left"`
}

type DownloadResourceRequest struct {
	ResourceName string `json:"resource_name"`
}

type CheckLicenseRequest struct {
	License      string `json:"license"`
	ResourceName string `json:"resource_name"`
}

type UpdateUserDataRequest struct {
	FirstName string         `json:"firstname" bson:"firstname"`
	LastName  string         `json:"lastname" bson:"lastname"`
	Email     string         `json:"email" bson:"email"`
	Phone     string         `json:"phone" bson:"phone"`
	Username  string         `json:"username" bson:"username"`
	License   string         `json:"license" bson:"license"`
	Resources []UserResource `json:"resources" bson:"resources"`
}

type ForgotPasswordRequest struct {
	Username    string `json:"username" bson:"username"`
	NewPassword string `json:"new_password" bson:"new_password"`
	SecretCode  string `json:"secret_code" bson:"secret_code"`
}

type AddOrderRequest struct {
	Username         string                 `json:"username" bson:"username"`
	RequestOrder     []RequestOrderResource `json:"request_order" bson:"request_order"`
	TransactionImage string                 `json:"transaction_image" bson:"transaction_image"`
}

type RequestOrderResource struct {
	Resource string `json:"resource" bson:"resource"`
	Plan     string `json:"plan" bson:"plan"`
}

type ConfirmOrderRequest struct {
	BillNumber string `json:"billno" bson:"billno"`
	Status     string `json:"status" bson:"status"`
}

type ClearPartnerBenefitRequest struct {
	ResourceName string `json:"resource_name" bson:"resource_name"`
}
