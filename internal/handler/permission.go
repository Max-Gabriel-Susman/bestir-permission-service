package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Max-Gabriel-Susman/bestir-permissionmaking-service/internal/foundation/web"
	"github.com/Max-Gabriel-Susman/bestir-permissionmaking-service/internal/permission"
	"github.com/go-chi/chi/v5"
)

type permission struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type permissionGroup struct {
	*permission.API
}

type ListpermissionsResponse struct {
	permissions []permission.permission `json:"permissions"`
}

func permissionEndpoints(app *web.App, api *permission.API) {
	ag := permissionGroup{API: api}

	app.Handle("GET", "/permission/{id}", ag.Getpermission)
	app.Handle("GET", "/permission", ag.Listpermissiones)
	app.Handle("POST", "/permission", ag.Createpermission)
	app.Handle("DELETE", "/permission/{id}", ag.Createpermission)
	app.Handle("PUT", "/permission/{id}", ag.Createpermission)
}

func (ag permissionGroup) Listpermissiones(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	permissions, err := ag.API.Listpermissiones(ctx)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, ListpermissionsResponse{
		permissions: permissions,
	}, http.StatusOK)
}

// permissions := []permission.permission{
// 	{
// 		ID:      69,
// 		Balance: 69,
// 	},
// 	{
// 		ID:      420,
// 		Balance: 420,
// 	},
// }

func (ag permissionGroup) Createpermission(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	fmt.Println("Create permission a invoked")
	var input permission.Incomingpermission
	if err := web.Decode(r.Body, &input); err != nil {
		return err
	}

	permission, err := ag.API.Createpermission(ctx, input)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, permission, http.StatusCreated)
}

func (ag permissionGroup) Getpermission(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	permissionID := chi.URLParam(r, "permission_id")
	if permissionID == "" {
		return nil
		// return handleMissingURLParameter(ctx, permissionID, permission)
	}

	return nil
}

func (ag permissionGroup) Updatepermission(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	fmt.Println("Create permission a invoked")
	var input permission.Incomingpermission
	if err := web.Decode(r.Body, &input); err != nil {
		return err
	}

	permission, err := ag.API.Updatepermission(ctx, input)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, permission, http.StatusCreated)
}

func (ag permissionGroup) Deletepermission(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	fmt.Println("Create permission a invoked")
	var input permission.Incomingpermission
	if err := web.Decode(r.Body, &input); err != nil {
		return err
	}

	permission, err := ag.API.Deletepermission(ctx, input)
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, permission, http.StatusCreated)
}
