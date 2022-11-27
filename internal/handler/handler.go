package handler

import (
	"net/http"

	"github.com/Max-Gabriel-Susman/bestir-permissionmaking-service/internal/foundation/database"
	"github.com/Max-Gabriel-Susman/bestir-permissionmaking-service/internal/foundation/web"
	"github.com/Max-Gabriel-Susman/bestir-permissionmaking-service/internal/permission"
)

var _ http.Handler = (*web.App)(nil)

// maybe we'll add gitsha and other params later
func API(d Deps) *web.App {
	app := web.NewApp()
	dbrConn := database.NewDBR(d.DB)
	permissionAPI := permission.NewAPI(permission.NewMySQLStore(dbrConn))
	permissionEndpoints(app, permissionAPI)
	return app
}
