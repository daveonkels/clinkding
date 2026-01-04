package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/daveonkels/clinkding/internal/client"
	"github.com/daveonkels/clinkding/internal/models"
)

type BookmarksAPI struct {
	client *client.Client
}

func NewBookmarksAPI(client *client.Client) *BookmarksAPI {
	return &BookmarksAPI{client: client}
}

func (a *BookmarksAPI) List(ctx context.Context, opts *ListOptions) (*models.BookmarkList, error) {
	params := url.Values{}

	if opts != nil {
		if opts.Query != "" {
			params.Add("q", opts.Query)
		}
		if opts.Limit > 0 {
			params.Add("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Offset > 0 {
			params.Add("offset", strconv.Itoa(opts.Offset))
		}
		if opts.ModifiedSince != "" {
			params.Add("modified_since", opts.ModifiedSince)
		}
		if opts.AddedSince != "" {
			params.Add("added_since", opts.AddedSince)
		}
		if opts.BundleID > 0 {
			params.Add("bundle", strconv.Itoa(opts.BundleID))
		}
	}

	path := "/api/bookmarks/"
	if opts != nil && opts.Archived {
		path = "/api/bookmarks/archived/"
	}

	path = a.client.BuildURL(path, params)

	var result models.BookmarkList
	if err := a.client.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (a *BookmarksAPI) Get(ctx context.Context, id int) (*models.Bookmark, error) {
	path := fmt.Sprintf("/api/bookmarks/%d/", id)
	var bookmark models.Bookmark
	if err := a.client.Get(ctx, path, &bookmark); err != nil {
		return nil, err
	}
	return &bookmark, nil
}

func (a *BookmarksAPI) Check(ctx context.Context, urlToCheck string) (*models.BookmarkCheck, error) {
	params := url.Values{}
	params.Add("url", urlToCheck)
	path := a.client.BuildURL("/api/bookmarks/check/", params)

	var result models.BookmarkCheck
	if err := a.client.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (a *BookmarksAPI) Create(ctx context.Context, bookmark *models.BookmarkCreate) (*models.Bookmark, error) {
	var result models.Bookmark
	if err := a.client.Post(ctx, "/api/bookmarks/", bookmark, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (a *BookmarksAPI) Update(ctx context.Context, id int, bookmark *models.BookmarkUpdate) (*models.Bookmark, error) {
	path := fmt.Sprintf("/api/bookmarks/%d/", id)
	var result models.Bookmark
	if err := a.client.Patch(ctx, path, bookmark, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (a *BookmarksAPI) Archive(ctx context.Context, id int) error {
	path := fmt.Sprintf("/api/bookmarks/%d/archive/", id)
	return a.client.Post(ctx, path, nil, nil)
}

func (a *BookmarksAPI) Unarchive(ctx context.Context, id int) error {
	path := fmt.Sprintf("/api/bookmarks/%d/unarchive/", id)
	return a.client.Post(ctx, path, nil, nil)
}

func (a *BookmarksAPI) Delete(ctx context.Context, id int) error {
	path := fmt.Sprintf("/api/bookmarks/%d/", id)
	return a.client.Delete(ctx, path)
}

type ListOptions struct {
	Query         string
	Limit         int
	Offset        int
	Archived      bool
	ModifiedSince string
	AddedSince    string
	BundleID      int
}
