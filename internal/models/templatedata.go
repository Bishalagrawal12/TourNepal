package models

import (
	"text/template"

	"github.com/Atul-Ranjan12/tourism/internal/forms"
)

type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
	Flash           string
	Warning         string
	Error           string
	Form            *forms.Form
	IsAuthenticated int
	TemplateFuncs   template.FuncMap
}
