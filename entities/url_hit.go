package entities

import (
	"time"
)

type UrlHit struct {
	Id          uint32         `db:"id"`
	UrlId       uint32         `db:"url_id" json:"url_id"`
	HitAt       time.Time      `db:"hit_at" json:"hit_at"`
	FromIP      string         `db:"from_ip" json:"from_ip"`
	FromCache   bool           `db:"from_cache" json:"from_cache"`
	HttpMethod  string         `db:"http_method" json:"http_method"`
	Proto       string         `db:"proto"`
	QueryParams []string       `db:"query_params" json:"query_params"`
	Headers     []string       `db:"headers"`
	UserAgent   string         `db:"user_agent" json:"user_agent"`
	Cookies     map[string]any `db:"cookies"`
}
