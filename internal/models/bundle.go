package models

import "time"

type Bundle struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DateAdded   time.Time `json:"date_added"`
}

type BundleList struct {
	Count    int      `json:"count"`
	Next     *string  `json:"next"`
	Previous *string  `json:"previous"`
	Results  []Bundle `json:"results"`
}

type BundleCreate struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type BundleUpdate struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}
