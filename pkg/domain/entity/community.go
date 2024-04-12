package entity

type Community struct {
	Location   string     `json:"location"`
	Trends     []Bot      `json:"trends"`
	Categories []Category `json:"category"`
}

func NewCommunity(location string, trends []Bot, categories []Category) *Community {
	return &Community{
		Location:   location,
		Trends:     trends,
		Categories: categories,
	}
}
