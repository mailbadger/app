package templates

import (
	"embed"
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin"
)

//go:embed html
var content embed.FS

// Init sets the embedded HTML files to the gin engine.
func Init(r *gin.Engine) error {
	templ, err := template.ParseFS(content, "html/*.html")
	if err != nil {
		return fmt.Errorf("templates: parse fs: %w", err)
	}
	r.SetHTMLTemplate(templ)
	return nil
}
