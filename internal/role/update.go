package permission

import (
	"context"

	"github.com/google/uuid"
)

func (api *API) Updatepermission(ctx context.Context, incomingpermission Incomingpermission) (permission, error) {
	id := uuid.New()

	permission := permission{
		ID:   id,
		Name: incomingpermission.Name,
	}

	err := api.Store.Updatepermission(ctx, permission)

	return permission, err
}
