package templates

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestInit(t *testing.T) {
	r := gin.New()

	err := Init(r)
	assert.Nil(t, err)

	rec := httptest.NewRecorder()

	err = r.HTMLRender.Instance("unsubscribe.html", gin.H{
		"email":  "foo@bar.com",
		"t":      "foo",
		"uuid":   "abcdefgh",
		"failed": false,
	}).Render(rec)

	assert.Nil(t, err)

	resp := rec.Result()
	defer resp.Body.Close()

	_, err = html.Parse(resp.Body)
	assert.Nil(t, err)
}
