package app

import (
	"app/agent_grid/cmd/app/middlewares"

	"github.com/gin-gonic/gin"
)

const (
	PublicApiV1    = "/public/v1/"
	ProtectedApiV1 = "/v1"
	PrivateApiV1   = "/private/v1"
)

func (app *App) addRoutes(router *gin.Engine, middlewares *middlewares.Middlewares) {
	app.addBaseRoutes(router, middlewares)
}
