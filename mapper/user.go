package mapper

import (
	"tomasweigenast.com/url-shortener/entities"
	"tomasweigenast.com/url-shortener/models"
)

func MapUser(e *entities.User) *models.User {
	return &models.User{
		Id:    e.Id,
		Name:  e.Name,
		Email: e.Email,
	}
}
