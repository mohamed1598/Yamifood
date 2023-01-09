package dbrepo

import (
	"context"
	"errors"
	"time"
	"yamifood/models"

	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDbRepository) CreateUser(newUser models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `insert into users (name,email,password,acct_created) values ($1,$2,$3,current_timestamp)`
	hashedPw, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), 12)
	_, err := m.Db.ExecContext(ctx, query, newUser.Name, newUser.Email, hashedPw)
	return err
}
func (m *postgresDbRepository) AuthenticateUser(userEmail, userPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var id int
	var hashedPw string
	query := `select id,password from users where email=$1`
	row := m.Db.QueryRowContext(ctx, query, userEmail)
	if err := row.Scan(&id, &hashedPw); err != nil {
		return id, "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPw), []byte(userPassword)); err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("password is not correct")
	} else if err != nil {
		return 0, "", err
	}
	return id, hashedPw, nil
}
func (m *postgresDbRepository) IsAdmin(id int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var email string
	query := `select email from users where id=$1`
	row := m.Db.QueryRowContext(ctx, query, id)
	if err := row.Scan(&email); err != nil {
		return false, err
	}
	return email == "mohamed@mail.com", nil
}
