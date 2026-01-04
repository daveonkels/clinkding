package models

import "time"

type Asset struct {
	ID           int       `json:"id"`
	BookmarkID   int       `json:"bookmark"`
	File         string    `json:"file"`
	DisplayName  string    `json:"display_name"`
	FileSize     int64     `json:"file_size"`
	Status       string    `json:"status"`
	DateCreated  time.Time `json:"date_created"`
}

type AssetList struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []Asset `json:"results"`
}
