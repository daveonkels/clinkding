package api

import (
	"context"
	"fmt"

	"github.com/daveonkels/clinkding/internal/client"
	"github.com/daveonkels/clinkding/internal/models"
)

type BundlesAPI struct {
	client *client.Client
}

func NewBundlesAPI(client *client.Client) *BundlesAPI {
	return &BundlesAPI{client: client}
}

func (a *BundlesAPI) List(ctx context.Context) (*models.BundleList, error) {
	var result models.BundleList
	if err := a.client.Get(ctx, "/api/bundles/", &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (a *BundlesAPI) Get(ctx context.Context, id int) (*models.Bundle, error) {
	path := fmt.Sprintf("/api/bundles/%d/", id)
	var bundle models.Bundle
	if err := a.client.Get(ctx, path, &bundle); err != nil {
		return nil, err
	}
	return &bundle, nil
}

func (a *BundlesAPI) Create(ctx context.Context, bundle *models.BundleCreate) (*models.Bundle, error) {
	var result models.Bundle
	if err := a.client.Post(ctx, "/api/bundles/", bundle, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (a *BundlesAPI) Update(ctx context.Context, id int, bundle *models.BundleUpdate) (*models.Bundle, error) {
	path := fmt.Sprintf("/api/bundles/%d/", id)
	var result models.Bundle
	if err := a.client.Patch(ctx, path, bundle, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (a *BundlesAPI) Delete(ctx context.Context, id int) error {
	path := fmt.Sprintf("/api/bundles/%d/", id)
	return a.client.Delete(ctx, path)
}
