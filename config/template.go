package config

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type TemplateManager struct {
	templates *template.Template
}

func firstChar(s string) string {
	if len(s) > 0 {
		return string(s[0])
	}
	return "" // Return empty string if input is empty
}

func subFn(a int, b int) int {
	return a - b
}

func addFn(a int, b int) int {
	return a + b
}

var funcMap = template.FuncMap{
	"firstChar": firstChar,
	"sub":       subFn,
	"add":       addFn,
}

func NewTemplateManager() error {
	tmpl := template.Must(template.New("").Funcs(funcMap).ParseGlob(filepath.Join(TEMPLATE_DIR, "*.html")))
	TMPL = &TemplateManager{templates: tmpl}
	return nil
}

func (tm *TemplateManager) Render(w http.ResponseWriter, tmpl string, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	err := tm.templates.ExecuteTemplate(w, tmpl, data)
	return NewInternalError(err)
}
