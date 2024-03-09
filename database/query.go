package database

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"tomasweigenast.com/url-shortener/entities"
	"tomasweigenast.com/url-shortener/utils"
)

func InsertUser(ctx context.Context, user *entities.User) error {
	tag, err := connection.Exec(ctx, queryInsertUser, user.Id, user.Name, user.Email, user.PasswordHash)
	if err != nil {
		return utils.TransformPgError(err)
	}

	if tag.RowsAffected() != 1 {
		return errors.New("unable to create user")
	}

	return nil
}

func FindUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	err := pgxscan.Get(ctx, connection, &user, queryGetUserByEmail, email)
	return &user, utils.TransformPgError(err)
}

func FindUserById(ctx context.Context, uid uint32) (*entities.User, error) {
	var user entities.User
	err := pgxscan.Get(ctx, connection, &user, queryGetUserById, uid)
	return &user, utils.TransformPgError(err)
}

func InsertSession(ctx context.Context, session *entities.UserSession) error {
	tag, err := connection.Exec(ctx, queryInsertSession, session.Id, session.Uid, session.RefreshToken, session.ExpiresAt)
	if err != nil {
		return utils.TransformPgError(err)
	}

	if tag.RowsAffected() != 1 {
		return errors.New("unable to create session")
	}

	return nil
}

func InsertUrl(ctx context.Context, link *entities.Url) error {
	tag, err := connection.Exec(ctx, queryInsertUrl, link.Id, link.Uid, link.Link, link.UrlPath, link.Name, link.ValidUntil)
	if err != nil {
		return utils.TransformPgError(err)
	}

	if tag.RowsAffected() != 1 {
		return errors.New("unable to create url")
	}

	return nil
}

// GetUrlByPath retrieves minimal information needed to redirect. It only returns the id and the link itself
func GetUrlByPath(ctx context.Context, path string) (*entities.Url, error) {
	var link entities.Url
	err := pgxscan.Get(ctx, connection, &link, queryGetLinkByPath, path)
	return &link, utils.TransformPgError(err)
}

func ListLinks(ctx context.Context, uid uint32) ([]entities.Url, error) {
	var links []entities.Url
	err := pgxscan.Select(ctx, connection, &links, queryGetUserLinks, uid)
	return links, utils.TransformPgError(err)
}
