package api

import (
	"context"
	"fmt"

	"github.com/daveonkels/clinkding/internal/client"
	"github.com/daveonkels/clinkding/internal/models"
)

type AssetsAPI struct {
	client *client.Client
}

func NewAssetsAPI(client *client.Client) *AssetsAPI {
	return &AssetsAPI{client: client}
}

func (a *AssetsAPI) List(ctx context.Context, bookmarkID int) (*models.AssetList, error) {
	path := fmt.Sprintf("/api/bookmarks/%d/assets/", bookmarkID)
	var result models.AssetList
	if err := a.client.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (a *AssetsAPI) Get(ctx context.Context, bookmarkID, assetID int) (*models.Asset, error) {
	path := fmt.Sprintf("/api/bookmarks/%d/assets/%d/", bookmarkID, assetID)
	var asset models.Asset
	if err := a.client.Get(ctx, path, &asset); err != nil {
		return nil, err
	}
	return &asset, nil
}

func (a *AssetsAPI) Upload(ctx context.Context, bookmarkID int, filePath string) (*models.Asset, error) {
	path := fmt.Sprintf("/api/bookmarks/%d/assets/upload/", bookmarkID)
	var result models.Asset
	if err := a.client.UploadFile(ctx, path, filePath, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (a *AssetsAPI) Download(ctx context.Context, bookmarkID, assetID int, outputPath string) error {
	path := fmt.Sprintf("/api/bookmarks/%d/assets/%d/download/", bookmarkID, assetID)
	return a.client.DownloadFile(ctx, path, outputPath)
}

func (a *AssetsAPI) Delete(ctx context.Context, bookmarkID, assetID int) error {
	path := fmt.Sprintf("/api/bookmarks/%d/assets/%d/", bookmarkID, assetID)
	return a.client.Delete(ctx, path)
}
