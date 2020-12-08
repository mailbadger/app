package templates

var (
	bucket = "bucketo za templates"
)

type Service interface {
}

type service struct {
}

func NewTemplateService() Service {
	return &service{}
}
