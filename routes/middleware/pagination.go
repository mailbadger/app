package middleware

import (
	"fmt"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

func Paginate(c *gin.Context) {
	var pagination Pagination

	pagination.Page = 0
	pagination.PerPage = DefaultPerPage
	pagination.Total = math.MaxUint32
	pagination.Collection = make([]interface{}, 0)

	if len(c.Request.URL.Query()["per_page"]) > 0 {
		if len(c.Request.URL.Query()["per_page"]) > 1 {
			panic("More than one per_page parameter attached to get url")
		}
		perpage, err := strconv.ParseUint(c.Request.URL.Query()["per_page"][0], 10, 32)
		if err != nil {
			panic(fmt.Sprintf("Error parsing 'per_page': %s", err))
		}

		pagination.PerPage = uint(perpage)

		//Lock on 100 if the user requests more than 100 items per page
		if pagination.PerPage > 100 {
			pagination.PerPage = 100
		}
	}
	if len(c.Request.URL.Query()["page"]) > 0 {
		if len(c.Request.URL.Query()["page"]) > 1 {
			panic("More than one page parameter attached to get url")
		}
		page, err := strconv.ParseUint(c.Request.URL.Query()["page"][0], 10, 32)
		if err != nil {
			panic(fmt.Sprintf("Error parsing 'page': %s", err))
		}
		pagination.Page = uint(page)
		pagination.Offset = uint(page * uint64(pagination.PerPage))
	}

	c.Set("pagination", pagination)
	c.Next()
}
