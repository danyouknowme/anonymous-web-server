package model

type Order struct {
	BillNumber       string          `json:"billno" bson:"billno"`
	Username         string          `json:"username" bson:"username"`
	Resources        []OrderResource `json:"resources" bson:"resources"`
	TransactionImage string          `json:"transaction_image" bson:"transaction_image"`
	Date             string          `json:"date" bson:"date"`
	TotalPrice       float64         `json:"total_price" bson:"total_price"`
	Status           string          `json:"status" bson:"status"`
}

type OrderResource struct {
	ResourceName  string  `json:"resource_name" bson:"resource_name"`
	ResourceLabel string  `json:"resource_label" bson:"resource_label"`
	Plan          string  `json:"plan" bson:"plan"`
	Price         float64 `json:"price" bson:"price"`
}
