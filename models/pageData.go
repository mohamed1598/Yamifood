package models

import "yamifood/pkg/forms"

type PageData struct {
	StrMap          map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	DataMap         map[string]interface{}
	CSRFToken       string
	Warning         string
	Error           string
	Form            *forms.Form
	Data            map[string]interface{}
	IsAuthenticated int
	IsAdmin         int
}
