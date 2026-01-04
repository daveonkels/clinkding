package api

import (
	"context"

	"github.com/daveonkels/clinkding/internal/client"
	"github.com/daveonkels/clinkding/internal/models"
)

type UserAPI struct {
	client *client.Client
}

func NewUserAPI(client *client.Client) *UserAPI {
	return &UserAPI{client: client}
}

func (a *UserAPI) GetProfile(ctx context.Context) (*models.UserProfile, error) {
	var profile models.UserProfile
	if err := a.client.Get(ctx, "/api/user/profile/", &profile); err != nil {
		return nil, err
	}
	return &profile, nil
}
