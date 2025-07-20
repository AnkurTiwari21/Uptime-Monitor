package router

import (
	"net/http"

	controllers "github.com/ankur12345678/uptime-monitor/Controllers"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes add all routing list here automatically get main router
func RegisterRoutes(ctrl controllers.BaseController) {
	ctrl.Router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Route Not Found"})
	})

	ctrl.Router.GET("/health", func(ctx *gin.Context) {
		// Send a ping to make sure the database connection is alive.
		db, err := ctrl.DB.DB()
		if err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"live": "not ok"})
			return
		}
		err = db.PingContext(ctx)
		if err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{"live": "not ok"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"live": "ok"})
	})

	// Register All routes
	InitRoutes(&ctrl)
}
