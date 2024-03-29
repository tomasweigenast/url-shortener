package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"tomasweigenast.com/url-shortener/entities"
	"tomasweigenast.com/url-shortener/utils"
)

func InsertUser(ctx context.Context, user *entities.User) error {
	tag, err := pool.Exec(ctx, queryInsertUser, user.Id, user.Name, user.Email, user.PasswordHash)
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
	err := pgxscan.Get(ctx, pool, &user, queryGetUserByEmail, email)
	return &user, utils.TransformPgError(err)
}

func FindUserById(ctx context.Context, uid uint32) (*entities.User, error) {
	var user entities.User
	err := pgxscan.Get(ctx, pool, &user, queryGetUserById, uid)
	return &user, utils.TransformPgError(err)
}

func InsertSession(ctx context.Context, session *entities.UserSession) error {
	tag, err := pool.Exec(ctx, queryInsertSession, session.Id, session.Uid, session.RefreshToken, session.ExpiresAt)
	if err != nil {
		return utils.TransformPgError(err)
	}

	if tag.RowsAffected() != 1 {
		return errors.New("unable to create session")
	}

	return nil
}

func InsertUrl(ctx context.Context, link *entities.Url) (err error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return utils.TransformPgError(err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, queryInsertUrl, link.Id, link.Uid, link.Link, link.UrlPath, link.Name, link.ValidUntil)
	if err != nil {
		return utils.TransformPgError(err)
	}

	_, err = tx.Exec(ctx, queryInsertUrlMetadata, link.Id)
	if err != nil {
		return utils.TransformPgError(err)
	}

	return nil
}

func DeleteUrl(ctx context.Context, id, uid uint32) (urlPath string, err error) {
	row := pool.QueryRow(ctx, queryDeleteLink, id, uid)

	err = row.Scan(&urlPath)
	if err != nil {
		return "", utils.TransformPgError(err)
	}

	return
}

// GetUrlByPath retrieves minimal information needed to redirect. It only returns the id and the link itself
func GetUrlByPath(ctx context.Context, path string) (*entities.Url, error) {
	var link entities.Url
	err := pgxscan.Get(ctx, pool, &link, queryGetLinkByPath, path)
	return &link, utils.TransformPgError(err)
}

func ListLinks(ctx context.Context, uid uint32) ([]entities.Url, error) {
	var links []entities.Url
	err := pgxscan.Select(ctx, pool, &links, queryGetUserLinks, uid)
	return links, utils.TransformPgError(err)
}

func InsertUrlHit(ctx context.Context, hit *entities.UrlHit) (err error) {

	tx, err := pool.Begin(ctx)
	if err != nil {
		return utils.TransformPgError(err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, queryInsertUrlHit, hit.Id, hit.UrlId, hit.HitAt, hit.FromIP, hit.FromCache, hit.HttpMethod, hit.Proto, hit.QueryParams, hit.Headers, hit.UserAgent, hit.Cookies)
	if err != nil {
		return utils.TransformPgError(err)
	}

	_, err = tx.Exec(ctx, queryUpdateUrlMetadata, hit.UrlId)
	if err != nil {
		return utils.TransformPgError(err)
	}

	return nil
}

func GetLinkData(ctx context.Context, id, uid uint32) (*entities.UrlData, error) {
	data := entities.UrlData{}
	err := pgxscan.Get(ctx, pool, &data, querySelectUrlMetadata, id, uid)
	if err != nil {
		return nil, utils.TransformPgError(err)
	}

	fmt.Println(data)

	return &data, utils.TransformPgError(err)
}
