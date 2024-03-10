package services

import (
	"context"
	"database/sql"
	"log"
	"time"

	"tomasweigenast.com/url-shortener/database"
	"tomasweigenast.com/url-shortener/entities"
	"tomasweigenast.com/url-shortener/mapper"
	"tomasweigenast.com/url-shortener/models"
	"tomasweigenast.com/url-shortener/utils"
)

type linksService struct {
	cache *utils.Cache
}

var linksServiceSingleton *linksService

func LinksService() *linksService {
	if linksServiceSingleton == nil {
		linksServiceSingleton = &linksService{
			cache: utils.NewCache(utils.CacheConfig{DefaultTtl: 1 * time.Hour, MaxEntries: 100, RenewTtlOnHit: true}),
		}
	}

	return linksServiceSingleton
}

func (*linksService) ListLinks(ctx context.Context, uid uint32) (*[]models.Url, error) {
	links, err := database.ListLinks(ctx, uid)
	if err != nil {
		return nil, err
	}

	urls := make([]models.Url, len(links))
	for i, link := range links {
		urls[i] = *mapper.MapUrl(&link)
	}

	return &urls, nil
}

func (*linksService) CreateLink(ctx context.Context, uid uint32, model models.CreateLink) (*models.Url, error) {
	// validate
	err := models.ModelValidator.Struct(model)
	if err != nil {
		return nil, err
	}

	urlPath := utils.RandomString(10)
	url := &entities.Url{
		Id:      utils.RandomId(),
		Uid:     uid,
		Link:    model.Link,
		UrlPath: urlPath,
		Name:    sql.Null[string]{V: model.Name, Valid: len(model.Name) > 0},
	}

	err = database.InsertUrl(ctx, url)
	if err != nil {
		return nil, err
	}

	return mapper.MapUrl(url), nil
}

func (ls *linksService) DeleteLink(ctx context.Context, id, uid uint32) error {
	path, err := database.DeleteUrl(ctx, id, uid)
	if err != nil {
		return err
	}

	ls.cache.Delete(path)
	return nil
}

// FetchUrl returns the redirect url for the given path
func (ls *linksService) FetchUrl(ctx context.Context, path string) (string, error) {
	// fetch from cache first
	url := ls.cache.Get(path)
	if url != nil {
		log.Println("cache hit:", path)
		return url.(string), nil
	}

	urlData, err := database.GetUrlByPath(ctx, path)
	if err != nil {
		return "", err
	}

	ls.cache.Put(path, urlData.Link)

	return urlData.Link, nil
}
