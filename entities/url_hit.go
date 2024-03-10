package entities

import (
	"database/sql"
	"time"
)

type UrlHit struct {
	Id          uint32           `db:"id"`
	UrlId       uint32           `db:"url_id"`
	HitAt       time.Time        `db:"hit_at"`
	FromIP      string           `db:"from_ip"`
	FromCache   bool             `db:"from_cache"`
	HttpMethod  string           `db:"http_method"`
	Proto       string           `db:"proto"`
	QueryParams []string         `db:"query_params"`
	Headers     []string         `db:"headers"`
	UserAgent   sql.Null[string] `db:"user_agent"`
	Cookies     map[string]any   `db:"cookies"`
}
