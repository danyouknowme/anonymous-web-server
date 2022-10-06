package model

type Download struct {
	Username     string `json:"username" bson:"username"`
	ResourceName string `json:"resource_name" bson:"resource_name"`
	Date         string `json:"date" bson:"date"`
}
