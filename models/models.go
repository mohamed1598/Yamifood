package models

import (
	"time"

)

type User struct {
	Name        string
	Email       string
	Password    string
	AcctCreated time.Time
	LastLogin   time.Time
	Id          int
}
 type Recipe struct{
	Name string
	Description string
	ImgUrl string
	Price float64
	Type string
	ID int
 }