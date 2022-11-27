package permission

import (
	"context"

	"github.com/google/uuid"
)

func (api *API) Createpermission(ctx context.Context, incomingpermission Incomingpermission) (permission, error) {
	id := uuid.New()

	permission := permission{
		ID:   id,
		Name: incomingpermission.Name,
	}

	err := api.Store.Createpermission(ctx, permission)

	return permission, err
}
