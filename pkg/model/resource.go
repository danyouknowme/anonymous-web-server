package model

type Resource struct {
	IsPublish   bool        `json:"is_publish" bson:"is_publish"`
	Name        string      `json:"name" bson:"name"`
	Label       string      `json:"label" bson:"label"`
	Description string      `json:"description" bson:"description"`
	Document    string      `json:"document" bson:"document"`
	Video       string      `json:"video" bson:"video"`
	Thumbnail   string      `json:"thumbnail" bson:"thumbnail"`
	Images      []string    `json:"images" bson:"images"`
	Plan        []Plan      `json:"plan" bson:"plan"`
	PatchNotes  []PatchNote `json:"patch_notes" bson:"patch_notes"`
}

type Plan struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type PatchNote struct {
	Version  string   `json:"version"`
	Download string   `json:"download"`
	Logs     []string `json:"logs"`
}
