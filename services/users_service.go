package services

import (
	"context"
	"errors"
	"fmt"

	"tomasweigenast.com/url-shortener/database"
	"tomasweigenast.com/url-shortener/entities"
	"tomasweigenast.com/url-shortener/mapper"
	"tomasweigenast.com/url-shortener/models"
	"tomasweigenast.com/url-shortener/utils"
)

type usersService struct{}

var usersServiceSingleton *usersService = &usersService{}

func UsersService() *usersService {
	return usersServiceSingleton
}

func (*usersService) CreateUser(ctx context.Context, model models.RegisterUser) (*models.User, error) {
	// validate
	err := models.ModelValidator.Struct(model)
	if err != nil {
		return nil, err
	}

	user := &entities.User{
		Id:           utils.RandomId(),
		Name:         model.Name,
		Email:        model.Email,
		PasswordHash: utils.HashPassword(model.Password),
	}

	err = database.InsertUser(ctx, user)
	if err != nil {
		err = utils.TransformPgError(err)
		if errors.Is(err, utils.ErrAlreadyExists) {
			return nil, fmt.Errorf("email %q already in use", model.Email)
		}

		return nil, err
	}

	return mapper.MapUser(user), nil
}

func (*usersService) GetUser(ctx context.Context, by GetUserBy) (*models.User, error) {
	var user *entities.User
	var err error

	if len(by.Email) != 0 {
		user, err = database.FindUserByEmail(ctx, by.Email)
	} else {
		user, err = database.FindUserById(ctx, by.Id)
	}

	if err != nil {
		return nil, utils.TransformPgError(err)
	}

	return mapper.MapUser(user), nil
}

type GetUserBy struct {
	Email string
	Id    uint32
}
