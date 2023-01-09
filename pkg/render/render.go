package render

import (
	"fmt"
	"net/http"
	"text/template"
	"yamifood/models"
	"yamifood/pkg/config"
	"yamifood/pkg/helpers"

	"github.com/justinas/nosurf"
)

var tmplCache = make(map[string]*template.Template)
var app *config.AppConfig

func NewAppConfig(a *config.AppConfig) {
	app = a
}

func AddCSRFData(pd *models.PageData, r *http.Request) *models.PageData {
	pd.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		pd.IsAuthenticated = 1
		if app.Session.Exists(r.Context(), "is_admin") {
			pd.IsAdmin = 1
		}
	}
	return pd
}
func RenderTemplate(w http.ResponseWriter, r *http.Request, t string, pd *models.PageData) {
	var tmpl *template.Template
	var err error
	_, inMap := tmplCache[t]
	if !inMap {
		err = makeTemplateCache(t)
		helpers.ErrorCheck(err)
	} else {
		fmt.Println("template in cache")
	}
	tmpl = tmplCache[t]
	pd = AddCSRFData(pd, r)
	err = tmpl.Execute(w, pd)
	helpers.ErrorCheck(err)
}
func makeTemplateCache(t string) error {
	templates := []string{
		fmt.Sprintf("./templates/%s", t),
		"templates/base.layout.tmpl",
	}
	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}
	tmplCache[t] = tmpl
	return nil
}
