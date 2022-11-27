package permission

import (
	"context"

	"github.com/google/uuid"
)

func (api *API) Deletepermission(ctx context.Context, incomingpermission Incomingpermission) (permission, error) {
	id := uuid.New()

	permission := permission{
		ID:   id,
		Name: incomingpermission.Name,
	}

	err := api.Store.Deletepermission(ctx, permission)

	return permission, err
}
