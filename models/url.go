package models

import "time"

type Url struct {
	Id         uint32     `json:"id"`
	Uid        uint32     `json:"uid"`
	Name       *string    `json:"name,omitempty"`
	ValidUntil *time.Time `json:"validUntil,omitempty"`
	Url        string     `json:"url"`
	RedirectTo string     `json:"redirectTo"`
}

type UrlMetadata struct {
	Url        Url       `json:"url,inline"`
	Hits       uint32    `json:"hits"`
	LastHitAt  time.Time `json:"lastHitAt"`
	LatestHits []UrlHit  `json:"latestHits"`
}

type UrlHit struct {
	DateTime        time.Time      `json:"dateTime"`
	FromIP          string         `json:"fromIp,omitempty"`
	HttpMethod      string         `json:"httpMethod"`
	HttpProtocol    string         `json:"httpProtocol"`
	UserAgent       string         `json:"userAgent"`
	QueryParameters []string       `json:"queryParameters"`
	Headers         []string       `json:"headers"`
	Cookies         map[string]any `json:"cookies"`
}
