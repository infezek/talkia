package usecase_community

type ListInput struct {
	Location   string `json:"location"`
	CategoryID string `json:"category"`
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
}
