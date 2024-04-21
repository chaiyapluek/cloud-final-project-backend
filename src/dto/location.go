package dto

type LocationResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type SingleLocationResponse struct {
	Id    string          `json:"id"`
	Name  string          `json:"name"`
	Menus []*MenuResponse `json:"menus"`
}

type MenuResponse struct {
	Id             string          `json:"id"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	Price          int             `json:"price"`
	IconImage      string          `json:"iconImage"`
	ThumbnailImage string          `json:"thumbnailImage"`
	Steps          []*StepResponse `json:"steps"`
}

type StepResponse struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Type        string            `json:"type"`
	Required    bool              `json:"required"`
	Min         int               `json:"min"`
	Max         int               `json:"max"`
	Options     []*OptionResponse `json:"options"`
}

type OptionResponse struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Price int    `json:"price"`
}
