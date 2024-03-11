package entities

import (
	"time"
)

type UrlData struct {
	Url        Url       `db:",inline"`
	Hits       uint32    `db:"hits"`
	LastHitAt  time.Time `db:"last_hit_at"`
	LatestHits []UrlHit  `db:"latest_hits"`
}
