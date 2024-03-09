package entities

import (
	"database/sql"
	"time"
)

type Url struct {
	Id         uint32              `db:"id"`
	Uid        uint32              `db:"uid"`
	Name       sql.Null[string]    `db:"name"`
	Link       string              `db:"link"`
	UrlPath    string              `db:"url_path"`
	CreatedAt  time.Time           `db:"created_at"`
	ValidUntil sql.Null[time.Time] `db:"valid_until"`
	IsDeleted  bool                `db:"is_deleted"`
}
