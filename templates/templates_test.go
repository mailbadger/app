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

	// render unsubscribe.html
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

	// render unsubscribe-success.html
	err = r.HTMLRender.Instance("unsubscribe-success.html", gin.H{}).Render(rec)

	assert.Nil(t, err)
	resp = rec.Result()
	defer resp.Body.Close()

	_, err = html.Parse(resp.Body)
	assert.Nil(t, err)

	err = r.HTMLRender.Instance("non-existent-file.html", gin.H{}).Render(rec)
	assert.NotNil(t, err)
}
