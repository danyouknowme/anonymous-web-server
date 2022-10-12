package model

type Benefit struct {
	ResourceName string  `json:"resource_name" bson:"resource_name"`
	Price        float64 `json:"price" bson:"price"`
}
