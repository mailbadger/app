package pagination

import (
	"net/url"
	"strconv"

	"github.com/news-maily/app/entities"
)

var DefaultPerPage int64 = 10

type Links struct {
	Previous *string `json:"previous"`
	Next     *string `json:"next"`
}

type Cursor struct {
	StartingAfter int64  `json:"-"`
	EndingBefore  int64  `json:"-"`
	Path          string `json:"-"`
	PerPage       int64  `json:"per_page"`
	Links         Links  `json:"links"`
	Model         entities.Model
	Collection    []interface{} `json:"collection"`
	Results       interface{}
}

// Append appends the object to the page collection.
func (c *Cursor) Append(obj interface{}) {
	c.Collection = append(c.Collection, obj)
}

func (c *Cursor) SetResults(res interface{}) {
	c.Results = res
}

// PopulateLinks populates the Links property with the query params needed for the
// previous and next urls. It uses the BasePath and encodes the 'per_page', 'ending_before' and 'starting_after'
// query parameters needed to create the links.
func (c *Cursor) PopulateLinks(prevID, nextID int64) {
	if prevID != 0 {
		params := url.Values{}
		params.Add("per_page", strconv.FormatInt(c.PerPage, 10))
		params.Add("ending_before", strconv.FormatInt(prevID, 10))
		l := c.Path + "?" + params.Encode()
		c.Links.Previous = &l
	}
	if nextID != 0 {
		params := url.Values{}
		params.Add("per_page", strconv.FormatInt(c.PerPage, 10))
		params.Add("starting_after", strconv.FormatInt(nextID, 10))
		l := c.Path + "?" + params.Encode()
		c.Links.Next = &l
	}
}
