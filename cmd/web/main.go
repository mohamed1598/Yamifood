package main

import (
	"encoding/gob"
	"net/http"
	"time"
	"yamifood/models"
	"yamifood/pkg/config"
	"yamifood/pkg/dbdriver"
	"yamifood/pkg/handlers"
	"yamifood/pkg/helpers"
	"yamifood/pkg/render"

	"github.com/alexedwards/scs/v2"
)

var sessionManager *scs.SessionManager
var app config.AppConfig

const connectionString = "host=localhost port=5432 dbname=yamifood_db user=postgres password=123456"

func main() {
	configSessions()
	db, err := configDb()
	helpers.ErrorCheck(err)
	defer db.SQL.Close()
	configHandlers(db)
	configRender()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	helpers.ErrorCheck(err)
}
func configHandlers(db *dbdriver.DB) {
	repo := handlers.NewRepository(&app, db)
	handlers.NewHandler(repo)
}
func configSessions() {
	gob.Register(models.User{})

	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.Secure = false
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	app.Session = sessionManager
}

func configRender() {
	render.NewAppConfig(&app)
}
func configDb() (*dbdriver.DB, error) {
	db, err := dbdriver.ConnectSQL(connectionString)
	return db, err
}
