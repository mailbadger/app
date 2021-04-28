package templates

import (
	"embed"
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin"
)

//go:embed html
var content embed.FS

var emailTmpls *template.Template

func init() {
	emailTmpls = template.Must(template.ParseFS(content, "html/email/*.html"))
}

// Init sets the embedded HTML files to the gin engine.
func Init(r *gin.Engine) error {
	templ, err := template.ParseFS(content, "html/*.html")
	if err != nil {
		return fmt.Errorf("templates: parse fs: %w", err)
	}

	r.SetHTMLTemplate(templ)
	return nil
}

// GetEmailTemplates returns the parsed email template files.
func GetEmailTemplates() *template.Template {
	return emailTmpls
}
