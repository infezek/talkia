package entity

type Pagination struct {
	Data    interface{} `json:"data"`
	PerPage int32       `json:"per_page"`
	Offset  int32       `json:"-"`
	Page    int32       `json:"page"`
	Total   int64       `json:"total,omitempty"`
}

func NewPagination(perPage, page int32) *Pagination {
	return &Pagination{
		PerPage: perPage,
		Page:    page,
	}
}

func (p *Pagination) SetTotal(total int64) {
	p.Total = total
}

func (p *Pagination) SetData(data interface{}) {
	p.Data = data
}
