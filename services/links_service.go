package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"tomasweigenast.com/url-shortener/database"
	"tomasweigenast.com/url-shortener/entities"
	"tomasweigenast.com/url-shortener/mapper"
	"tomasweigenast.com/url-shortener/models"
	"tomasweigenast.com/url-shortener/taskmanager"
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

func (ls *linksService) RegisterHit(ctx context.Context, linkId uint32, req *http.Request) error {
	queryParams := []string{}
	headers := []string{}
	cookies := map[string]any{}
	userAgent := req.UserAgent()
	for k, val := range req.URL.Query() {
		queryParams = append(queryParams, fmt.Sprintf("%s:%s", k, strings.Join(val, ",")))
	}

	for k, val := range req.Header {
		headers = append(headers, fmt.Sprintf("%s:%s", k, strings.Join(val, ",")))
	}

	for _, cookie := range req.Cookies() {
		cookies[cookie.Name] = map[string]any{
			"path":      cookie.Path,
			"expires":   cookie.Expires,
			"domain":    cookie.Domain,
			"max_age":   cookie.MaxAge,
			"same_site": cookie.SameSite,
			"http_only": cookie.HttpOnly,
			"secure":    cookie.Secure,
			"value":     cookie.Value,
		}
	}

	hit := entities.UrlHit{
		Id:          utils.RandomId(),
		UrlId:       linkId,
		HitAt:       time.Now(),
		FromIP:      req.RemoteAddr,
		FromCache:   false,
		HttpMethod:  req.Method,
		Proto:       req.Proto,
		QueryParams: queryParams,
		Headers:     headers,
		UserAgent: sql.Null[string]{
			V:     userAgent,
			Valid: len(userAgent) > 0,
		},
		Cookies: cookies,
	}

	taskmanager.EnqueueUrlHit(hit)
	return nil
}

// FetchUrl returns the redirect url for the given path
func (ls *linksService) FetchUrl(ctx context.Context, path string) (url string, id uint32, err error) {
	// fetch from cache first
	url_data := ls.cache.Get(path)
	if url_data != nil {
		log.Println("cache hit:", path)
		data := url_data.(*entities.Url)
		return data.Link, data.Id, nil
	}

	urlData, err := database.GetUrlByPath(ctx, path)
	if err != nil {
		return "", 0, err
	}

	ls.cache.Put(path, urlData)

	return urlData.Link, urlData.Id, nil
}
