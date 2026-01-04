package models

import "time"

type Bookmark struct {
	ID                 int       `json:"id"`
	URL                string    `json:"url"`
	Title              string    `json:"title"`
	Description        string    `json:"description"`
	Notes              string    `json:"notes"`
	WebsiteTitle       string    `json:"website_title"`
	WebsiteDescription string    `json:"website_description"`
	IsArchived         bool      `json:"is_archived"`
	Unread             bool      `json:"unread"`
	Shared             bool      `json:"shared"`
	TagNames           []string  `json:"tag_names"`
	DateAdded          time.Time `json:"date_added"`
	DateModified       time.Time `json:"date_modified"`
}

type BookmarkList struct {
	Count    int        `json:"count"`
	Next     *string    `json:"next"`
	Previous *string    `json:"previous"`
	Results  []Bookmark `json:"results"`
}

type BookmarkCreate struct {
	URL         string   `json:"url"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Notes       string   `json:"notes,omitempty"`
	TagNames    []string `json:"tag_names,omitempty"`
	IsArchived  bool     `json:"is_archived,omitempty"`
	Unread      bool     `json:"unread,omitempty"`
	Shared      bool     `json:"shared,omitempty"`
}

type BookmarkUpdate struct {
	URL         string   `json:"url,omitempty"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Notes       string   `json:"notes,omitempty"`
	TagNames    []string `json:"tag_names,omitempty"`
	Unread      *bool    `json:"unread,omitempty"`
	Shared      *bool    `json:"shared,omitempty"`
}

type BookmarkCheck struct {
	Bookmark *Bookmark        `json:"bookmark"`
	Metadata BookmarkMetadata `json:"metadata"`
	AutoTags []string         `json:"auto_tags"`
}

type BookmarkMetadata struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
