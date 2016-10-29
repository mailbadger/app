package pagination

var DefaultPerPage uint = 10

type Pagination struct {
	Page       uint          `json:"page"`
	Offset     uint          `json:"offset"`
	PerPage    uint          `json:"per_page"`
	Total      uint64        `json:"total"`
	Collection []interface{} `json:"collection"`
}

func (pagination *Pagination) SetTotal(total uint64) {
	pagination.Total = total
}

func (pagination *Pagination) Append(obj interface{}) {
	pagination.Collection = append(pagination.Collection, obj)
}
