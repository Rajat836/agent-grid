package app

import "github.com/gin-gonic/gin"

var appVersion = "1.0.0"

// SetVersion sets the application version to be returned in health checks
func SetVersion(version string) {
	appVersion = version
}

// GetVersion returns the current application version
func GetVersion() string {
	return appVersion
}

func RegisterRoutes(engine *gin.Engine) {
	// Health routes (public, no auth required)
	registerHealthRoutes(engine, appVersion)

	// TODO: Register other feature routes here
	// registerEventRoutes(engine)
	// registerUserRoutes(engine)
}
