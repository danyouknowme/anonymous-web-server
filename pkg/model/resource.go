package model

type Resource struct {
	IsPublish   bool        `json:"is_publish"`
	Name        string      `json:"name"`
	Label       string      `json:"label"`
	Description string      `json:"description"`
	Document    string      `json:"document"`
	Video       string      `json:"video"`
	Thumbnail   string      `json:"thumbnail"`
	Images      []string    `json:"images"`
	Plan        []Plan      `json:"plan"`
	PatchNotes  []PatchNote `json:"patch_notes"`
}

type Plan struct {
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

type PatchNote struct {
	Version  string   `json:"version"`
	Download string   `json:"download"`
	Logs     []string `json:"logs"`
}
