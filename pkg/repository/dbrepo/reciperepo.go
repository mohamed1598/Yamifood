package dbrepo

import (
	"context"
	"time"
	"yamifood/models"
)

func (m *postgresDbRepository) CreateRecipe(newRecipe models.Recipe) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `insert into recipes(name ,description,imgUrl,price,type) values($1,$2,$3,$4,$5)`
	_, err := m.Db.ExecContext(ctx, query, newRecipe.Name, newRecipe.Description, newRecipe.ImgUrl, newRecipe.Price, newRecipe.Type)
	return err
}
func (m *postgresDbRepository) GetAllRecipes() ([]models.Recipe, error) {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := m.Db.Query(`select id,name,description,imgUrl,price,type from recipes`)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var recipes []models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		err = rows.Scan(&recipe.ID, &recipe.Name, &recipe.Description, &recipe.ImgUrl, &recipe.Price, &recipe.Type)
		if err != nil {
			panic(err)
		}
		recipes = append(recipes, recipe)
	}
	err = rows.Err()
	return recipes, err
}
func (m *postgresDbRepository) GetRecipesByType(recipeType string) ([]models.Recipe, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `select id,name,description,imgUrl,price,type from recipes where type=$1`
	rows, err := m.Db.QueryContext(ctx, query, recipeType)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var recipes []models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		err = rows.Scan(&recipe.ID, &recipe.Name, &recipe.Description, &recipe.ImgUrl, &recipe.Price, &recipe.Type)
		if err != nil {
			panic(err)
		}
		recipes = append(recipes, recipe)
	}
	err = rows.Err()
	return recipes, err
}

func (m *postgresDbRepository) GetRecipeById(id int) (models.Recipe, error) {
	var recipe models.Recipe
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()
	query := `select id,name,description,imgUrl,price,type from recipes where id=$1`
	row := m.Db.QueryRowContext(ctx, query, id)
	err := row.Scan(&recipe.ID, &recipe.Name, &recipe.Description, &recipe.ImgUrl, &recipe.Price, &recipe.Type)
	return recipe, err

}

func (m *postgresDbRepository) DeleteRecipe(id int) error {
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()
	query := `delete from recipes where id=$1`
	_, err := m.Db.ExecContext(ctx, query, id)
	return err
}

func (m *postgresDbRepository) UpdateRecipe(recipe models.Recipe) error {
	ctx, close := context.WithTimeout(context.Background(), 5*time.Second)
	defer close()
	query := `update recipes set name=$1,description=$2,imgUrl=$3,price=$4,type=$5 where id=$6`
	_, err := m.Db.ExecContext(ctx, query, &recipe.Name, &recipe.Description, &recipe.ImgUrl, &recipe.Price, &recipe.Type, &recipe.ID)
	return err
}
