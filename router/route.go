package router

import (
	controllers "github.com/ankur12345678/uptime-monitor/Controllers"
	"github.com/ankur12345678/uptime-monitor/Controllers/middlewares"
	"github.com/ankur12345678/uptime-monitor/pkg/logger"
)

func InitRoutes(ctrl *controllers.BaseController) {
	v1RouteGroup := ctrl.Router.Group("/v1")

	//AUTH routes
	v1RouteGroup.POST("/signup", ctrl.SignUpHandler)
	v1RouteGroup.POST("/login", ctrl.LoginHandler)
	v1RouteGroup.POST("/refresh", middlewares.HandleAuth, ctrl.HandleRefresh)
	v1RouteGroup.POST("/logout", middlewares.HandleAuth, ctrl.HandleLogOut)

	fullAuthV1Routes := v1RouteGroup.Group("", middlewares.HandleAuth)

	//Website regitering/testing routes
	fullAuthV1Routes.POST("/register-website", ctrl.RegisterWebsite)
	fullAuthV1Routes.POST("/test-website", ctrl.TestWebsiteLiveliness)

	logger.Info("Initializing Routes : Success.....")
}
