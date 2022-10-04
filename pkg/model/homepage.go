package model

type Homepage struct {
	Title         string   `json:"title"`
	ResourceName  string   `json:"resourceName"`
	ResourceLabel string   `json:"resourceLabel"`
	Description   []string `json:"description"`
}
