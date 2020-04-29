package storage

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/news-maily/app/entities"
)

// Direction of the pagination.
type Direction int

// Possible directions for the pagination.
const (
	Start Direction = iota
	Forward
	Backward
)

// DefaultPerPage is a default value of number of items per page
var DefaultPerPage int64 = 10

// Links represent the previous and next links used when iterating through the
// collection.
type Links struct {
	Previous *string `json:"previous"`
	Next     *string `json:"next"`
}

// PaginationCursor represents the paginated results by the given model.
type PaginationCursor struct {
	Scopes        []func(*gorm.DB) *gorm.DB `json:"-"`
	Query         *gorm.DB                  `json:"-"`
	StartingAfter int64                     `json:"-"`
	EndingBefore  int64                     `json:"-"`
	Path          string                    `json:"-"`
	Resource      string                    `json:"-"`
	Direction     Direction                 `json:"-"`
	PerPage       int64                     `json:"per_page"`
	Total         int64                     `json:"total"`
	Links         Links                     `json:"links"`
	Collection    interface{}               `json:"collection"`
}

// NewPaginationCursor creates new PaginationCursor object.
func NewPaginationCursor(path string, perPage int64) *PaginationCursor {
	if perPage <= 0 || perPage > 100 {
		perPage = DefaultPerPage
	}

	return &PaginationCursor{
		Path:      path,
		PerPage:   perPage,
		Direction: Start,
	}
}

// PopulateLinks populates the Links property with the query params needed for the
// previous and next urls. It uses the BasePath and encodes the 'per_page', 'ending_before' and 'starting_after'
// query parameters needed to create the links.
func (c *PaginationCursor) PopulateLinks(prevID, nextID string) {

	c.Links = Links{}

	if prevID != "" && prevID != "0" {
		params := url.Values{}
		params.Add("per_page", strconv.FormatInt(c.PerPage, 10))
		params.Add("ending_before", prevID)
		l := c.Path + "?" + params.Encode()
		c.Links.Previous = &l
	}
	if nextID != "" && nextID != "0" {
		params := url.Values{}
		params.Add("per_page", strconv.FormatInt(c.PerPage, 10))
		params.Add("starting_after", nextID)
		l := c.Path + "?" + params.Encode()
		c.Links.Next = &l
	}
}

// SetCollection sets the collection in the cursor. Usually when setting a collection, it is empty, and
// gets populated when invoking the Paginate() method.
func (c *PaginationCursor) SetCollection(collection interface{}) {
	c.Collection = collection
}

// SetScopes sets the pagination query scopes.
func (c *PaginationCursor) SetScopes(scopes []func(*gorm.DB) *gorm.DB) {
	c.Scopes = scopes
}

// SetQuery sets the main query.
func (c *PaginationCursor) SetQuery(query *gorm.DB) {
	c.Query = query
}

// SetResource sets the pagination resource.
func (c *PaginationCursor) SetResource(r string) {
	c.Resource = r
}

// SetTotal sets the total number of items in the collection.
func (c *PaginationCursor) SetTotal(total int64) {
	c.Total = total
}

// SetStartingAfter sets the ID of the object that the page should start after.
func (c *PaginationCursor) SetStartingAfter(sa int64) {
	c.StartingAfter = sa
	c.Direction = Forward
}

// SetEndingBefore sets the ID of the object that the page should end before.
func (c *PaginationCursor) SetEndingBefore(eb int64) {
	c.EndingBefore = eb
	c.Direction = Backward
}

// SetPerPage sets the number for total items to be fetched per page.
func (c *PaginationCursor) SetPerPage(perPage int64) {
	if perPage <= 0 || perPage > 100 {
		perPage = DefaultPerPage
	}
	c.PerPage = perPage
}

func (db *store) Paginate(p *PaginationCursor, userID int64) error {
	var last *entities.Model

	switch p.Direction {
	case Backward:
		m, err := db.GetOne(p.EndingBefore, userID, p.Resource, p.Scopes...)
		if err != nil {
			return fmt.Errorf("paginate: get one: %w", err)
		}
		p.Query.Joins(fmt.Sprintf("INNER JOIN (?) as r ON %s.id = r.id", p.Resource),
			db.DB.
				Table(p.Resource).
				Select("id").
				Where("(created_at > ? OR (created_at = ? AND id > ?)) AND created_at < ?",
					m.CreatedAt,
					m.CreatedAt,
					m.ID,
					time.Now(),
				).Order("created_at, id asc").Limit(p.PerPage).QueryExpr(),
		).Find(p.Collection)

		last, err = db.GetLast(userID, p.Resource, p.Scopes...)
		if err != nil {
			return fmt.Errorf("paginate: get last: %w", err)
		}

	case Forward:
		m, err := db.GetOne(p.StartingAfter, userID, p.Resource, p.Scopes...)
		if err != nil {
			return fmt.Errorf("paginate: get one: %w", err)
		}

		p.Query.Where(`(created_at < ? OR (created_at = ? AND id < ?)) AND created_at < ?`,
			m.CreatedAt,
			m.CreatedAt,
			m.ID,
			time.Now(),
		).Find(p.Collection)

		// when it is descending order we'll need the first record (last from behind) in order
		// to check if it matches the last record from the current page. If they're the same
		// the 'next' link will be nil.
		last, err = db.GetFirst(userID, p.Resource, p.Scopes...)
		if err != nil {
			return fmt.Errorf("paginate: get first: %w", err)
		}
	case Start:
		p.Query.Find(p.Collection)
	}

	total, err := db.GetTotal(userID, p.Resource, p.Scopes...)
	if err != nil {
		return fmt.Errorf("paginate: get total: %w", err)
	}

	p.SetTotal(total)

	models := interfaceToSlice(p.Collection)

	prevID, nextID, err := findPrevAndNextIDs(p, models, last)
	if err != nil {
		return fmt.Errorf("paginate: find prev and next: %w", err)
	}

	p.PopulateLinks(strconv.FormatInt(prevID, 10), strconv.FormatInt(nextID, 10))

	return nil
}

func (db *store) GetOne(id, userID int64, table string, scopes ...func(*gorm.DB) *gorm.DB) (*entities.Model, error) {
	var model entities.Model
	scopes = append(scopes, BelongsToUser(userID))
	err := db.Table(table).Scopes(scopes...).Where("id = ?", id).Find(&model).Error
	return &model, err
}

func (db *store) GetTotal(userID int64, table string, scopes ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var count int64
	scopes = append(scopes, BelongsToUser(userID))
	err := db.Table(table).Scopes(scopes...).Count(&count).Error
	return count, err
}

func (db *store) GetFirst(userID int64, table string, scopes ...func(*gorm.DB) *gorm.DB) (*entities.Model, error) {
	var model entities.Model
	scopes = append(scopes, BelongsToUser(userID))
	err := db.Table(table).
		Scopes(scopes...).
		Order("created_at, id").
		Limit(1).
		Find(&model).
		Error
	return &model, err
}

func (db *store) GetLast(userID int64, table string, scopes ...func(*gorm.DB) *gorm.DB) (*entities.Model, error) {
	var model entities.Model
	scopes = append(scopes, BelongsToUser(userID))
	err := db.Table(table).
		Scopes(scopes...).
		Order("created_at desc, id desc").
		Limit(1).
		Find(&model).
		Error
	return &model, err
}

func findPrevAndNextIDs(p *PaginationCursor, models []interface{}, last *entities.Model) (int64, int64, error) {
	var (
		prevID, nextID int64
	)

	if p.Direction == Start {
		if len(models) == int(p.PerPage) && len(models) < int(p.Total) {
			nextID = models[len(models)-1].(entities.ModelInterface).GetID()
		}

		return prevID, nextID, nil
	}

	if len(models) > 0 {
		if p.Direction == Backward {
			nextID = models[len(models)-1].(entities.ModelInterface).GetID()
			if len(models) == int(p.PerPage) {
				firstFromCol := models[0].(entities.ModelInterface)
				if last.ID != firstFromCol.GetID() {
					prevID = firstFromCol.GetID()
				}
			}
		}

		if p.Direction == Forward {
			prevID = models[0].(entities.ModelInterface).GetID()
			if len(models) == int(p.PerPage) {
				lastFromCol := models[len(models)-1].(entities.ModelInterface)
				if last.ID != lastFromCol.GetID() {
					nextID = lastFromCol.GetID()
				}
			}
		}
	}

	return prevID, nextID, nil
}

func interfaceToSlice(slicePtr interface{}) []interface{} {
	ptr := reflect.ValueOf(slicePtr)
	if ptr.Kind() != reflect.Ptr {
		panic("interfaceSlice() is not given a pointer type")
	}
	ind := reflect.Indirect(ptr)

	if ind.Kind() != reflect.Slice {
		panic("interfaceSlice() indirect type is not a slice")
	}

	ret := make([]interface{}, ind.Len())

	for i := 0; i < ind.Len(); i++ {
		ret[i] = ind.Index(i).Interface()
	}

	return ret
}
