package mapper

import (
	"log"
	"net/url"
	"os"
	"time"

	"tomasweigenast.com/url-shortener/entities"
	"tomasweigenast.com/url-shortener/models"
)

func MapUrl(e *entities.Url) *models.Url {
	baseUrl := os.Getenv("BASE_URL")
	if len(baseUrl) == 0 {
		log.Fatalf("BASE_URL environment variable is not present")
	}

	var name *string
	var validUntil *time.Time

	if e.Name.Valid {
		name = &e.Name.V
	}

	if e.ValidUntil.Valid {
		validUntil = &e.ValidUntil.V
	}

	url, err := url.JoinPath(baseUrl, e.UrlPath)
	if err != nil {
		log.Fatalf("unable to join path, this should not happen: %s", err)
	}

	return &models.Url{
		Id:         e.Id,
		Uid:        e.Uid,
		Name:       name,
		ValidUntil: validUntil,
		Url:        url,
		RedirectTo: e.Link,
	}
}
