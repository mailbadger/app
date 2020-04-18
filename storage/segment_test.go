package storage

import (
	"testing"

	"github.com/news-maily/app/entities"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSegment(t *testing.T) {
	db := openTestDb()
	defer func() {
		err := db.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	store := From(db)

	//Test create list
	l := &entities.Segment{
		Name:   "foo",
		UserID: 1,
	}

	err := store.CreateSegment(l)
	assert.Nil(t, err)

	//Test get list
	l, err = store.GetSegment(l.ID, 1)
	assert.Nil(t, err)
	assert.Equal(t, l.Name, "foo")

	//Test update list
	l.Name = "bar"
	err = store.UpdateSegment(l)
	assert.Nil(t, err)
	assert.Equal(t, l.Name, "bar")

	//Test get list by name
	l, err = store.GetSegmentByName("bar", 1)
	assert.Nil(t, err)
	assert.Equal(t, l.Name, "bar")

	//Test list validation when name is invalid
	l.Name = ""
	l.Validate()
	assert.Equal(t, l.Errors["name"], "The segment name cannot be empty.")

	//Test get lists
	p := NewPaginationCursor("/api/segments", 10)
	err = store.GetSegments(1, p)
	assert.Nil(t, err)
	col := p.Collection.(*[]entities.Segment)
	assert.NotNil(t, col)
	assert.NotEmpty(t, *col)

	//Test append subscribers to list
	s := &entities.Subscriber{
		Name:   "john",
		Email:  "john@example.com",
		UserID: 1,
	}
	err = store.CreateSubscriber(s)
	assert.Nil(t, err)

	l.Subscribers = append(l.Subscribers, *s)

	err = store.AppendSubscribers(l)
	assert.Nil(t, err)

	// store.GetSubscribersBySegmentID(l.ID, l.UserID, p)
	// assert.Nil(t, err)
	// assert.NotEmpty(t, p.Collection)

	//Test detach subscribers from list
	err = store.DetachSubscribers(l)
	assert.Nil(t, err)

	// Test delete list
	err = store.DeleteSegment(1, 1)
	assert.Nil(t, err)
}
