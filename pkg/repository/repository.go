package repository

import "yamifood/models"

type DatabaseRepo interface {
	CreateUser(newUser models.User) error
	AuthenticateUser(testEmail, testPassword string) (int, string, error)
	IsAdmin(id int) (bool, error)
	CreateRecipe(newRecipe models.Recipe) error
	GetAllRecipes() ([]models.Recipe, error)
	GetRecipeById(id int) (models.Recipe, error)
	DeleteRecipe(id int) error
	UpdateRecipe(recipe models.Recipe) error
	GetRecipesByType(recipeType string) ([]models.Recipe, error)
}
