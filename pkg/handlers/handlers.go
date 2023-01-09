package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"yamifood/models"
	"yamifood/pkg/config"
	"yamifood/pkg/dbdriver"
	"yamifood/pkg/forms"
	"yamifood/pkg/helpers"
	"yamifood/pkg/render"
	"yamifood/pkg/repository"
	"yamifood/pkg/repository/dbrepo"

	"github.com/go-chi/chi/v5"
)

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

var Repo *Repository

func NewRepository(ac *config.AppConfig, db *dbdriver.DB) *Repository {
	return &Repository{
		App: ac,
		DB:  dbrepo.NewPostgresRepo(db.SQL, ac),
	}
}
func NewHandler(r *Repository) {
	Repo = r
}

func (m *Repository) HomeHandler(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "index.page.tmpl", &models.PageData{})
}
func (m *Repository) AboutHandler(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "about.page.tmpl", &models.PageData{})
}
func (m *Repository) SignUpHandler(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "signup.page.tmpl", &models.PageData{
		Form: &forms.Form{
			Errors: make(map[string][]string),
		},
	})
}
func (m *Repository) PostSignUpHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	helpers.ErrorCheck(err)
	newUser := models.User{
		Name:     r.Form.Get("name"),
		Email:    r.Form.Get("email"),
		Password: r.Form.Get("password"),
	}
	form := forms.New(r.PostForm)
	form.HasRequired("name", "email", "password")
	form.MinLength("name", 5, r)
	form.MinLength("email", 5, r)
	form.MinLength("password", 5, r)
	form.IsEmail("email")
	if !form.Valid() {
		data := make(map[string]interface{})
		data["RegisteredUser"] = data
		render.RenderTemplate(w, r, "signup.page.tmpl", &models.PageData{Form: form, Data: data})
		return
	}
	err = m.DB.CreateUser(newUser)
	helpers.ErrorCheck(err)
	render.RenderTemplate(w, r, "signupsuccessfully.page.tmpl", &models.PageData{})
}
func (m *Repository) SignupSuccessfullyHandler(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "signupsuccessfully.page.tmpl", &models.PageData{})
}
func (m *Repository) LoginHandler(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "login.page.tmpl", &models.PageData{
		Form: &forms.Form{
			Errors: make(map[string][]string),
		},
	})
}
func (m *Repository) PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	helpers.ErrorCheck(err)
	testemail := r.Form.Get("email")
	testpassword := r.Form.Get("password")
	form := forms.New(r.PostForm)
	form.HasRequired("email", "password")
	form.MinLength("email", 5, r)
	form.MinLength("password", 5, r)
	form.IsEmail("email")
	if !form.Valid() {
		data := make(map[string]interface{})
		data["loggingInUser"] = models.User{Email: testemail, Password: testpassword}
		render.RenderTemplate(w, r, "login.page.tmpl", &models.PageData{Data: data, Form: form})
		return
	}
	id, _, err := m.DB.AuthenticateUser(testemail, testpassword)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", err.Error())
		data := make(map[string]interface{})
		data["loggingInUser"] = models.User{Email: testemail, Password: testpassword}
		data["error"] = err.Error()
		render.RenderTemplate(w, r, "login.page.tmpl", &models.PageData{Data: data, Form: form})
		return
	}
	isAdmin, _ := m.DB.IsAdmin(id)
	if isAdmin {
		m.App.Session.Put(r.Context(), "is_admin", isAdmin)
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "valid login")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (m *Repository) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (m *Repository) CreateRecipeHandler(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "createrecipe.page.tmpl", &models.PageData{
		Form: &forms.Form{
			Errors: make(map[string][]string),
		},
	})
}
func (m *Repository) PostCreateRecipeHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	isvalid, form := validateRecipe(r)
	if !isvalid {
		data := make(map[string]interface{})
		render.RenderTemplate(w, r, "createrecipe.page.tmpl", &models.PageData{Form: form, Data: data})
		return
	}
	recipePrice, err := strconv.ParseFloat(r.Form.Get("price"), 64)
	if err != nil {
		log.Println(err)
		return
	}
	recipe := models.Recipe{
		Name:        r.Form.Get("name"),
		Description: r.Form.Get("description"),
		ImgUrl:      r.Form.Get("imgUrl"),
		Price:       recipePrice,
		Type:        r.Form.Get("type"),
	}
	err = m.DB.CreateRecipe(recipe)
	if err != nil {
		log.Println(err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) ManageRecipeHandler(w http.ResponseWriter, r *http.Request) {
	userId := 0
	if m.App.Session.Exists(r.Context(), "user_id") {
		userId = m.App.Session.GetInt(r.Context(), "user_id")
	}
	isAdmin, _ := m.DB.IsAdmin(userId)
	if !isAdmin {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	recipes, err := m.DB.GetAllRecipes()
	if err != nil {
		fmt.Println(err)
	}
	data := make(map[string]interface{})
	data["recipes"] = recipes
	render.RenderTemplate(w, r, "managerecipe.page.tmpl", &models.PageData{Data: data})
}

func (m *Repository) DeleteRecipeHandler(w http.ResponseWriter, r *http.Request) {
	recipeId := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(recipeId, 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	err = m.DB.DeleteRecipe(int(id))
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/manage-recipe", http.StatusSeeOther)
}

func (m *Repository) EditRecipeHandler(w http.ResponseWriter, r *http.Request) {
	recipeId := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(recipeId, 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	recipe, err := m.DB.GetRecipeById(int(id))
	if err != nil {
		log.Fatal(err)
	}
	data := make(map[string]interface{})
	data["recipe"] = recipe
	// fmt.Println(data["recipe"].Name)
	render.RenderTemplate(w, r, "editrecipe.page.tmpl", &models.PageData{
		Form: &forms.Form{
			Errors: make(map[string][]string),
		},
		Data: data,
	})
}

func (m *Repository) PostEditRecipeHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	isvalid, _ := validateRecipe(r)
	if !isvalid {
		http.Redirect(w, r, "/edit-recipe/"+r.Form.Get("id"), http.StatusSeeOther)
		return
	}
	recipePrice, err := strconv.ParseFloat(r.Form.Get("price"), 64)
	if err != nil {
		log.Println(err)
		return
	}
	recipeId, _ := strconv.ParseInt(r.Form.Get("id"), 0, 0)
	recipe := models.Recipe{
		Name:        r.Form.Get("name"),
		Description: r.Form.Get("description"),
		ImgUrl:      r.Form.Get("imgUrl"),
		Price:       recipePrice,
		Type:        r.Form.Get("type"),
		ID:          int(recipeId),
	}
	err = m.DB.UpdateRecipe(recipe)
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/manage-recipe", http.StatusSeeOther)
}

func (m *Repository) MenuHandler(w http.ResponseWriter, r *http.Request) {
	recipeType := chi.URLParam(r, "type")
	var recipes []models.Recipe
	var err error
	if recipeType == "all" {
		recipes, err = m.DB.GetAllRecipes()
	} else {
		recipes, err = m.DB.GetRecipesByType(recipeType)
	}

	if err != nil {
		fmt.Println(err)
	}
	data := make(map[string]interface{})
	data["recipes"] = recipes
	data["recipeType"] = recipeType
	render.RenderTemplate(w, r, "menu.page.tmpl", &models.PageData{Data: data})
}

func validateRecipe(r *http.Request) (bool, *forms.Form) {
	form := forms.New(r.PostForm)
	form.HasRequired("name", "description", "imgUrl", "price", "type")
	form.MinLength("name", 5, r)
	form.MinLength("description", 20, r)
	form.MinLength("type", 5, r)
	return form.Valid(), form
}
