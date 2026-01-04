package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/daveonkels/clinkding/internal/client"
	"github.com/daveonkels/clinkding/internal/models"
)

type TagsAPI struct {
	client *client.Client
}

func NewTagsAPI(client *client.Client) *TagsAPI {
	return &TagsAPI{client: client}
}

func (a *TagsAPI) List(ctx context.Context, limit, offset int) (*models.TagList, error) {
	params := url.Values{}
	if limit > 0 {
		params.Add("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		params.Add("offset", strconv.Itoa(offset))
	}

	path := a.client.BuildURL("/api/tags/", params)

	var result models.TagList
	if err := a.client.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (a *TagsAPI) Get(ctx context.Context, id int) (*models.Tag, error) {
	path := fmt.Sprintf("/api/tags/%d/", id)
	var tag models.Tag
	if err := a.client.Get(ctx, path, &tag); err != nil {
		return nil, err
	}
	return &tag, nil
}

func (a *TagsAPI) Create(ctx context.Context, name string) (*models.Tag, error) {
	tagCreate := &models.TagCreate{Name: name}
	var result models.Tag
	if err := a.client.Post(ctx, "/api/tags/", tagCreate, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
