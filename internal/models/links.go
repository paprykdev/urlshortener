package models

import (
	"time"
)

type Link struct {
	ID        string    `json:"id"`
	Url       string    `json:"url"`
	ShortCode string    `json:"short_code"`
	CreatedAt time.Time `json:"created_at"`
}
