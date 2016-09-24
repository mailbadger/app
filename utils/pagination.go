package utils

var DefaultPerPage uint = 10

type Pagination struct {
	Page       uint
	Offset     uint
	PerPage    uint
	Total      uint64
	Collection []interface{}
}

func (pagination *Pagination) SetTotal(total uint64) {
	pagination.Total = total
}

func (pagination *Pagination) Append(obj interface{}) {
	pagination.Collection = append(pagination.Collection, obj)
}
