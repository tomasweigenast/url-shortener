package entities

import "time"

type UrlMetadata struct {
	Id         uint32    `db:"id"`
	Hits       uint32    `db:"hits"`
	LastUpdate time.Time `db:"last_update"`
}
