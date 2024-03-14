package entity

import (
	"time"
)

type Schedule struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	PublishAt  time.Time `json:"publish_at"`
	UnlockedAt time.Time `json:"unlocked_at"`
	IsUnlocked bool      `json:"is_unlocked"`
}
