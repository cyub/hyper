package router

import (
	"net/http"
	"net/http/pprof"
	"runtime"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ProfileProvider for regsiter router
func ProfileProvider() func(r *gin.Engine) {
	return func(r *gin.Engine) {
		profile := r.Group("/_profile")
		{
			profile.GET("/", pprofHandler(pprof.Index))
			profile.GET("/cmdline", pprofHandler(pprof.Cmdline))
			profile.GET("/profile", pprofHandler(pprof.Profile))
			profile.POST("/symbol", pprofHandler(pprof.Symbol))
			profile.GET("/symbol", pprofHandler(pprof.Symbol))
			profile.GET("/trace", pprofHandler(pprof.Trace))
			profile.GET("/allocs", pprofHandler(pprof.Handler("allocs").ServeHTTP))
			profile.GET("/block", pprofHandler(pprof.Handler("block").ServeHTTP))
			profile.GET("/goroutine", pprofHandler(pprof.Handler("goroutine").ServeHTTP))
			profile.GET("/heap", pprofHandler(pprof.Handler("heap").ServeHTTP))
			profile.GET("/mutex", pprofHandler(pprof.Handler("mutex").ServeHTTP))
			profile.GET("/threadcreate", pprofHandler(pprof.Handler("threadcreate").ServeHTTP))
			profile.PUT("/setblockrate", setBlockRateHandler)
		}
	}
}

func pprofHandler(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func setBlockRateHandler(c *gin.Context) {
	rate, err := strconv.Atoi(c.PostForm("rate"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	runtime.SetBlockProfileRate(rate)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
