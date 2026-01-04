package models

import "time"

type Tag struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	DateAdded     time.Time `json:"date_added"`
	BookmarkCount int       `json:"bookmark_count"`
}

type TagList struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []Tag   `json:"results"`
}

type TagCreate struct {
	Name string `json:"name"`
}
