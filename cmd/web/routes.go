package main

import (
	"yamifood/pkg/config"
	"yamifood/pkg/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)

	mux.Use(NoSurf)
	mux.Use(SetupSession)
	//home page
	mux.Get("/", handlers.Repo.HomeHandler)
	mux.Get("/about", handlers.Repo.AboutHandler)

	//user
	mux.Get("/signup", handlers.Repo.SignUpHandler)
	mux.Post("/signup", handlers.Repo.PostSignUpHandler)
	mux.Get("/signupsuccessfully", handlers.Repo.SignupSuccessfullyHandler)
	mux.Get("/login", handlers.Repo.LoginHandler)
	mux.Post("/login", handlers.Repo.PostLoginHandler)
	mux.Get("/logout", handlers.Repo.LogoutHandler)

	// recipes
	mux.Get("/create-recipe", handlers.Repo.CreateRecipeHandler)
	mux.Post("/create-recipe", handlers.Repo.PostCreateRecipeHandler)
	mux.Get("/manage-recipe", handlers.Repo.ManageRecipeHandler)
	mux.Get("/delete-recipe/{id}", handlers.Repo.DeleteRecipeHandler)
	mux.Get("/edit-recipe/{id}", handlers.Repo.EditRecipeHandler)
	mux.Post("/edit-recipe", handlers.Repo.PostEditRecipeHandler)

	//menu
	mux.Get("/menu/{type}", handlers.Repo.MenuHandler)
	return mux
}
