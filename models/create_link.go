package models

type CreateLink struct {
	Name string `json:"name" validate:"ascii,max=20"`
	Link string `json:"link" validate:"required,url"`
}
